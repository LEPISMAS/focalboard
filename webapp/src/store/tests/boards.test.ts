// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {
    reducer as boardsReducer,
    fetchBoardMembers,
    updateMembersEnsuringBoardsAndUsers,
    updateBoards,
    setCurrent,
    setLinkToChannel,
    updateMembers,
    addMyBoardMemberships,
    updateMembersHandler,
    getBoards,
    getMySortedBoards,
    getTemplates,
    getSortedTemplates,
    getBoard,
    isLoadingBoard,
    getCurrentBoardId,
    getCurrentBoard,
    getCurrentBoardMembers,
    getMyBoardMembership,
    getCurrentLinkToChannel,
} from '../boards'
import {reducer as usersReducer} from '../users'
import {createBoard} from '../../blocks/board'

jest.mock('../../octoClient', () => ({
    default: {
        getBoardMembers: jest.fn().mockResolvedValue([]),
        getTeamUsersList: jest.fn().mockResolvedValue([]),
        getBoard: jest.fn().mockResolvedValue(null),
    },
}))

jest.mock('../../utils', () => ({
    Utils: {
        log: jest.fn(),
        logError: jest.fn(),
        createGuid: jest.fn(() => 'mock-guid'),
    },
    IDType: {
        BlockID: 'block',
    },
}))

function makeStore() {
    return configureStore({
        reducer: {
            boards: boardsReducer,
            users: usersReducer,
        },
    })
}

function makeBoard(id: string, title: string, isTemplate = false, deleteAt = 0) {
    const board = createBoard()
    board.id = id
    board.title = title
    board.isTemplate = isTemplate
    board.deleteAt = deleteAt
    return board
}

function makeMember(boardId: string, userId: string, schemeAdmin = false, schemeEditor = true) {
    return {
        boardId,
        userId,
        roles: '',
        minimumRole: '' as any,
        schemeAdmin,
        schemeEditor,
        schemeCommenter: false,
        schemeViewer: true,
        synthetic: false,
    }
}

describe('boards reducer', () => {
    test('initial state is empty', () => {
        const store = makeStore()
        expect(store.getState().boards.boards).toEqual({})
        expect(store.getState().boards.templates).toEqual({})
        expect(store.getState().boards.membersInBoards).toEqual({})
        expect(store.getState().boards.myBoardMemberships).toEqual({})
        expect(store.getState().boards.loadingBoard).toBe(false)
    })

    test('setCurrent sets the current board id', () => {
        const store = makeStore()
        store.dispatch(setCurrent('board1'))
        expect(store.getState().boards.current).toBe('board1')
    })

    test('setLinkToChannel sets the link to channel', () => {
        const store = makeStore()
        store.dispatch(setLinkToChannel('channel1'))
        expect(store.getState().boards.linkToChannel).toBe('channel1')
    })

    test('updateBoards adds non-template board', () => {
        const store = makeStore()
        const board = makeBoard('b1', 'Board 1')
        store.dispatch(updateBoards([board]))
        expect(store.getState().boards.boards['b1']).toBeTruthy()
    })

    test('updateBoards adds template board to templates', () => {
        const store = makeStore()
        const tmpl = makeBoard('tmpl1', 'Template 1', true)
        store.dispatch(updateBoards([tmpl]))
        expect(store.getState().boards.templates['tmpl1']).toBeTruthy()
        expect(store.getState().boards.boards['tmpl1']).toBeUndefined()
    })

    test('updateBoards removes board when deleteAt != 0', () => {
        const store = makeStore()
        const board = makeBoard('b1', 'Board 1')
        store.dispatch(updateBoards([board]))
        const deleted = makeBoard('b1', 'Board 1', false, 999)
        store.dispatch(updateBoards([deleted]))
        expect(store.getState().boards.boards['b1']).toBeUndefined()
    })

    test('addMyBoardMemberships adds memberships', () => {
        const store = makeStore()
        const member = makeMember('b1', 'u1', true)
        store.dispatch(addMyBoardMemberships([member]))
        expect(store.getState().boards.myBoardMemberships['b1']).toBeTruthy()
    })

    test('addMyBoardMemberships removes when no scheme perms', () => {
        const store = makeStore()
        const member = makeMember('b1', 'u1')
        store.dispatch(addMyBoardMemberships([member]))
        const removedMember = {
            ...member,
            schemeAdmin: false,
            schemeEditor: false,
            schemeViewer: false,
            schemeCommenter: false,
        }
        store.dispatch(addMyBoardMemberships([removedMember]))
        expect(store.getState().boards.myBoardMemberships['b1']).toBeUndefined()
    })

    test('updateMembers updates member in board', () => {
    const state: any = {
        membersInBoards: {b1: {}},
        myBoardMemberships: {},
        boards: {},
        templates: {},
        loadingBoard: false,
        current: '',
        linkToChannel: '',
    }

    const member = makeMember('b1', 'u1')

    updateMembersHandler(state, {
        payload: [member],
    } as any)

    expect(state.membersInBoards.b1.u1).toBeTruthy()
})

    test('updateMembers with empty payload does nothing', () => {
        const store = makeStore()
        store.dispatch(updateMembers([]))
        expect(store.getState().boards.membersInBoards).toEqual({})
    })

    test('updateMembers removes member with no perms', () => {
    const state: any = {
        membersInBoards: {b1: {}},
        myBoardMemberships: {},
        boards: {},
        templates: {},
        loadingBoard: false,
        current: '',
        linkToChannel: '',
    }

    const member = makeMember('b1', 'u1')

    updateMembersHandler(state, {
        payload: [member],
    } as any)

    const noPermMember = {
        ...member,
        schemeAdmin: false,
        schemeEditor: false,
        schemeViewer: false,
        schemeCommenter: false,
    }

    updateMembersHandler(state, {
        payload: [noPermMember],
    } as any)

    expect(state.membersInBoards.b1.u1).toBeUndefined()
})

    test('extraReducer: loadBoardData.pending sets loadingBoard=true', () => {
        const store = makeStore()
        store.dispatch({type: 'loadBoardData/pending'})
        expect(store.getState().boards.loadingBoard).toBe(true)
    })

    test('extraReducer: loadBoardData.fulfilled sets loadingBoard=false', () => {
        const store = makeStore()
        store.dispatch({type: 'loadBoardData/pending'})
        store.dispatch({type: 'loadBoardData/fulfilled', payload: {blocks: []}})
        expect(store.getState().boards.loadingBoard).toBe(false)
    })

    test('extraReducer: loadBoardData.rejected sets loadingBoard=false', () => {
        const store = makeStore()
        store.dispatch({type: 'loadBoardData/pending'})
        store.dispatch({type: 'loadBoardData/rejected', error: {}})
        expect(store.getState().boards.loadingBoard).toBe(false)
    })

    test('extraReducer: initialReadOnlyLoad.fulfilled sets board', () => {
        const store = makeStore()
        const board = makeBoard('b1', 'Board 1')
        store.dispatch({
            type: 'initialReadOnlyLoad/fulfilled',
            payload: {board, blocks: []},
        })
        expect(store.getState().boards.boards['b1']).toBeTruthy()
    })

    test('extraReducer: initialReadOnlyLoad.fulfilled with template board', () => {
        const store = makeStore()
        const tmpl = makeBoard('tmpl1', 'Template', true)
        store.dispatch({
            type: 'initialReadOnlyLoad/fulfilled',
            payload: {board: tmpl, blocks: []},
        })
        expect(store.getState().boards.templates['tmpl1']).toBeTruthy()
    })

    test('extraReducer: initialReadOnlyLoad.fulfilled with no board', () => {
        const store = makeStore()
        store.dispatch({
            type: 'initialReadOnlyLoad/fulfilled',
            payload: {board: null, blocks: []},
        })
        expect(Object.keys(store.getState().boards.boards)).toHaveLength(0)
    })

    test('extraReducer: initialLoad.fulfilled populates boards and templates', () => {
        const store = makeStore()
        const b1 = makeBoard('b1', 'Board 1')
        const tmpl1 = makeBoard('tmpl1', 'Template 1', true)
        const member = makeMember('b1', 'u1')
        store.dispatch({
            type: 'initialLoad/fulfilled',
            payload: {
                team: {id: 't1'},
                teams: [],
                boards: [b1],
                boardsMemberships: [member],
                boardTemplates: [tmpl1],
                limits: null,
                myConfig: [],
            },
        })
        expect(store.getState().boards.boards['b1']).toBeTruthy()
        expect(store.getState().boards.templates['tmpl1']).toBeTruthy()
        expect(store.getState().boards.myBoardMemberships['b1']).toBeTruthy()
    })

    test('extraReducer: loadBoards.fulfilled populates boards', () => {
        const store = makeStore()
        const b1 = makeBoard('b1', 'Board 1')
        store.dispatch({
            type: 'loadBoards/fulfilled',
            payload: {boards: [b1]},
        })
        expect(store.getState().boards.boards['b1']).toBeTruthy()
    })

    test('extraReducer: loadMyBoardsMemberships.fulfilled populates myBoardMemberships', () => {
        const store = makeStore()
        const member = makeMember('b1', 'u1')
        store.dispatch({
            type: 'loadMyBoardsMemberships/fulfilled',
            payload: {boardsMemberships: [member]},
        })
        expect(store.getState().boards.myBoardMemberships['b1']).toBeTruthy()
    })

    test('extraReducer: boardMembers/fetch/fulfilled populates membersInBoards', () => {
        const store = makeStore()
        const member = makeMember('b1', 'u1')
        store.dispatch({
            type: 'boardMembers/fetch/fulfilled',
            payload: [member],
        })
        expect(store.getState().boards.membersInBoards['b1']['u1']).toBeTruthy()
    })

    test('extraReducer: boardMembers/fetch/fulfilled with empty payload does nothing', () => {
        const store = makeStore()
        store.dispatch({
            type: 'boardMembers/fetch/fulfilled',
            payload: [],
        })
        expect(store.getState().boards.membersInBoards).toEqual({})
    })
})

describe('boards selectors', () => {
    test('getBoards returns empty object initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getBoards(state)).toEqual({})
    })

    test('getMySortedBoards returns empty array initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getMySortedBoards(state)).toEqual([])
    })

    test('getMySortedBoards returns boards with membership sorted by title', () => {
        const store = makeStore()
        const b1 = makeBoard('b1', 'Zebra Board')
        const b2 = makeBoard('b2', 'Apple Board')
        const m1 = makeMember('b1', 'u1')
        const m2 = makeMember('b2', 'u1')
        store.dispatch(updateBoards([b1, b2]))
        store.dispatch(addMyBoardMemberships([m1, m2]))
        const state = store.getState() as any
        const sorted = getMySortedBoards(state)
        expect(sorted[0].title).toBe('Apple Board')
        expect(sorted[1].title).toBe('Zebra Board')
    })

    test('getTemplates returns empty object initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getTemplates(state)).toEqual({})
    })

    test('getSortedTemplates returns templates sorted by title', () => {
        const store = makeStore()
        const t1 = makeBoard('t1', 'Zeta Template', true)
        const t2 = makeBoard('t2', 'Alpha Template', true)
        store.dispatch(updateBoards([t1, t2]))
        const state = store.getState() as any
        const sorted = getSortedTemplates(state)
        expect(sorted[0].title).toBe('Alpha Template')
    })

    test('getBoard returns null for unknown boardId', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getBoard('unknown')(state)).toBeNull()
    })

    test('getBoard returns board from boards', () => {
        const store = makeStore()
        const board = makeBoard('b1', 'Board 1')
        store.dispatch(updateBoards([board]))
        const state = store.getState() as any
        expect(getBoard('b1')(state)?.id).toBe('b1')
    })

    test('getBoard returns board from templates', () => {
        const store = makeStore()
        const tmpl = makeBoard('tmpl1', 'Template', true)
        store.dispatch(updateBoards([tmpl]))
        const state = store.getState() as any
        expect(getBoard('tmpl1')(state)?.id).toBe('tmpl1')
    })

    test('isLoadingBoard returns false initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(isLoadingBoard(state)).toBe(false)
    })

    test('getCurrentBoardId returns empty string initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentBoardId(state)).toBe('')
    })

    test('getCurrentBoard returns undefined when no board selected', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentBoard(state)).toBeUndefined()
    })

    test('getCurrentBoard returns selected board', () => {
        const store = makeStore()
        const board = makeBoard('b1', 'Board 1')
        store.dispatch(updateBoards([board]))
        store.dispatch(setCurrent('b1'))
        const state = store.getState() as any
        expect(getCurrentBoard(state)?.id).toBe('b1')
    })

    test('getCurrentBoardMembers returns empty object when no board selected', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentBoardMembers(state)).toEqual({})
    })

    test('getCurrentBoardMembers returns members for current board', () => {
    const state: any = {
        boards: {
            current: 'b1',
            membersInBoards: {
                b1: {
                    u1: makeMember('b1', 'u1'),
                },
            },
            myBoardMemberships: {},
            boards: {},
            templates: {},
            loadingBoard: false,
            linkToChannel: '',
        },
        users: {},
    }

    expect(getCurrentBoardMembers(state).u1).toBeTruthy()
})

    test('getMyBoardMembership returns null for unknown board', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getMyBoardMembership('unknown')(state)).toBeNull()
    })

    test('getMyBoardMembership returns membership for known board', () => {
        const store = makeStore()
        const member = makeMember('b1', 'u1')
        store.dispatch(addMyBoardMemberships([member]))
        const state = store.getState() as any
        expect(getMyBoardMembership('b1')(state)).toBeTruthy()
    })

    test('getCurrentLinkToChannel returns empty string initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentLinkToChannel(state)).toBe('')
    })
})

describe('updateMembersHandler', () => {
    test('handles empty payload', () => {
        const state: any = {membersInBoards: {}, myBoardMemberships: {}, boards: {}, templates: {}, loadingBoard: false, current: '', linkToChannel: ''}
        const action: any = {payload: []}
        updateMembersHandler(state, action)
        expect(state.membersInBoards).toEqual({})
    })

    test('updates myBoardMemberships when member found', () => {
        const state: any = {
            membersInBoards: {'b1': {}},
            myBoardMemberships: {'b1': {boardId: 'b1', userId: 'u1', schemeAdmin: false, schemeEditor: true, schemeViewer: true, schemeCommenter: false}},
            boards: {},
            templates: {},
            loadingBoard: false,
            current: '',
            linkToChannel: '',
        }
        const member = makeMember('b1', 'u1', true)
        const action: any = {payload: [member]}
        updateMembersHandler(state, action)
        expect(state.myBoardMemberships['b1'].schemeAdmin).toBe(true)
    })

    test('removes myBoardMembership when no perms', () => {
        const state: any = {
            membersInBoards: {'b1': {}},
            myBoardMemberships: {'b1': {boardId: 'b1', userId: 'u1', schemeAdmin: false, schemeEditor: true, schemeViewer: true, schemeCommenter: false}},
            boards: {},
            templates: {},
            loadingBoard: false,
            current: '',
            linkToChannel: '',
        }
        const noPermMember = {
            ...makeMember('b1', 'u1'),
            schemeAdmin: false,
            schemeEditor: false,
            schemeViewer: false,
            schemeCommenter: false,
        }
        const action: any = {payload: [noPermMember]}
        updateMembersHandler(state, action)
        expect(state.myBoardMemberships['b1']).toBeUndefined()
    })
})


describe('boards async thunk coverage', () => {

    test('updateMembersEnsuringBoardsAndUsers marks board deleted when current user membership has no permissions', async () => {
        const member = makeMember('board-removed', 'me-user', false, false)
        member.schemeViewer = false
        member.schemeCommenter = false

        const dispatch = jest.fn()
        const getState = jest.fn(() => ({
            users: {
                me: {id: 'me-user'},
                boardUsers: {},
            },
            boards: {
                boards: {},
            },
        }))

        const result = await updateMembersEnsuringBoardsAndUsers([member])(dispatch, getState, undefined)

        expect(result.type).toBe('updateMembersEnsuringBoardsAndUsers/fulfilled')
        expect(result.payload).toEqual([member])
        expect(dispatch).toHaveBeenCalled()
    })

    test('updateMembersEnsuringBoardsAndUsers removes deleted board user', async () => {
        const member = makeMember('board1', 'user-deleted', false, false)
        member.schemeViewer = false
        member.schemeCommenter = false

        const dispatch = jest.fn()
        const getState = jest.fn(() => ({
            users: {
                me: null,
                boardUsers: {},
            },
            boards: {
                boards: {},
            },
        }))

        const result = await updateMembersEnsuringBoardsAndUsers([member])(dispatch, getState, undefined)

        expect(result.type).toBe('updateMembersEnsuringBoardsAndUsers/fulfilled')
        expect(result.payload).toEqual([member])
        expect(dispatch).toHaveBeenCalled()
    })
})

