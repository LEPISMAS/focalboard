// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {reducer, setChannel, getCurrentChannel, ChannelTypeOpen, ChannelTypePrivate, ChannelTypeDirectMessage, ChannelTypeGroupMessage} from '../channels'

function makeStore() {
    return configureStore({reducer: {channels: reducer}})
}

describe('channels constants', () => {
    test('ChannelTypeOpen is "O"', () => {
        expect(ChannelTypeOpen).toBe('O')
    })

    test('ChannelTypePrivate is "P"', () => {
        expect(ChannelTypePrivate).toBe('P')
    })

    test('ChannelTypeDirectMessage is "D"', () => {
        expect(ChannelTypeDirectMessage).toBe('D')
    })

    test('ChannelTypeGroupMessage is "G"', () => {
        expect(ChannelTypeGroupMessage).toBe('G')
    })
})

describe('channels reducer', () => {
    test('initial state has current as null', () => {
        const store = makeStore()
        expect(store.getState().channels.current).toBeNull()
    })

    test('setChannel sets the current channel', () => {
        const store = makeStore()
        const channel = {id: 'ch1', name: 'general', display_name: 'General', type: ChannelTypeOpen}
        store.dispatch(setChannel(channel))
        expect(store.getState().channels.current).toEqual(channel)
    })

    test('setChannel does nothing if same channel object reference', () => {
        const store = makeStore()
        const channel = {id: 'ch1', name: 'general', display_name: 'General', type: ChannelTypeOpen}
        store.dispatch(setChannel(channel))
        const stateBefore = store.getState().channels.current
        store.dispatch(setChannel(channel))
        expect(store.getState().channels.current).toEqual(stateBefore)
    })

    test('setChannel updates to a different channel', () => {
        const store = makeStore()
        const channel1 = {id: 'ch1', name: 'general', display_name: 'General', type: ChannelTypeOpen}
        const channel2 = {id: 'ch2', name: 'random', display_name: 'Random', type: ChannelTypePrivate}
        store.dispatch(setChannel(channel1))
        store.dispatch(setChannel(channel2))
        expect(store.getState().channels.current).toEqual(channel2)
    })

    test('setChannel works with Direct Message type', () => {
        const store = makeStore()
        const channel = {id: 'dm1', name: 'dm', display_name: 'DM', type: ChannelTypeDirectMessage}
        store.dispatch(setChannel(channel))
        expect(store.getState().channels.current).toEqual(channel)
    })

    test('setChannel works with Group Message type', () => {
        const store = makeStore()
        const channel = {id: 'gm1', name: 'gm', display_name: 'GM', type: ChannelTypeGroupMessage}
        store.dispatch(setChannel(channel))
        expect(store.getState().channels.current).toEqual(channel)
    })
})

describe('getCurrentChannel selector', () => {
    test('returns null initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentChannel(state)).toBeNull()
    })

    test('returns the current channel after setting', () => {
        const store = makeStore()
        const channel = {id: 'ch1', name: 'general', display_name: 'General', type: ChannelTypeOpen}
        store.dispatch(setChannel(channel))
        const state = store.getState() as any
        expect(getCurrentChannel(state)).toEqual(channel)
    })
})
