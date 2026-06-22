// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {reducer, getGlobalTemplates} from '../globalTemplates'

// Mock octoClient to avoid real API calls
jest.mock('../../octoClient', () => ({
    default: {
        getTeamTemplates: jest.fn().mockResolvedValue([]),
    },
}))

function makeStore() {
    return configureStore({reducer: {globalTemplates: reducer}})
}

const createMockBoard = (id: string, title: string) => ({
    id,
    teamId: 'team1',
    title,
    description: '',
    icon: '',
    showDescription: false,
    isTemplate: true,
    templateVersion: 1,
    properties: {},
    cardProperties: [],
    createdBy: '',
    modifiedBy: '',
    type: 'P' as any,
    minimumRole: '' as any,
    createAt: 1000,
    updateAt: 1000,
    deleteAt: 0,
})

describe('globalTemplates reducer', () => {
    test('initial state has empty templates array', () => {
        const store = makeStore()
        expect(store.getState().globalTemplates.value).toEqual([])
    })

    test('extraReducer: fetchGlobalTemplates.fulfilled sets templates from payload', () => {
        const store = makeStore()
        const templates = [
            createMockBoard('tmpl1', 'Template Alpha'),
            createMockBoard('tmpl2', 'Template Beta'),
        ]
        store.dispatch({
            type: 'globalTemplates/fetch/fulfilled',
            payload: templates,
        })
        expect(store.getState().globalTemplates.value).toHaveLength(2)
        expect(store.getState().globalTemplates.value[0].id).toBe('tmpl1')
    })

    test('extraReducer: fetchGlobalTemplates.fulfilled with null payload sets empty array', () => {
        const store = makeStore()
        store.dispatch({
            type: 'globalTemplates/fetch/fulfilled',
            payload: null,
        })
        expect(store.getState().globalTemplates.value).toEqual([])
    })

    test('extraReducer: fetchGlobalTemplates.fulfilled with empty array', () => {
        const store = makeStore()
        store.dispatch({
            type: 'globalTemplates/fetch/fulfilled',
            payload: [],
        })
        expect(store.getState().globalTemplates.value).toEqual([])
    })

    test('extraReducer: sets templates replacing previous ones', () => {
        const store = makeStore()
        const firstTemplates = [createMockBoard('tmpl1', 'Template 1')]
        store.dispatch({type: 'globalTemplates/fetch/fulfilled', payload: firstTemplates})

        const secondTemplates = [createMockBoard('tmpl2', 'Template 2'), createMockBoard('tmpl3', 'Template 3')]
        store.dispatch({type: 'globalTemplates/fetch/fulfilled', payload: secondTemplates})

        expect(store.getState().globalTemplates.value).toHaveLength(2)
        expect(store.getState().globalTemplates.value[0].id).toBe('tmpl2')
    })
})

describe('getGlobalTemplates selector', () => {
    test('returns empty array initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getGlobalTemplates(state)).toEqual([])
    })

    test('returns templates after dispatch', () => {
        const store = makeStore()
        const templates = [createMockBoard('tmpl1', 'Template A')]
        store.dispatch({type: 'globalTemplates/fetch/fulfilled', payload: templates})
        const state = store.getState() as any
        expect(getGlobalTemplates(state)).toHaveLength(1)
        expect(getGlobalTemplates(state)[0].title).toBe('Template A')
    })
})
