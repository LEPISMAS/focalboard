// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {render} from '@testing-library/react'
import {Router, Route} from 'react-router-dom'
import {createMemoryHistory} from 'history'
import {Provider as ReduxProvider} from 'react-redux'
import configureStore from 'redux-mock-store'
import thunk from 'redux-thunk'
import {mocked} from 'jest-mock'

import {getCurrentBoardId, getBoards} from '../../store/boards'
import {getCurrentBoardViews, setCurrent as setCurrentView} from '../../store/views'
import {getSidebarCategories} from '../../store/sidebar'
import {UserSettings} from '../../userSettings'
import TeamToBoardAndViewRedirect from '../boardPage/teamToBoardAndViewRedirect'

jest.mock('../../store/boards')
jest.mock('../../store/views', () => {
    const original = jest.requireActual('../../store/views')
    return {
        ...original,
        setCurrent: jest.fn(() => ({type: 'setCurrentView'})),
        getCurrentBoardViews: jest.fn(),
    }
})
jest.mock('../../store/sidebar')
jest.mock('../../userSettings', () => {
    return {
        UserSettings: {
            lastTeamId: 'last-team-id',
            lastBoardId: {},
            lastViewId: {},
            setLastViewId: jest.fn(),
        },
    }
})

const mockedGetCurrentBoardId = mocked(getCurrentBoardId, true)
const mockedGetBoards = mocked(getBoards, true)
const mockedGetCurrentBoardViews = mocked(getCurrentBoardViews, true)
const mockedGetSidebarCategories = mocked(getSidebarCategories, true)

describe('pages/boardPage/teamToBoardAndViewRedirect', () => {
    const mockStore = configureStore([thunk])
    const store = mockStore({})
    let history = createMemoryHistory()

    beforeEach(() => {
        jest.clearAllMocks()
        history = createMemoryHistory()
        Object.defineProperty(UserSettings, 'lastBoardId', {value: {}, configurable: true})
        Object.defineProperty(UserSettings, 'lastViewId', {value: {}, configurable: true})
    })

    test('redirects to last visited board if boardId is missing in URL', () => {
        history.push('/team/team_id_1')
        mockedGetCurrentBoardId.mockReturnValue('')
        mockedGetCurrentBoardViews.mockReturnValue([])
        mockedGetSidebarCategories.mockReturnValue([])
        mockedGetBoards.mockReturnValue({})

        Object.defineProperty(UserSettings, 'lastBoardId', {value: {team_id_1: 'last-board-id'}, configurable: true})
        history.replace = jest.fn()

        render(
            <ReduxProvider store={store}>
                <Router history={history}>
                    <Route path='/team/:teamId/:boardId?/:viewId?'>
                        <TeamToBoardAndViewRedirect/>
                    </Route>
                </Router>
            </ReduxProvider>,
        )

        expect(history.replace).toHaveBeenCalledWith('/team/team_id_1/last-board-id')
    })

    test('redirects to first board in categories if last visited board is missing', () => {
        history.push('/team/team_id_1')
        mockedGetCurrentBoardId.mockReturnValue('')
        mockedGetCurrentBoardViews.mockReturnValue([])
        mockedGetBoards.mockReturnValue({
            'board-1': {id: 'board-1', title: 'Board 1'} as any,
        })
        mockedGetSidebarCategories.mockReturnValue([
            {
                boardMetadata: [
                    {boardID: 'board-1', hidden: false},
                ],
            } as any,
        ])

        history.replace = jest.fn()

        render(
            <ReduxProvider store={store}>
                <Router history={history}>
                    <Route path='/team/:teamId/:boardId?/:viewId?'>
                        <TeamToBoardAndViewRedirect/>
                    </Route>
                </Router>
            </ReduxProvider>,
        )

        expect(history.replace).toHaveBeenCalledWith('/team/team_id_1/board-1')
    })

    test('redirects to last visited view if viewId is missing in URL', () => {
        history.push('/team/team_id_1/board-1')
        mockedGetCurrentBoardId.mockReturnValue('board-1')
        mockedGetCurrentBoardViews.mockReturnValue([
            {id: 'view-1', title: 'View 1'} as any,
            {id: 'view-2', title: 'View 2'} as any,
        ])
        mockedGetSidebarCategories.mockReturnValue([])
        mockedGetBoards.mockReturnValue({})

        Object.defineProperty(UserSettings, 'lastViewId', {value: {'board-1': 'view-2'}, configurable: true})
        history.replace = jest.fn()

        render(
            <ReduxProvider store={store}>
                <Router history={history}>
                    <Route path='/team/:teamId/:boardId?/:viewId?'>
                        <TeamToBoardAndViewRedirect/>
                    </Route>
                </Router>
            </ReduxProvider>,
        )

        expect(setCurrentView).toHaveBeenCalledWith('view-2')
        expect(history.replace).toHaveBeenCalledWith('/team/team_id_1/board-1/view-2')
    })

    test('redirects to first view in list if last visited view is missing', () => {
        history.push('/team/team_id_1/board-1')
        mockedGetCurrentBoardId.mockReturnValue('board-1')
        mockedGetCurrentBoardViews.mockReturnValue([
            {id: 'view-1', title: 'View 1'} as any,
        ])
        mockedGetSidebarCategories.mockReturnValue([])
        mockedGetBoards.mockReturnValue({})

        history.replace = jest.fn()

        render(
            <ReduxProvider store={store}>
                <Router history={history}>
                    <Route path='/team/:teamId/:boardId?/:viewId?'>
                        <TeamToBoardAndViewRedirect/>
                    </Route>
                </Router>
            </ReduxProvider>,
        )

        expect(setCurrentView).toHaveBeenCalledWith('view-1')
        expect(history.replace).toHaveBeenCalledWith('/team/team_id_1/board-1/view-1')
    })
})
