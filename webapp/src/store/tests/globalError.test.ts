// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {reducer, setGlobalError, getGlobalError} from '../globalError'

// We mock initialLoad to avoid async thunk execution via client
jest.mock('../initialLoad', () => ({
    initialReadOnlyLoad: {
        rejected: {
            type: 'initialReadOnlyLoad/rejected',
            match: (action: any) => action.type === 'initialReadOnlyLoad/rejected',
        },
        fulfilled: {
            type: 'initialReadOnlyLoad/fulfilled',
            match: (action: any) => action.type === 'initialReadOnlyLoad/fulfilled',
        },
    },
    initialLoad: {
        rejected: {
            type: 'initialLoad/rejected',
            match: (action: any) => action.type === 'initialLoad/rejected',
        },
        fulfilled: {
            type: 'initialLoad/fulfilled',
            match: (action: any) => action.type === 'initialLoad/fulfilled',
        },
    },
    loadBoardData: {
        fulfilled: {
            type: 'loadBoardData/fulfilled',
            match: (action: any) => action.type === 'loadBoardData/fulfilled',
        },
        pending: {
            type: 'loadBoardData/pending',
            match: (action: any) => action.type === 'loadBoardData/pending',
        },
        rejected: {
            type: 'loadBoardData/rejected',
            match: (action: any) => action.type === 'loadBoardData/rejected',
        },
    },
    loadBoards: {
        fulfilled: {
            type: 'loadBoards/fulfilled',
            match: (action: any) => action.type === 'loadBoards/fulfilled',
        },
    },
    loadMyBoardsMemberships: {
        fulfilled: {
            type: 'loadMyBoardsMemberships/fulfilled',
            match: (action: any) => action.type === 'loadMyBoardsMemberships/fulfilled',
        },
    },
    getUserBlockSubscriptions: jest.fn(),
    getUserBlockSubscriptionList: jest.fn(),
}))

function makeStore() {
    return configureStore({reducer: {globalError: reducer}})
}

describe('globalError reducer', () => {
    test('initial state is empty string', () => {
        const store = makeStore()
        expect(store.getState().globalError.value).toBe('')
    })

    test('setGlobalError sets the error value', () => {
        const store = makeStore()
        store.dispatch(setGlobalError('Something went wrong'))
        expect(store.getState().globalError.value).toBe('Something went wrong')
    })

    test('setGlobalError can clear the error', () => {
        const store = makeStore()
        store.dispatch(setGlobalError('error'))
        store.dispatch(setGlobalError(''))
        expect(store.getState().globalError.value).toBe('')
    })

    test('setGlobalError replaces existing error', () => {
        const store = makeStore()
        store.dispatch(setGlobalError('first error'))
        store.dispatch(setGlobalError('second error'))
        expect(store.getState().globalError.value).toBe('second error')
    })

    test('extraReducer: initialReadOnlyLoad.rejected sets error from action.error.message', () => {
        const store = makeStore()
        store.dispatch({
            type: 'initialReadOnlyLoad/rejected',
            error: {message: 'read only load failed'},
        })
        expect(store.getState().globalError.value).toBe('read only load failed')
    })

    test('extraReducer: initialLoad.rejected sets error from action.error.message', () => {
        const store = makeStore()
        store.dispatch({
            type: 'initialLoad/rejected',
            error: {message: 'initial load failed'},
        })
        expect(store.getState().globalError.value).toBe('initial load failed')
    })

    test('extraReducer: rejected with no message sets empty string', () => {
        const store = makeStore()
        store.dispatch({
            type: 'initialLoad/rejected',
            error: {},
        })
        expect(store.getState().globalError.value).toBe('')
    })
})

describe('getGlobalError selector', () => {
    test('returns the current global error value', () => {
        const store = makeStore()
        store.dispatch(setGlobalError('test error'))
        const state = store.getState() as any
        expect(getGlobalError(state)).toBe('test error')
    })

    test('returns empty string initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getGlobalError(state)).toBe('')
    })
})

