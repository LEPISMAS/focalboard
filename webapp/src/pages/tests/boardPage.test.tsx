// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {render, screen, fireEvent, waitFor, act} from '@testing-library/react'
import {Router, Route} from 'react-router-dom'
import {createMemoryHistory} from 'history'
import {mocked} from 'jest-mock'
import {wrapIntl} from '../../testUtils'

import BoardPage from '../boardPage/boardPage'
import {useWebsockets} from '../../hooks/websockets'
import octoClient from '../../octoClient'
import {UserSettings} from '../../userSettings'
import {initialLoad, initialReadOnlyLoad, loadBoardData} from '../../store/initialLoad'
import {fetchBoardMembers, setCurrent as setCurrentBoard} from '../../store/boards'
import {setCurrent as setCurrentView} from '../../store/views'

jest.mock('../../components/workspace', () => () => <div data-testid="workspace">Workspace</div>)
jest.mock('../../components/messages/versionMessage', () => () => <div data-testid="versionMessage">VersionMessage</div>)
jest.mock('../../components/confirmationDialogBox', () => (props: any) => (
    <div data-testid="confirmationDialog">
        ConfirmationDialog: {props.dialogBox.heading}
        <button onClick={props.dialogBox.onConfirm}>Confirm</button>
        <button onClick={props.dialogBox.onClose}>Close</button>
    </div>
))

jest.mock('../boardPage/teamToBoardAndViewRedirect', () => () => <div data-testid="teamRedirect" />)
jest.mock('../boardPage/backwardCompatibilityQueryParamsRedirect', () => () => <div data-testid="backwardRedirect" />)
jest.mock('../boardPage/setWindowTitleAndIcon', () => () => <div data-testid="setWindowTitle" />)
jest.mock('../boardPage/undoRedoHotKeys', () => () => <div data-testid="undoRedo" />)
jest.mock('../boardPage/websocketConnection', () => () => <div data-testid="websocketConnection" />)

jest.mock('../../hooks/websockets')
jest.mock('../../octoClient')

const mockDispatch = jest.fn((action) => {
    if (typeof action === 'function') {
        return action(mockDispatch, () => mockState, undefined)
    }
    return action
})

const mockState = {
    boards: {
        current: 'board-1',
    },
    views: {
        current: 'view-1',
    },
    users: {
        me: { id: 'user-1', permissions: ['manage_system'] } as any,
    },
    sidebar: {
        hiddenBoardIDs: [] as string[],
        categoryAttributes: [] as any[],
    },
}

jest.mock('../../store/hooks', () => ({
    useAppDispatch: () => mockDispatch,
    useAppSelector: jest.fn((selector) => selector(mockState)),
}))

jest.mock('../../store/boards', () => ({
    getCurrentBoardId: jest.fn(() => mockState.boards.current),
    fetchBoardMembers: jest.fn(() => ({ type: 'fetchBoardMembers' })),
    setCurrent: jest.fn(() => ({ type: 'setCurrentBoard' })),
    updateBoards: jest.fn(() => ({ type: 'updateBoards' })),
    updateMembersEnsuringBoardsAndUsers: jest.fn(() => ({ type: 'updateMembersEnsuringBoardsAndUsers' })),
    addMyBoardMemberships: jest.fn(() => ({ type: 'addMyBoardMemberships' })),
}))

jest.mock('../../store/views', () => ({
    getCurrentViewId: jest.fn(() => mockState.views.current),
    setCurrent: jest.fn(() => ({ type: 'setCurrentView' })),
    updateViews: jest.fn(() => ({ type: 'updateViews' })),
}))

jest.mock('../../store/users', () => ({
    getMe: jest.fn(() => mockState.users.me),
    followBlock: jest.fn(() => ({ type: 'followBlock' })),
    unfollowBlock: jest.fn(() => ({ type: 'unfollowBlock' })),
}))

jest.mock('../../store/sidebar', () => ({
    getHiddenBoardIDs: jest.fn(() => mockState.sidebar.hiddenBoardIDs),
    getCategoryOfBoard: jest.fn(() => () => null),
}))

jest.mock('../../store/initialLoad', () => ({
    initialLoad: jest.fn(() => ({ type: 'initialLoad' })),
    initialReadOnlyLoad: jest.fn(() => ({ type: 'initialReadOnlyLoad' })),
    loadBoardData: jest.fn(() => ({ type: 'loadBoardData', payload: { blocks: [{ id: 'block-1' }] } })),
}))

jest.mock('../../store/cards', () => ({ updateCards: jest.fn(() => ({ type: 'updateCards' })) }))
jest.mock('../../store/comments', () => ({ updateComments: jest.fn(() => ({ type: 'updateComments' })) }))
jest.mock('../../store/attachments', () => ({ updateAttachments: jest.fn(() => ({ type: 'updateAttachments' })) }))
jest.mock('../../store/contents', () => ({ updateContents: jest.fn(() => ({ type: 'updateContents' })) }))
jest.mock('../../store/teams', () => ({ setTeam: jest.fn(() => ({ type: 'setTeam' })) }))
jest.mock('../../store/globalError', () => ({ setGlobalError: jest.fn(() => ({ type: 'setGlobalError' })) }))

jest.mock('../../userSettings', () => {
    return {
        UserSettings: {
            mobileWarningClosed: false,
            lastTeamId: '',
            lastBoardId: {},
            lastViewId: {},
            setLastBoardID: jest.fn(),
            setLastViewId: jest.fn(),
        },
    }
})

const mockedUseWebsockets = mocked(useWebsockets, true)
const mockedOctoClient = mocked(octoClient, true)

describe('pages/boardPage/boardPage', () => {
    let history = createMemoryHistory()
    let mockWsClient: any
    let registeredCallbacks: Record<string, any> = {}

    beforeEach(() => {
        jest.clearAllMocks()
        history = createMemoryHistory()
        UserSettings.mobileWarningClosed = false
        UserSettings.lastTeamId = ''
        UserSettings.lastBoardId = {}
        UserSettings.lastViewId = {}

        mockState.boards.current = 'board-1'
        mockState.views.current = 'view-1'
        mockState.users.me = { id: 'user-1', permissions: ['manage_system'] } as any
        mockState.sidebar.hiddenBoardIDs = []

        registeredCallbacks = {}
        mockWsClient = {
            addOnChange: jest.fn((cb, type) => {
                registeredCallbacks[type] = cb
            }),
            removeOnChange: jest.fn(),
            addOnReconnect: jest.fn((cb) => {
                registeredCallbacks['reconnect'] = cb
            }),
            removeOnReconnect: jest.fn(),
            setOnFollowBlock: jest.fn((cb) => {
                registeredCallbacks['follow'] = cb
            }),
            setOnUnfollowBlock: jest.fn((cb) => {
                registeredCallbacks['unfollow'] = cb
            }),
        }

        mockedUseWebsockets.mockImplementation((teamId, callback) => {
            callback(mockWsClient)
        })
    })

    test('renders board page and subcomponents', () => {
        history.push('/team/team-1/board-1/view-1')

        render(
            wrapIntl(
                <Router history={history}>
                    <Route path="/team/:teamId/:boardId/:viewId">
                        <BoardPage />
                    </Route>
                </Router>
            )
        )

        expect(screen.getByTestId('workspace')).toBeDefined()
        expect(screen.getByTestId('versionMessage')).toBeDefined()
        expect(screen.getByTestId('teamRedirect')).toBeDefined()
        expect(screen.getByText('Mobile web support is currently in early beta.', {exact: false})).toBeDefined()
    })

    test('hides mobile warning when close is clicked', () => {
        history.push('/team/team-1/board-1/view-1')

        render(
            wrapIntl(
                <Router history={history}>
                    <Route path="/team/:teamId/:boardId/:viewId">
                        <BoardPage />
                    </Route>
                </Router>
            )
        )

        const closeBtn = screen.getByTitle('Close')
        fireEvent.click(closeBtn)

        expect(UserSettings.mobileWarningClosed).toBe(true)
        expect(screen.queryByText('Mobile web support is currently in early beta.', {exact: false})).toBeNull()
    })

    test('handles websockets incremental update registration and events', () => {
        history.push('/team/team-1/board-1/view-1')

        render(
            wrapIntl(
                <Router history={history}>
                    <Route path="/team/:teamId/:boardId/:viewId">
                        <BoardPage />
                    </Route>
                </Router>
            )
        )

        expect(registeredCallbacks['block']).toBeDefined()
        expect(registeredCallbacks['board']).toBeDefined()
        expect(registeredCallbacks['boardMembers']).toBeDefined()
        expect(registeredCallbacks['reconnect']).toBeDefined()
        expect(registeredCallbacks['follow']).toBeDefined()
        expect(registeredCallbacks['unfollow']).toBeDefined()

        // Trigger incremental updates
        act(() => {
            registeredCallbacks['block'](mockWsClient, [{ id: 'block-1', type: 'card' }])
            registeredCallbacks['board'](mockWsClient, [{ id: 'board-1', teamId: 'team-1' }])
            registeredCallbacks['boardMembers'](mockWsClient, [{ userId: 'user-1', boardId: 'board-1' }])
            registeredCallbacks['reconnect']()
            registeredCallbacks['follow'](mockWsClient, { subscriberId: 'user-1', blockId: 'block-1' })
            registeredCallbacks['unfollow'](mockWsClient, { subscriberId: 'user-1', blockId: 'block-1' })
        })

        // Verify dispatches happened
        expect(mockDispatch).toHaveBeenCalled()
    })

    test('handles loadOrJoinBoard when blocks list is empty', async () => {
        mocked(loadBoardData as any).mockImplementation(() => ({
            type: 'loadBoardData',
            payload: { blocks: [] },
        }))

        mockedOctoClient.joinBoard.mockImplementation(async (boardId, allowAdmin) => {
            if (allowAdmin) {
                return { id: 'member-1' } as any
            }
            return null
        })

        history.push('/team/team-1/board-1/view-1')

        render(
            wrapIntl(
                <Router history={history}>
                    <Route path="/team/:teamId/:boardId/:viewId">
                        <BoardPage />
                    </Route>
                </Router>
            )
        )

        // Wait for join board process
        await waitFor(() => {
            expect(screen.getByTestId('confirmationDialog')).toBeDefined()
        })

        // Confirm join
        fireEvent.click(screen.getByText('Confirm'))

        await waitFor(() => {
            expect(mockedOctoClient.joinBoard).toHaveBeenCalledWith('board-1', true)
        })
    })

    test('renders board page in readonly mode', () => {
        mockState.users.me = null as any

        history.push('/team/team-1/board-1/view-1')

        render(
            wrapIntl(
                <Router history={history}>
                    <Route path="/team/:teamId/:boardId/:viewId">
                        <BoardPage readonly={true} />
                    </Route>
                </Router>
            )
        )

        expect(initialReadOnlyLoad).toHaveBeenCalled()
    })
})
