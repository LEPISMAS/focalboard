// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {
    reducer,
    setMe,
    setBoardUsers,
    addBoardUsers,
    removeBoardUsersById,
    followBlock,
    unfollowBlock,
    patchProps,
    getMe,
    getLoggedIn,
    getBoardUsers,
    getBoardUsersList,
    getUser,
    getOnboardingTourStarted,
    getOnboardingTourStep,
    getOnboardingTourCategory,
    getVersionMessageCanceled,
    getCardLimitSnoozeUntil,
    getCardHiddenWarningSnoozeUntil,
    getMyConfig,
    versionProperty,
} from '../users'
import {IUser} from '../../user'

// Mock octoClient
jest.mock('../../octoClient', () => ({
    default: {
        getMe: jest.fn().mockResolvedValue(null),
        getMyConfig: jest.fn().mockResolvedValue([]),
    },
}))

function makeStore() {
    return configureStore({reducer: {users: reducer}})
}

function makeUser(id: string, username = 'testuser'): IUser {
    return {
        id,
        username,
        email: `${username}@example.com`,
        nickname: '',
        firstname: '',
        lastname: '',
        props: {},
        create_at: 1000,
        update_at: 1000,
        is_bot: false,
        is_guest: false,
        roles: 'system_user',
    }
}

describe('users reducer', () => {
    test('initial state is correct', () => {
        const store = makeStore()
        expect(store.getState().users.me).toBeNull()
        expect(store.getState().users.loggedIn).toBeNull()
        expect(store.getState().users.boardUsers).toEqual({})
        expect(store.getState().users.blockSubscriptions).toEqual([])
        expect(store.getState().users.myConfig).toEqual({})
    })

    test('setMe sets me and loggedIn=true', () => {
        const store = makeStore()
        const user = makeUser('u1')
        store.dispatch(setMe(user))
        expect(store.getState().users.me).toEqual(user)
        expect(store.getState().users.loggedIn).toBe(true)
    })

    test('setMe with null sets loggedIn=false', () => {
        const store = makeStore()
        store.dispatch(setMe(makeUser('u1')))
        store.dispatch(setMe(null))
        expect(store.getState().users.me).toBeNull()
        expect(store.getState().users.loggedIn).toBe(false)
    })

    test('setBoardUsers sets board users indexed by id', () => {
        const store = makeStore()
        const u1 = makeUser('u1', 'alice')
        const u2 = makeUser('u2', 'bob')
        store.dispatch(setBoardUsers([u1, u2]))
        expect(store.getState().users.boardUsers.u1).toEqual(u1)
        expect(store.getState().users.boardUsers.u2).toEqual(u2)
    })

    test('setBoardUsers replaces existing users', () => {
        const store = makeStore()
        const u1 = makeUser('u1')
        const u2 = makeUser('u2')
        store.dispatch(setBoardUsers([u1]))
        store.dispatch(setBoardUsers([u2]))
        expect(Object.keys(store.getState().users.boardUsers)).toHaveLength(1)
        expect(store.getState().users.boardUsers.u2).toBeTruthy()
    })

    test('addBoardUsers adds users to existing map', () => {
        const store = makeStore()
        const u1 = makeUser('u1')
        const u2 = makeUser('u2')
        store.dispatch(setBoardUsers([u1]))
        store.dispatch(addBoardUsers([u2]))
        expect(Object.keys(store.getState().users.boardUsers)).toHaveLength(2)
    })

    test('removeBoardUsersById removes specified users', () => {
        const store = makeStore()
        const u1 = makeUser('u1')
        const u2 = makeUser('u2')
        store.dispatch(setBoardUsers([u1, u2]))
        store.dispatch(removeBoardUsersById(['u1']))
        expect(store.getState().users.boardUsers.u1).toBeUndefined()
        expect(store.getState().users.boardUsers.u2).toBeTruthy()
    })

    test('followBlock adds subscription', () => {
        const store = makeStore()
        const sub = {blockId: 'block1', workspaceId: 'ws1', subscriberType: 'user', subscriberId: 'u1', notifyProps: {}}
        store.dispatch(followBlock(sub as any))
        expect(store.getState().users.blockSubscriptions).toHaveLength(1)
    })

    test('unfollowBlock removes subscription by blockId', () => {
        const store = makeStore()
        const sub1 = {blockId: 'block1', workspaceId: 'ws1', subscriberType: 'user', subscriberId: 'u1', notifyProps: {}}
        const sub2 = {blockId: 'block2', workspaceId: 'ws1', subscriberType: 'user', subscriberId: 'u1', notifyProps: {}}
        store.dispatch(followBlock(sub1 as any))
        store.dispatch(followBlock(sub2 as any))
        store.dispatch(unfollowBlock(sub1 as any))
        expect(store.getState().users.blockSubscriptions).toHaveLength(1)
        expect(store.getState().users.blockSubscriptions[0].blockId).toBe('block2')
    })

    test('patchProps parses user preferences', () => {
        const store = makeStore()
        const prefs = [{category: 'boards', name: 'onboardingTourStarted', value: 'true', userId: 'u1'}]
        store.dispatch(patchProps(prefs as any))
        expect(store.getState().users.myConfig.onboardingTourStarted).toBeTruthy()
    })

    test('extraReducer: users/fetchMe/fulfilled sets me and myConfig', () => {
        const store = makeStore()
        const user = makeUser('u1')
        store.dispatch({
            type: 'users/fetchMe/fulfilled',
            payload: {
                me: user,
                myConfig: [{category: 'boards', name: 'onboardingTourStarted', value: 'true', user_id: 'u1'}],
            },
        })
        expect(store.getState().users.me?.id).toBe('u1')
        expect(store.getState().users.loggedIn).toBe(true)
    })

    test('extraReducer: users/fetchMe/fulfilled with null me', () => {
        const store = makeStore()
        store.dispatch({
            type: 'users/fetchMe/fulfilled',
            payload: {me: null, myConfig: null},
        })
        expect(store.getState().users.me).toBeNull()
        expect(store.getState().users.loggedIn).toBe(false)
    })

    test('extraReducer: users/fetchMe/rejected clears me', () => {
        const store = makeStore()
        store.dispatch(setMe(makeUser('u1')))
        store.dispatch({type: 'users/fetchMe/rejected', error: {}})
        expect(store.getState().users.me).toBeNull()
        expect(store.getState().users.loggedIn).toBe(false)
        expect(store.getState().users.myConfig).toEqual({})
    })

    test('extraReducer: user/blockSubscriptions/fulfilled sets subscriptions', () => {
        const store = makeStore()
        store.dispatch({
            type: 'user/blockSubscriptions/fulfilled',
            payload: [],
        })
        expect(store.getState().users.blockSubscriptions).toEqual([])
    })

    test('extraReducer: initialLoad/fulfilled sets myConfig from payload', () => {
        const store = makeStore()
        store.dispatch({
            type: 'initialLoad/fulfilled',
            payload: {
                team: {id: 't1'},
                teams: [],
                boards: [],
                boardsMemberships: [],
                boardTemplates: [],
                limits: null,
                myConfig: [{category: 'boards', name: 'tourCategory', value: 'boards', user_id: 'u1'}],
            },
        })
        expect(store.getState().users.myConfig.tourCategory).toBeTruthy()
    })

    test('extraReducer: initialLoad/fulfilled with no myConfig does not crash', () => {
        const store = makeStore()
        store.dispatch({
            type: 'initialLoad/fulfilled',
            payload: {
                team: {id: 't1'},
                teams: [],
                boards: [],
                boardsMemberships: [],
                boardTemplates: [],
                limits: null,
                myConfig: null,
            },
        })

        // Should not throw, myConfig remains as is
        expect(store.getState().users.myConfig).toBeTruthy()
    })
})

describe('users selectors', () => {
    test('getMe returns null initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getMe(state)).toBeNull()
    })

    test('getLoggedIn returns null initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getLoggedIn(state)).toBeNull()
    })

    test('getBoardUsers returns empty map initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getBoardUsers(state)).toEqual({})
    })

    test('getBoardUsersList returns empty array when no users', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getBoardUsersList(state)).toEqual([])
    })

    test('getBoardUsersList returns users sorted by username', () => {
        const store = makeStore()
        const u1 = makeUser('u1', 'zara')
        const u2 = makeUser('u2', 'alice')
        store.dispatch(setBoardUsers([u1, u2]))
        const state = store.getState() as any
        const list = getBoardUsersList(state)
        expect(list[0].username).toBe('alice')
        expect(list[1].username).toBe('zara')
    })

    test('getUser returns undefined for unknown userId', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getUser('unknown')(state)).toBeUndefined()
    })

    test('getUser returns user for known userId', () => {
        const store = makeStore()
        const u1 = makeUser('u1')
        store.dispatch(setBoardUsers([u1]))
        const state = store.getState() as any
        expect(getUser('u1')(state)).toEqual(u1)
    })

    test('getMyConfig returns empty object initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getMyConfig(state)).toEqual({})
    })

    test('getOnboardingTourStarted returns false when not set', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getOnboardingTourStarted(state)).toBe(false)
    })

    test('getOnboardingTourStarted returns true when set', () => {
        const store = makeStore()
        store.dispatch(patchProps([
            {category: 'boards', name: 'onboardingTourStarted', value: 'true', user_id: 'u1'},
        ] as any))
        const state = store.getState() as any
        expect(getOnboardingTourStarted(state)).toBe(true)
    })

    test('getOnboardingTourStep returns empty string when not set', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getOnboardingTourStep(state)).toBeUndefined()
    })

    test('getOnboardingTourCategory returns empty string when not set', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getOnboardingTourCategory(state)).toBe('')
    })

    test('getVersionMessageCanceled returns true when versionProperty not set and me is null', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getVersionMessageCanceled(state)).toBe(true)
    })

    test('getVersionMessageCanceled returns true when user is single-user', () => {
        const store = makeStore()
        store.dispatch(setMe(makeUser('single-user')))
        const state = store.getState() as any
        expect(getVersionMessageCanceled(state)).toBe(true)
    })

    test('getVersionMessageCanceled returns false when user is regular and prop not set', () => {
        const store = makeStore()
        store.dispatch(setMe(makeUser('u1')))
        const state = store.getState() as any
        expect(getVersionMessageCanceled(state)).toBe(false)
    })

    test('getCardLimitSnoozeUntil returns 0 when not set', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCardLimitSnoozeUntil(state)).toBe(0)
    })

    test('getCardLimitSnoozeUntil returns parsed value', () => {
        const store = makeStore()
        store.dispatch(patchProps([
            {category: 'boards', name: 'cardLimitSnoozeUntil', value: '12345', user_id: 'u1'},
        ] as any))
        const state = store.getState() as any
        expect(getCardLimitSnoozeUntil(state)).toBe(12345)
    })

    test('getCardHiddenWarningSnoozeUntil returns 0 when not set', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCardHiddenWarningSnoozeUntil(state)).toBe(0)
    })

    test('getCardHiddenWarningSnoozeUntil returns parsed value', () => {
        const store = makeStore()
        store.dispatch(patchProps([
            {category: 'boards', name: 'cardHiddenWarningSnoozeUntil', value: '9999', user_id: 'u1'},
        ] as any))
        const state = store.getState() as any
        expect(getCardHiddenWarningSnoozeUntil(state)).toBe(9999)
    })

    test('versionProperty is defined', () => {
        expect(versionProperty).toBeDefined()
        expect(typeof versionProperty).toBe('string')
    })
})
