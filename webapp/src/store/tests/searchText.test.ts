// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {reducer, setSearchText, getSearchText} from '../searchText'

function makeStore() {
    return configureStore({reducer: {searchText: reducer}})
}

describe('searchText reducer', () => {
    test('initial state is empty string', () => {
        const store = makeStore()
        expect(store.getState().searchText.value).toBe('')
    })

    test('setSearchText sets value', () => {
        const store = makeStore()
        store.dispatch(setSearchText('hello'))
        expect(store.getState().searchText.value).toBe('hello')
    })

    test('setSearchText replaces existing value', () => {
        const store = makeStore()
        store.dispatch(setSearchText('first'))
        store.dispatch(setSearchText('second'))
        expect(store.getState().searchText.value).toBe('second')
    })

    test('setSearchText with empty string clears value', () => {
        const store = makeStore()
        store.dispatch(setSearchText('hello'))
        store.dispatch(setSearchText(''))
        expect(store.getState().searchText.value).toBe('')
    })
})

describe('getSearchText selector', () => {
    test('returns current search text', () => {
        const store = makeStore()
        store.dispatch(setSearchText('test search'))
        const state = store.getState() as any
        expect(getSearchText(state)).toBe('test search')
    })

    test('returns empty string initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getSearchText(state)).toBe('')
    })
})
