// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {reducer, getLanguage} from '../language'

// Mock i18n module to avoid real browser storage
jest.mock('../../i18n', () => ({
    getCurrentLanguage: jest.fn().mockReturnValue('es'),
    storeLanguage: jest.fn(),
}))

function makeStore() {
    return configureStore({reducer: {language: reducer}})
}

describe('language reducer', () => {
    test('initial state is "en"', () => {
        const store = makeStore()
        expect(store.getState().language.value).toBe('en')
    })

    test('setLanguage action (via direct dispatch) updates language', () => {
        const store = makeStore()
        store.dispatch({type: 'language/setLanguage', payload: 'fr'})
        expect(store.getState().language.value).toBe('fr')
    })

    test('extraReducer: language/fetch/fulfilled sets language from payload', () => {
        const store = makeStore()
        store.dispatch({
            type: 'language/fetch/fulfilled',
            payload: 'de',
        })
        expect(store.getState().language.value).toBe('de')
    })

    test('extraReducer: language/store/fulfilled sets language from payload', () => {
        const store = makeStore()
        store.dispatch({
            type: 'language/store/fulfilled',
            payload: 'ja',
        })
        expect(store.getState().language.value).toBe('ja')
    })

    test('language can be set to different values sequentially', () => {
        const store = makeStore()
        store.dispatch({type: 'language/fetch/fulfilled', payload: 'es'})
        expect(store.getState().language.value).toBe('es')
        store.dispatch({type: 'language/store/fulfilled', payload: 'pt'})
        expect(store.getState().language.value).toBe('pt')
    })
})

describe('getLanguage selector', () => {
    test('returns default language initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getLanguage(state)).toBe('en')
    })

    test('returns updated language after fetch', () => {
        const store = makeStore()
        store.dispatch({type: 'language/fetch/fulfilled', payload: 'ko'})
        const state = store.getState() as any
        expect(getLanguage(state)).toBe('ko')
    })
})
