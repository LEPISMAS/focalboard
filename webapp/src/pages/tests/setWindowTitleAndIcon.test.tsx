// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {render} from '@testing-library/react'
import {Provider as ReduxProvider} from 'react-redux'
import configureStore from 'redux-mock-store'
import thunk from 'redux-thunk'
import {mocked} from 'jest-mock'

import {getCurrentBoard} from '../../store/boards'
import {getCurrentView} from '../../store/views'
import {Utils} from '../../utils'
import SetWindowTitleAndIcon from '../boardPage/setWindowTitleAndIcon'

jest.mock('../../store/boards')
jest.mock('../../store/views')
jest.mock('../../utils')

const mockedGetCurrentBoard = mocked(getCurrentBoard, true)
const mockedGetCurrentView = mocked(getCurrentView, true)
const mockedUtils = mocked(Utils, true)

describe('pages/boardPage/setWindowTitleAndIcon', () => {
    const mockStore = configureStore([thunk])
    const store = mockStore({})

    beforeEach(() => {
        jest.resetAllMocks()
        document.title = ''
    })

    test('sets Focalboard as title and favicon to undefined if no board is selected', () => {
        mockedGetCurrentBoard.mockReturnValue(null as any)
        mockedGetCurrentView.mockReturnValue(null as any)

        render(
            <ReduxProvider store={store}>
                <SetWindowTitleAndIcon/>
            </ReduxProvider>,
        )

        expect(document.title).toBe('Focalboard')
        expect(mockedUtils.setFavicon).toHaveBeenCalledWith(undefined)
    })

    test('sets board title and icon', () => {
        mockedGetCurrentBoard.mockReturnValue({
            title: 'My Board',
            icon: '🌟',
        } as any)
        mockedGetCurrentView.mockReturnValue(null as any)

        render(
            <ReduxProvider store={store}>
                <SetWindowTitleAndIcon/>
            </ReduxProvider>,
        )

        expect(document.title).toBe('My Board')
        expect(mockedUtils.setFavicon).toHaveBeenCalledWith('🌟')
    })

    test('sets board and view title in document title', () => {
        mockedGetCurrentBoard.mockReturnValue({
            title: 'My Board',
            icon: '🌟',
        } as any)
        mockedGetCurrentView.mockReturnValue({
            title: 'Board View',
        } as any)

        render(
            <ReduxProvider store={store}>
                <SetWindowTitleAndIcon/>
            </ReduxProvider>,
        )

        expect(document.title).toBe('My Board | Board View')
        expect(mockedUtils.setFavicon).toHaveBeenCalledWith('🌟')
    })
})
