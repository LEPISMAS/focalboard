// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {
    reducer as viewsReducer,
    updateViews,
    setCurrent,
    updateView,
    getViews,
    getSortedViews,
    getViewsByBoard,
    getView,
    getCurrentViewId,
    getCurrentView,
    getCurrentBoardViews,
    getCurrentViewGroupBy,
    getCurrentViewDisplayBy,
} from '../views'
import {reducer as boardsReducer} from '../boards'
import {createBoardView} from '../../blocks/boardView'
import {createBoard} from '../../blocks/board'

function makeStore() {
    return configureStore({
        reducer: {
            views: viewsReducer,
            boards: boardsReducer,
        },
    })
}

function makeView(id: string, parentId: string, title: string, deleteAt = 0) {
    const view = createBoardView()
    view.id = id
    view.boardId = parentId
    view.parentId = parentId
    view.title = title
    view.deleteAt = deleteAt
    return view
}

describe('views reducer', () => {
    test('initial state is empty', () => {
        const store = makeStore()
        expect(store.getState().views.views).toEqual({})
        expect(store.getState().views.current).toBe('')
    })

    test('setCurrent sets the current view id', () => {
        const store = makeStore()
        store.dispatch(setCurrent('view1'))
        expect(store.getState().views.current).toBe('view1')
    })

    test('updateViews adds views with deleteAt=0', () => {
        const store = makeStore()
        const view = makeView('v1', 'board1', 'My View')
        store.dispatch(updateViews([view]))
        expect(store.getState().views.views['v1']).toBeTruthy()
    })

    test('updateViews removes views with deleteAt != 0', () => {
        const store = makeStore()
        const view = makeView('v1', 'board1', 'My View')
        store.dispatch(updateViews([view]))
        const deleted = makeView('v1', 'board1', 'My View', 999)
        store.dispatch(updateViews([deleted]))
        expect(store.getState().views.views['v1']).toBeUndefined()
    })

    test('updateView sets a single view', () => {
        const store = makeStore()
        const view = makeView('v1', 'board1', 'View 1')
        store.dispatch(updateView(view))
        expect(store.getState().views.views['v1']).toBeTruthy()
    })

    test('updateViews smart update preserves reference equality for unchanged fields', () => {
        const store = makeStore()
        const view = makeView('v1', 'board1', 'View 1')
        store.dispatch(updateViews([view]))
        const updatedView = makeView('v1', 'board1', 'View 1 Updated')
        store.dispatch(updateViews([updatedView]))
        expect(store.getState().views.views['v1'].title).toBe('View 1 Updated')
    })

    test('extraReducer: initialReadOnlyLoad.fulfilled populates views', () => {
        const store = makeStore()
        const view = makeView('v1', 'board1', 'Readonly View')
        store.dispatch({
            type: 'initialReadOnlyLoad/fulfilled',
            payload: {board: {id: 'board1'}, blocks: [view]},
        })
        expect(store.getState().views.views['v1']).toBeTruthy()
    })

    test('extraReducer: initialReadOnlyLoad.fulfilled ignores non-view blocks', () => {
        const store = makeStore()
        const card = {id: 'card1', type: 'card', parentId: 'board1'}
        store.dispatch({
            type: 'initialReadOnlyLoad/fulfilled',
            payload: {board: {id: 'board1'}, blocks: [card]},
        })
        expect(Object.keys(store.getState().views.views)).toHaveLength(0)
    })

    test('extraReducer: loadBoardData.fulfilled populates views', () => {
        const store = makeStore()
        const view = makeView('v1', 'board1', 'View 1')
        store.dispatch({
            type: 'loadBoardData/fulfilled',
            payload: {blocks: [view]},
        })
        expect(store.getState().views.views['v1']).toBeTruthy()
    })

    test('extraReducer: loadBoardData.fulfilled clears previous views', () => {
        const store = makeStore()
        const view = makeView('v1', 'board1', 'View 1')
        store.dispatch(updateViews([view]))
        store.dispatch({type: 'loadBoardData/fulfilled', payload: {blocks: []}})
        expect(Object.keys(store.getState().views.views)).toHaveLength(0)
    })
})

describe('views selectors', () => {
    test('getViews returns empty map initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getViews(state)).toEqual({})
    })

    test('getSortedViews returns empty array initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getSortedViews(state)).toEqual([])
    })

    test('getSortedViews returns views sorted by title', () => {
        const store = makeStore()
        const v1 = makeView('v1', 'board1', 'Zebra')
        const v2 = makeView('v2', 'board1', 'Apple')
        store.dispatch(updateViews([v1, v2]))
        const state = store.getState() as any
        const sorted = getSortedViews(state)
        expect(sorted[0].title).toBe('Apple')
        expect(sorted[1].title).toBe('Zebra')
    })

    test('getViewsByBoard groups views by parentId', () => {
        const store = makeStore()
        const v1 = makeView('v1', 'board1', 'View 1')
        const v2 = makeView('v2', 'board1', 'View 2')
        const v3 = makeView('v3', 'board2', 'View 3')
        store.dispatch(updateViews([v1, v2, v3]))
        const state = store.getState() as any
        const byBoard = getViewsByBoard(state)
        expect(byBoard['board1']).toHaveLength(2)
        expect(byBoard['board2']).toHaveLength(1)
    })

    test('getView returns null for unknown viewId', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getView('unknown')(state)).toBeNull()
    })

    test('getView returns the view for known viewId', () => {
        const store = makeStore()
        const view = makeView('v1', 'board1', 'View 1')
        store.dispatch(updateViews([view]))
        const state = store.getState() as any
        expect(getView('v1')(state)).toBeTruthy()
    })

    test('getCurrentViewId returns empty string initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentViewId(state)).toBe('')
    })

    test('getCurrentView returns undefined when no views', () => {
        const store = makeStore()
        store.dispatch(setCurrent('v1'))
        const state = store.getState() as any
        expect(getCurrentView(state)).toBeUndefined()
    })

    test('getCurrentView returns the current view', () => {
        const store = makeStore()
        const view = makeView('v1', 'board1', 'View 1')
        store.dispatch(updateViews([view]))
        store.dispatch(setCurrent('v1'))
        const state = store.getState() as any
        expect(getCurrentView(state)?.id).toBe('v1')
    })

    test('getCurrentBoardViews returns empty array when no board is set', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentBoardViews(state)).toEqual([])
    })

    test('getCurrentBoardViews returns views for current board', () => {
        const store = makeStore()
        const board = createBoard()
        const v1 = makeView('v1', board.id, 'View 1')
        const v2 = makeView('v2', 'other-board', 'Other View')
        store.dispatch({
            type: 'boards/setCurrent',
            payload: board.id,
        })
        store.dispatch(updateViews([v1, v2]))
        const state = store.getState() as any
        const current = getCurrentBoardViews(state)
        expect(current).toHaveLength(1)
        expect(current[0].id).toBe('v1')
    })

    test('getCurrentViewGroupBy returns undefined when no current board or view', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentViewGroupBy(state)).toBeUndefined()
    })

    test('getCurrentViewDisplayBy returns undefined when no current board or view', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentViewDisplayBy(state)).toBeUndefined()
    })
})
