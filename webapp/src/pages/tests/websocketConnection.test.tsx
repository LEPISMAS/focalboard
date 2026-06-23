// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {render, screen, act} from '@testing-library/react'
import {Provider as ReduxProvider} from 'react-redux'
import configureStore from 'redux-mock-store'
import thunk from 'redux-thunk'
import {wrapIntl} from '../../testUtils'
import wsClient from '../../wsclient'
import WebsocketConnection from '../boardPage/websocketConnection'

jest.mock('../../wsclient', () => {
    return {
        __esModule: true,
        default: {
            addOnStateChange: jest.fn(),
            removeOnStateChange: jest.fn(),
        },
        WSClient: jest.fn(),
    }
})

describe('pages/boardPage/websocketConnection', () => {
    const mockStore = configureStore([thunk])
    const store = mockStore({
        users: {
            me: {
                id: 'my-user-id',
            },
        },
    })

    beforeEach(() => {
        jest.clearAllMocks()
        jest.useFakeTimers()
    })

    afterEach(() => {
        jest.useRealTimers()
    })

    test('renders null by default and registers listener', () => {
        const {container} = render(
            <ReduxProvider store={store}>
                {wrapIntl(<WebsocketConnection />)}
            </ReduxProvider>
        )

        expect(container.firstChild).toBeNull()
        expect(wsClient.addOnStateChange).toHaveBeenCalled()
    })

    test('shows banner after 5 seconds on close state change', () => {
        let stateChangeCallback: any = null
        ;(wsClient.addOnStateChange as jest.Mock).mockImplementation((cb) => {
            stateChangeCallback = cb
        })

        render(
            <ReduxProvider store={store}>
                {wrapIntl(<WebsocketConnection />)}
            </ReduxProvider>
        )

        expect(stateChangeCallback).not.toBeNull()

        // Trigger close state change
        act(() => {
            stateChangeCallback(wsClient, 'close')
        })

        // Before 5 seconds, should be null
        expect(screen.queryByText(/Websocket connection closed/)).toBeNull()

        // Advance timers by 5000ms
        act(() => {
            jest.advanceTimersByTime(5000)
        })

        expect(screen.getByText(/Websocket connection closed/)).toBeDefined()
    })

    test('hides banner if open state is received before/after close', () => {
        let stateChangeCallback: any = null
        ;(wsClient.addOnStateChange as jest.Mock).mockImplementation((cb) => {
            stateChangeCallback = cb
        })

        render(
            <ReduxProvider store={store}>
                {wrapIntl(<WebsocketConnection />)}
            </ReduxProvider>
        )

        // Close connection
        act(() => {
            stateChangeCallback(wsClient, 'close')
        })

        // Open connection before timeout fires
        act(() => {
            stateChangeCallback(wsClient, 'open')
        })

        act(() => {
            jest.advanceTimersByTime(5000)
        })

        expect(screen.queryByText(/Websocket connection closed/)).toBeNull()
    })

    test('unregisters listener on unmount', () => {
        const {unmount} = render(
            <ReduxProvider store={store}>
                {wrapIntl(<WebsocketConnection />)}
            </ReduxProvider>
        )

        unmount()
        expect(wsClient.removeOnStateChange).toHaveBeenCalled()
    })
})
