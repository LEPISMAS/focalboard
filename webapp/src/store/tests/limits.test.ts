// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {reducer, setCardLimitTimestamp, getLimits, getCardLimitTimestamp} from '../limits'

function makeStore() {
    return configureStore({reducer: {limits: reducer}})
}

const defaultLimits = {
    cards: 0,
    used_cards: 0,
    card_limit_timestamp: 0,
    views: 0,
}

describe('limits reducer', () => {
    test('initial state has default limits', () => {
        const store = makeStore()
        expect(store.getState().limits.limits).toEqual(defaultLimits)
    })

    test('setLimits action (via direct dispatch) sets limits', () => {
        const store = makeStore()
        const newLimits = {cards: 100, used_cards: 50, card_limit_timestamp: 12345, views: 10}
        store.dispatch({type: 'limits/setLimits', payload: newLimits})
        expect(store.getState().limits.limits).toEqual(newLimits)
    })

    test('setCardLimitTimestamp updates only the timestamp', () => {
        const store = makeStore()
        const newLimits = {cards: 100, used_cards: 50, card_limit_timestamp: 12345, views: 10}
        store.dispatch({type: 'limits/setLimits', payload: newLimits})
        store.dispatch(setCardLimitTimestamp(99999))
        expect(store.getState().limits.limits.card_limit_timestamp).toBe(99999)
        expect(store.getState().limits.limits.cards).toBe(100)
    })

    test('extraReducer: initialLoad.fulfilled sets limits from payload', () => {
        const store = makeStore()
        const limits = {cards: 200, used_cards: 75, card_limit_timestamp: 54321, views: 20}
        store.dispatch({
            type: 'initialLoad/fulfilled',
            payload: {
                team: {id: 'team1', title: 'Team 1'},
                teams: [],
                boards: [],
                boardsMemberships: [],
                boardTemplates: [],
                limits,
                myConfig: [],
            },
        })
        expect(store.getState().limits.limits).toEqual(limits)
    })

    test('extraReducer: initialLoad.fulfilled with no limits uses default', () => {
        const store = makeStore()
        store.dispatch({
            type: 'initialLoad/fulfilled',
            payload: {
                team: {id: 'team1', title: 'Team 1'},
                teams: [],
                boards: [],
                boardsMemberships: [],
                boardTemplates: [],
                limits: null,
                myConfig: [],
            },
        })
        expect(store.getState().limits.limits).toEqual(defaultLimits)
    })

    test('setCardLimitTimestamp to 0 resets timestamp', () => {
        const store = makeStore()
        store.dispatch(setCardLimitTimestamp(99999))
        store.dispatch(setCardLimitTimestamp(0))
        expect(store.getState().limits.limits.card_limit_timestamp).toBe(0)
    })
})

describe('getLimits selector', () => {
    test('returns default limits initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getLimits(state)).toEqual(defaultLimits)
    })

    test('returns updated limits after dispatch', () => {
        const store = makeStore()
        const limits = {cards: 50, used_cards: 10, card_limit_timestamp: 111, views: 5}
        store.dispatch({type: 'limits/setLimits', payload: limits})
        const state = store.getState() as any
        expect(getLimits(state)).toEqual(limits)
    })
})

describe('getCardLimitTimestamp selector', () => {
    test('returns 0 initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCardLimitTimestamp(state)).toBe(0)
    })

    test('returns updated timestamp after setCardLimitTimestamp', () => {
        const store = makeStore()
        store.dispatch(setCardLimitTimestamp(777))
        const state = store.getState() as any
        expect(getCardLimitTimestamp(state)).toBe(777)
    })
})
