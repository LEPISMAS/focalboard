// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {reducer, updateComments, getCardComments, getLastCardComment, getLastCommentByCard} from '../comments'
import {createCommentBlock} from '../../blocks/commentBlock'

function makeStore() {
    return configureStore({reducer: {comments: reducer}})
}

function makeComment(id: string, parentId: string, deleteAt = 0, createAt = 1000) {
    const comment = createCommentBlock()
    comment.id = id
    comment.parentId = parentId
    comment.deleteAt = deleteAt
    comment.createAt = createAt
    return comment
}

describe('comments reducer', () => {
    test('initial state is empty', () => {
        const store = makeStore()
        expect(store.getState().comments.comments).toEqual({})
        expect(store.getState().comments.commentsByCard).toEqual({})
    })

    test('updateComments adds new comment to new parent', () => {
        const store = makeStore()
        const comment = makeComment('c1', 'card1')
        store.dispatch(updateComments([comment]))
        expect(store.getState().comments.comments.c1).toBeTruthy()
        expect(store.getState().comments.commentsByCard.card1).toHaveLength(1)
    })

    test('updateComments adds to existing parent', () => {
        const store = makeStore()
        const c1 = makeComment('c1', 'card1', 0, 1000)
        const c2 = makeComment('c2', 'card1', 0, 1001)
        store.dispatch(updateComments([c1]))
        store.dispatch(updateComments([c2]))
        expect(store.getState().comments.commentsByCard.card1).toHaveLength(2)
    })

    test('updateComments updates existing comment in parent', () => {
        const store = makeStore()
        const c1 = makeComment('c1', 'card1')
        store.dispatch(updateComments([c1]))
        const updatedC1 = {...c1, title: 'updated'}
        store.dispatch(updateComments([updatedC1]))

        // Should still have length 1 (updated in place)
        expect(store.getState().comments.commentsByCard.card1).toHaveLength(1)
    })

    test('updateComments deletes comment when deleteAt != 0', () => {
        const store = makeStore()
        const c1 = makeComment('c1', 'card1')
        store.dispatch(updateComments([c1]))
        const deletedC1 = makeComment('c1', 'card1', 999)
        store.dispatch(updateComments([deletedC1]))
        expect(store.getState().comments.comments.c1).toBeUndefined()
        expect(store.getState().comments.commentsByCard.card1).toHaveLength(0)
    })

    test('updateComments delete when no byCard entry', () => {
        const store = makeStore()
        const deletedC = makeComment('c99', 'card99', 999)
        store.dispatch(updateComments([deletedC]))
        expect(store.getState().comments.comments.c99).toBeUndefined()
    })

    test('extraReducer: initialReadOnlyLoad.fulfilled populates and sorts comments', () => {
        const store = makeStore()
        const c1 = makeComment('c1', 'card1', 0, 2000)
        const c2 = makeComment('c2', 'card1', 0, 1000)
        const nonComment = {id: 'card1', type: 'card', parentId: 'board1'}
        store.dispatch({
            type: 'initialReadOnlyLoad/fulfilled',
            payload: {board: {id: 'board1'}, blocks: [c1, c2, nonComment]},
        })
        expect(store.getState().comments.comments.c1).toBeTruthy()
        expect(store.getState().comments.comments.c2).toBeTruthy()
        const byCard = store.getState().comments.commentsByCard.card1

        // Sorted by createAt: c2 (1000) first, c1 (2000) second
        expect(byCard[0].id).toBe('c2')
        expect(byCard[1].id).toBe('c1')
    })

    test('extraReducer: loadBoardData.fulfilled populates and sorts comments', () => {
        const store = makeStore()
        const c1 = makeComment('c1', 'card1', 0, 3000)
        const c2 = makeComment('c2', 'card1', 0, 1500)
        store.dispatch({
            type: 'loadBoardData/fulfilled',
            payload: {blocks: [c1, c2]},
        })
        const byCard = store.getState().comments.commentsByCard.card1
        expect(byCard[0].id).toBe('c2')
        expect(byCard[1].id).toBe('c1')
    })

    test('extraReducer: loadBoardData.fulfilled clears previous comments', () => {
        const store = makeStore()
        const c1 = makeComment('c1', 'card1')
        store.dispatch(updateComments([c1]))
        store.dispatch({type: 'loadBoardData/fulfilled', payload: {blocks: []}})
        expect(Object.keys(store.getState().comments.comments)).toHaveLength(0)
    })
})

describe('getCardComments selector', () => {
    test('returns empty array for unknown cardId', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCardComments('card-none')(state)).toEqual([])
    })

    test('returns comments for known cardId', () => {
        const store = makeStore()
        const c1 = makeComment('c1', 'card1')
        store.dispatch(updateComments([c1]))
        const state = store.getState() as any
        expect(getCardComments('card1')(state)).toHaveLength(1)
    })
})

describe('getLastCardComment selector', () => {
    test('returns undefined for unknown cardId', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getLastCardComment('card-none')(state)).toBeUndefined()
    })

    test('returns last comment for known cardId', () => {
        const store = makeStore()
        const c1 = makeComment('c1', 'card1', 0, 1000)
        const c2 = makeComment('c2', 'card1', 0, 2000)
        store.dispatch(updateComments([c1]))
        store.dispatch(updateComments([c2]))
        const state = store.getState() as any
        const last = getLastCardComment('card1')(state)
        expect(last?.id).toBe('c2')
    })
})

describe('getLastCommentByCard selector', () => {
    test('returns empty object when no comments', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getLastCommentByCard(state)).toEqual({})
    })

    test('returns last comment for each card', () => {
        const store = makeStore()
        const c1 = makeComment('c1', 'card1', 0, 1000)
        const c2 = makeComment('c2', 'card1', 0, 2000)
        const c3 = makeComment('c3', 'card2', 0, 1500)
        store.dispatch(updateComments([c1]))
        store.dispatch(updateComments([c2]))
        store.dispatch(updateComments([c3]))
        const state = store.getState() as any
        const result = getLastCommentByCard(state)
        expect(result.card1?.id).toBe('c2')
        expect(result.card2?.id).toBe('c3')
    })
})
