// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {reducer as contentsReducer, updateContents, getContentsById, getContents, getCardContents, getLastCardContent} from '../contents'
import {reducer as cardsReducer} from '../cards'
import {createContentBlock} from '../../blocks/contentBlock'
import {createCard} from '../../blocks/card'

function makeStore() {
    return configureStore({reducer: {contents: contentsReducer, cards: cardsReducer}})
}

function makeContent(id: string, parentId: string, deleteAt = 0, createAt = 1000, type = 'text') {
    const block = createContentBlock()
    block.id = id
    block.parentId = parentId
    block.type = type as any
    block.deleteAt = deleteAt
    block.createAt = createAt
    return block
}

describe('contents reducer', () => {
    test('initial state is empty', () => {
        const store = makeStore()
        expect(store.getState().contents.contents).toEqual({})
        expect(store.getState().contents.contentsByCard).toEqual({})
    })

    test('updateContents adds new content to new parent', () => {
        const store = makeStore()
        const content = makeContent('cnt1', 'card1')
        store.dispatch(updateContents([content]))
        expect(store.getState().contents.contents['cnt1']).toBeTruthy()
        expect(store.getState().contents.contentsByCard['card1']).toHaveLength(1)
    })

    test('updateContents adds to existing parent', () => {
        const store = makeStore()
        const c1 = makeContent('cnt1', 'card1', 0, 1000)
        const c2 = makeContent('cnt2', 'card1', 0, 1001)
        store.dispatch(updateContents([c1]))
        store.dispatch(updateContents([c2]))
        expect(store.getState().contents.contentsByCard['card1']).toHaveLength(2)
    })

    test('updateContents updates existing content', () => {
        const store = makeStore()
        const c1 = makeContent('cnt1', 'card1')
        store.dispatch(updateContents([c1]))
        const updated = {...c1, title: 'updated title'}
        store.dispatch(updateContents([updated]))
        expect(store.getState().contents.contentsByCard['card1']).toHaveLength(1)
    })

    test('updateContents deletes content when deleteAt != 0', () => {
        const store = makeStore()
        const c1 = makeContent('cnt1', 'card1')
        store.dispatch(updateContents([c1]))
        const deleted = makeContent('cnt1', 'card1', 999)
        store.dispatch(updateContents([deleted]))
        expect(store.getState().contents.contents['cnt1']).toBeUndefined()
        expect(store.getState().contents.contentsByCard['card1']).toHaveLength(0)
    })

    test('updateContents delete when no byCard entry', () => {
        const store = makeStore()
        const deleted = makeContent('cnt99', 'card99', 999)
        store.dispatch(updateContents([deleted]))
        expect(store.getState().contents.contents['cnt99']).toBeUndefined()
    })

    test('extraReducer: initialReadOnlyLoad.fulfilled populates contents, ignores board/view/comment types', () => {
        const store = makeStore()
        const textBlock = makeContent('cnt1', 'card1', 0, 1000, 'text')
        const boardBlock = {...makeContent('board1', 'team1', 0, 1000, 'board')}
        const viewBlock = {...makeContent('view1', 'board1', 0, 1000, 'view')}
        const commentBlock = {...makeContent('cmt1', 'card1', 0, 1000, 'comment')}
        store.dispatch({
            type: 'initialReadOnlyLoad/fulfilled',
            payload: {
                board: {id: 'board1'},
                blocks: [textBlock, boardBlock, viewBlock, commentBlock],
            },
        })
        expect(store.getState().contents.contents['cnt1']).toBeTruthy()
        expect(store.getState().contents.contents['board1']).toBeUndefined()
        expect(store.getState().contents.contents['view1']).toBeUndefined()
        expect(store.getState().contents.contents['cmt1']).toBeUndefined()
    })

    test('extraReducer: loadBoardData.fulfilled populates contents', () => {
        const store = makeStore()
        const c1 = makeContent('cnt1', 'card1', 0, 1000)
        const c2 = makeContent('cnt2', 'card1', 0, 2000)
        store.dispatch({
            type: 'loadBoardData/fulfilled',
            payload: {blocks: [c2, c1]},
        })
        // Should be sorted by createAt
        const byCard = store.getState().contents.contentsByCard['card1']
        expect(byCard[0].id).toBe('cnt1')
        expect(byCard[1].id).toBe('cnt2')
    })

    test('extraReducer: loadBoardData.fulfilled clears previous contents', () => {
        const store = makeStore()
        const c1 = makeContent('cnt1', 'card1')
        store.dispatch(updateContents([c1]))
        store.dispatch({type: 'loadBoardData/fulfilled', payload: {blocks: []}})
        expect(Object.keys(store.getState().contents.contents)).toHaveLength(0)
    })
})

describe('getContentsById selector', () => {
    test('returns empty object initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getContentsById(state)).toEqual({})
    })

    test('returns contents after dispatch', () => {
        const store = makeStore()
        const c1 = makeContent('cnt1', 'card1')
        store.dispatch(updateContents([c1]))
        const state = store.getState() as any
        expect(getContentsById(state)['cnt1']).toBeTruthy()
    })
})

describe('getContents selector', () => {
    test('returns empty array initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getContents(state)).toEqual([])
    })

    test('returns all contents as array', () => {
        const store = makeStore()
        const c1 = makeContent('cnt1', 'card1')
        const c2 = makeContent('cnt2', 'card2')
        store.dispatch(updateContents([c1]))
        store.dispatch(updateContents([c2]))
        const state = store.getState() as any
        expect(getContents(state)).toHaveLength(2)
    })
})

describe('getCardContents selector', () => {
    test('returns empty array for unknown cardId', () => {
        const store = makeStore()
        const state = store.getState() as any
        const selector = getCardContents('card-none')
        expect(selector(state)).toEqual([])
    })

    test('returns contents in contentOrder from card', () => {
        const store = makeStore()
        const c1 = makeContent('cnt1', 'card1')
        const c2 = makeContent('cnt2', 'card1')

        // Add a card with known contentOrder
        const card = createCard()
        card.id = 'card1'
        card.boardId = 'board1'
        card.fields.contentOrder = ['cnt1', 'cnt2']

        store.dispatch({
            type: 'cards/setCurrent',
            payload: 'card1',
        })
        store.dispatch({
            type: 'loadBoardData/fulfilled',
            payload: {blocks: [c1, c2, {...card, type: 'card', fields: {...card.fields, isTemplate: false}}]},
        })

        const state = store.getState() as any
        const selector = getCardContents('card1')
        const result = selector(state)
        expect(result).toHaveLength(2)
    })

    test('getLastCardContent returns undefined for unknown cardId', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getLastCardContent('unknown')(state)).toBeUndefined()
    })

    test('getLastCardContent returns last content block', () => {
        const store = makeStore()
        const c1 = makeContent('cnt1', 'card1', 0, 1000)
        const c2 = makeContent('cnt2', 'card1', 0, 2000)
        store.dispatch(updateContents([c1]))
        store.dispatch(updateContents([c2]))
        const state = store.getState() as any
        const last = getLastCardContent('card1')(state)
        expect(last?.id).toBe('cnt2')
    })
})
