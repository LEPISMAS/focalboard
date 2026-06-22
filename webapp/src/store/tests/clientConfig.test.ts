// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {reducer, setClientConfig, getClientConfig} from '../clientConfig'
import {ShowUsername} from '../../utils'

// Mock octoClient to avoid real API calls
jest.mock('../../octoClient', () => ({
    default: {
        getClientConfig: jest.fn().mockResolvedValue({
            telemetry: true,
            telemetryid: 'test-id',
            enablePublicSharedBoards: true,
            teammateNameDisplay: 'username',
            featureFlags: {},
            maxFileSize: 1024,
        }),
    },
}))

function makeStore() {
    return configureStore({reducer: {clientConfig: reducer}})
}

const defaultConfig = {
    telemetry: false,
    telemetryid: '',
    enablePublicSharedBoards: false,
    teammateNameDisplay: ShowUsername,
    featureFlags: {},
    maxFileSize: 0,
}

describe('clientConfig reducer', () => {
    test('initial state has default ClientConfig values', () => {
        const store = makeStore()
        expect(store.getState().clientConfig.value).toEqual(defaultConfig)
    })

    test('setClientConfig updates state with provided config', () => {
        const store = makeStore()
        const newConfig = {
            telemetry: true,
            telemetryid: 'abc-123',
            enablePublicSharedBoards: true,
            teammateNameDisplay: 'full_name',
            featureFlags: {someFlag: true},
            maxFileSize: 2048,
        }
        store.dispatch(setClientConfig(newConfig as any))
        expect(store.getState().clientConfig.value).toEqual(newConfig)
    })

    test('setClientConfig replaces existing config', () => {
        const store = makeStore()
        const firstConfig = {...defaultConfig, telemetry: true, telemetryid: 'first'}
        const secondConfig = {...defaultConfig, telemetry: false, telemetryid: 'second'}
        store.dispatch(setClientConfig(firstConfig as any))
        store.dispatch(setClientConfig(secondConfig as any))
        expect(store.getState().clientConfig.value.telemetryid).toBe('second')
    })

    test('extraReducer: fetchClientConfig.fulfilled sets config from payload', () => {
        const store = makeStore()
        const payload = {
            telemetry: true,
            telemetryid: 'fetched-id',
            enablePublicSharedBoards: true,
            teammateNameDisplay: 'username',
            featureFlags: {},
            maxFileSize: 512,
        }
        store.dispatch({
            type: 'clientConfig/fetchClientConfig/fulfilled',
            payload,
        })
        expect(store.getState().clientConfig.value).toEqual(payload)
    })

    test('extraReducer: fetchClientConfig.fulfilled with null payload uses default config', () => {
        const store = makeStore()
        store.dispatch({
            type: 'clientConfig/fetchClientConfig/fulfilled',
            payload: null,
        })
        expect(store.getState().clientConfig.value).toEqual(defaultConfig)
    })
})

describe('getClientConfig selector', () => {
    test('returns the current client config', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getClientConfig(state)).toEqual(defaultConfig)
    })

    test('returns updated config after dispatch', () => {
        const store = makeStore()
        const newConfig = {...defaultConfig, telemetryid: 'xyz'}
        store.dispatch(setClientConfig(newConfig as any))
        const state = store.getState() as any
        expect(getClientConfig(state).telemetryid).toBe('xyz')
    })
})
