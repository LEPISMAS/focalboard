// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/* eslint-disable @typescript-eslint/no-explicit-any */

import {
    ACTION_AUTH,
    ACTION_SUBSCRIBE_BLOCKS,
    ACTION_SUBSCRIBE_TEAM,
    ACTION_UNSUBSCRIBE_BLOCKS,
    ACTION_UNSUBSCRIBE_TEAM,
    MMWebSocketClient,
    WSClient,
} from './wsclient'
import {Utils} from './utils'
import {TestBlockFactory} from './test/testBlockFactory'

console.log = jest.fn()
console.error = jest.fn()
console.warn = jest.fn()

type PluginListeners = {
    connect?: () => void
    reconnect?: () => void
    error?: (event: Event) => void
    close?: (count: number) => void
}

function pluginClient(listeners: PluginListeners): MMWebSocketClient {
    return {
        conn: {readyState: 1} as WebSocket,
        sendMessage: jest.fn(),
        addFirstConnectListener: jest.fn((callback) => {
            listeners.connect = callback
        }),
        addReconnectListener: jest.fn((callback) => {
            listeners.reconnect = callback
        }),
        addErrorListener: jest.fn((callback) => {
            listeners.error = callback
        }),
        addCloseListener: jest.fn((callback) => {
            listeners.close = callback
        }),
    }
}

describe('WSClient', () => {
    beforeEach(() => {
        jest.useFakeTimers()
        jest.spyOn(Utils, 'assertFailure').mockImplementation(() => undefined)
    })

    afterEach(() => {
        jest.restoreAllMocks()
        jest.useRealTimers()
    })

    test('manages plugin connection, authentication, team, and block subscriptions', () => {
        const listeners: PluginListeners = {}
        const plugin = pluginClient(listeners)
        const client = new WSClient('http://focalboard.test')
        const states = jest.fn()
        const reconnect = jest.fn()
        const errors = jest.fn()

        client.initPlugin('focalboard', '1.0.0', plugin)
        client.addOnStateChange(states)
        client.addOnReconnect(reconnect)
        client.addOnError(errors)
        client.subscribeToTeam('team-1')
        client.subscribeToTeam('team-1')
        client.open()
        listeners.connect?.()

        expect(states).toHaveBeenCalledWith(client, 'open')
        expect(plugin.sendMessage).toHaveBeenCalledWith(
            `custom_focalboard_${ACTION_SUBSCRIBE_TEAM}`,
            {teamId: 'team-1'},
        )

        client.authenticate('session-token')
        client.subscribeToBlocks('team-1', ['block-1'], 'read-token')
        client.unsubscribeFromBlocks('team-1', ['block-1'], 'read-token')
        client.unsubscribeToTeam('team-1')
        client.unsubscribeToTeam('team-1')
        client.unsubscribeToTeam('missing')

        const actions = (plugin.sendMessage as jest.Mock).mock.calls.map(([action]) => action)
        expect(actions).toEqual(expect.arrayContaining([
            `custom_focalboard_${ACTION_AUTH}`,
            `custom_focalboard_${ACTION_SUBSCRIBE_BLOCKS}`,
            `custom_focalboard_${ACTION_UNSUBSCRIBE_BLOCKS}`,
            `custom_focalboard_${ACTION_UNSUBSCRIBE_TEAM}`,
        ]))

        const event = new Event('error')
        listeners.error?.(event)
        expect(errors).toHaveBeenCalledWith(client, event)

        listeners.reconnect?.()
        expect(reconnect).toHaveBeenCalledWith(client)

        listeners.close?.(1)
        expect(states).toHaveBeenLastCalledWith(client, 'close')
        client.close()
        expect(client.hasConn()).toBe(true)
    })

    test('rejects block commands without a connection and empty authentication', () => {
        const client = new WSClient('http://focalboard.test')
        client.authenticate('')
        client.subscribeToBlocks('team-1', ['block-1'])
        client.unsubscribeFromBlocks('team-1', ['block-1'])

        expect(Utils.assertFailure).toHaveBeenCalledTimes(3)
        expect(client.hasConn()).toBe(false)
        client.close()
    })

    test('adds, batches, dispatches, and removes every change handler type', () => {
        const client = new WSClient('http://focalboard.test')
        const handler = jest.fn()
        const card = TestBlockFactory.createCard()
        const messages = [
            [{block: card}, card, 'block'],
            [{category: {id: 'category-1'}}, {id: 'category-1'}, 'category'],
            [{blockCategories: [{boardID: 'board-1', categoryID: 'category-1'}]}, [{boardID: 'board-1', categoryID: 'category-1'}], 'blockCategories'],
            [{board: {id: 'board-1'}}, {id: 'board-1'}, 'board'],
            [{member: {boardId: 'board-1', userId: 'user-1'}}, {boardId: 'board-1', userId: 'user-1'}, 'boardMembers'],
            [{categoryOrder: ['category-1']}, ['category-1'], 'categoryOrder'],
        ] as const

        for (const [, , type] of messages) {
            client.addOnChange(handler, type as any)
        }

        const fixWSData = jest.spyOn(Utils, 'fixWSData')
        for (const [message, data, type] of messages) {
            fixWSData.mockReturnValueOnce([data as any, type as any])
            client.updateHandler(message as any)
        }
        client.updateHandler({teamId: 'different-team', block: card})

        jest.advanceTimersByTime(100)
        expect(handler).toHaveBeenCalledTimes(6)

        for (const [, , type] of messages) {
            client.removeOnChange(handler, type as any)
            client.removeOnChange(handler, type as any)
        }
    })

    test('dispatches config, card-limit, follow, unfollow, and plugin-version events', () => {
        const listeners: PluginListeners = {}
        const plugin = pluginClient(listeners)
        const client = new WSClient('http://focalboard.test')
        const configHandler = jest.fn()
        const limitHandler = jest.fn()
        const follow = jest.fn()
        const unfollow = jest.fn()
        const versionChanged = jest.fn()

        client.addOnConfigChange(configHandler)
        client.addOnCardLimitTimestampChange(limitHandler)
        client.updateClientConfigHandler({} as any)
        client.updateCardLimitTimestampHandler({action: 'limit', timestamp: 42})
        expect(configHandler).toHaveBeenCalledWith(client, {})
        expect(limitHandler).toHaveBeenCalledWith(client, 42)
        client.removeOnConfigChange(configHandler)
        client.removeOnCardLimitTimestampChange(limitHandler)

        client.setOnFollowBlock(follow)
        client.setOnUnfollowBlock(unfollow)
        client.updateSubscriptionHandler({action: 'follow'} as any)
        client.updateSubscriptionHandler({action: 'follow', subscription: {blockId: 'block-1'}} as any)
        client.updateSubscriptionHandler({action: 'unfollow', subscription: {blockId: 'block-1', deleteAt: 1}} as any)
        expect(follow).toHaveBeenCalledTimes(1)
        expect(unfollow).toHaveBeenCalledTimes(1)

        client.initPlugin('focalboard', '1.0.0', plugin)
        client.open()
        client.setOnAppVersionChangeHandler(versionChanged)
        client.pluginStatusesChangedHandler({plugin_statuses: [{plugin_id: 'other', version: '2.0.0'}]})
        client.pluginStatusesChangedHandler({plugin_statuses: [{plugin_id: 'focalboard', version: '2.0.0'}]})
        expect(versionChanged).toHaveBeenCalledWith(true)
        jest.advanceTimersByTime(1000)

        client.removeOnReconnect(jest.fn())
        client.removeOnStateChange(jest.fn())
        client.removeOnError(jest.fn())
    })
})
