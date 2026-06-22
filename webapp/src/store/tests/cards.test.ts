// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {
    reducer as cardsReducer,    updateCards,
    addCard,
    addTemplate,
    setCurrent,
    setLimitTimestamp,
    showCardHiddenWarning,
    getCards,
    getSortedCards,
    getTemplates,
    getSortedTemplates,
    getCard,
    getCurrentBoardCards,
    getCurrentBoardTemplates,
    getCurrentCard,
    getCardLimitTimestamp,
    getCardHiddenWarning,
    getCurrentViewCardsSortedFilteredAndGrouped,
    getCurrentViewCardsSortedFilteredAndGroupedWithoutLimit,
    getCurrentBoardHiddenCardsCount,
} from '../cards'
import {reducer as boardsReducer} from '../boards'
import {reducer as viewsReducer} from '../views'
import {reducer as commentsReducer} from '../comments'
import {reducer as usersReducer} from '../users'
import {reducer as searchTextReducer} from '../searchText'
import {createCard} from '../../blocks/card'
import {createBoard} from '../../blocks/board'
import {createBoardView} from '../../blocks/boardView'
import {Constants} from '../../constants'

// Mock octoClient
jest.mock('../../octoClient', () => ({
    default: {
        getBlocksWithBlockID: jest.fn().mockResolvedValue([]),
    },
}))

jest.mock('../../utils', () => ({
    Utils: {
        log: jest.fn(),
        logError: jest.fn(),
        createGuid: jest.fn(() => 'mock-guid'),
        blockTypeToIDType: jest.fn(() => 'block'),
    },
    IDType: {
        BlockID: 'block',
    },
}))
function makeStore() {
    return configureStore({
        reducer: {
            cards: cardsReducer,
            boards: boardsReducer,
            views: viewsReducer,
            comments: commentsReducer,
            users: usersReducer,
            searchText: searchTextReducer,
        },
    })
}

function makeCard(id: string, boardId: string, title = 'Card Title', isTemplate = false, deleteAt = 0) {
    const card = createCard()
    card.id = id
    card.boardId = boardId
    card.title = title
    card.deleteAt = deleteAt
    card.fields.isTemplate = isTemplate
    return card
}

describe('cards reducer', () => {
    test('initial state is empty', () => {
        const store = makeStore()
        expect(store.getState().cards.cards).toEqual({})
        expect(store.getState().cards.templates).toEqual({})
        expect(store.getState().cards.current).toBe('')
        expect(store.getState().cards.limitTimestamp).toBe(0)
        expect(store.getState().cards.cardHiddenWarning).toBe(false)
    })

    test('setCurrent sets the current card id', () => {
        const store = makeStore()
        store.dispatch(setCurrent('card1'))
        expect(store.getState().cards.current).toBe('card1')
    })

    test('showCardHiddenWarning sets cardHiddenWarning', () => {
        const store = makeStore()
        store.dispatch(showCardHiddenWarning(true))
        expect(store.getState().cards.cardHiddenWarning).toBe(true)
        store.dispatch(showCardHiddenWarning(false))
        expect(store.getState().cards.cardHiddenWarning).toBe(false)
    })

    test('addCard adds a card', () => {
        const store = makeStore()
        const card = makeCard('c1', 'board1')
        store.dispatch(addCard(card))
        expect(store.getState().cards.cards['c1']).toBeTruthy()
    })

    test('addTemplate adds a template', () => {
        const store = makeStore()
        const card = makeCard('tmpl1', 'board1', 'Template', true)
        store.dispatch(addTemplate(card))
        expect(store.getState().cards.templates['tmpl1']).toBeTruthy()
    })

    test('updateCards adds non-template card', () => {
        const store = makeStore()
        const card = makeCard('c1', 'board1')
        store.dispatch(updateCards([card]))
        expect(store.getState().cards.cards['c1']).toBeTruthy()
    })

    test('updateCards adds template card to templates', () => {
        const store = makeStore()
        const card = makeCard('tmpl1', 'board1', 'Template Card', true)
        store.dispatch(updateCards([card]))
        expect(store.getState().cards.templates['tmpl1']).toBeTruthy()
        expect(store.getState().cards.cards['tmpl1']).toBeUndefined()
    })

    test('updateCards removes card when deleteAt != 0', () => {
        const store = makeStore()
        const card = makeCard('c1', 'board1')
        store.dispatch(updateCards([card]))
        const deleted = makeCard('c1', 'board1', 'Card', false, 999)
        store.dispatch(updateCards([deleted]))
        expect(store.getState().cards.cards['c1']).toBeUndefined()
    })

    test('setLimitTimestamp sets the limitTimestamp', () => {
        const store = makeStore()
        store.dispatch(setLimitTimestamp({timestamp: 12345, templates: {}}))
        expect(store.getState().cards.limitTimestamp).toBe(12345)
    })

    test('setLimitTimestamp limits cards older than timestamp', () => {
        const store = makeStore()
        const oldCard = makeCard('c1', 'board1', 'Old Card')
        oldCard.updateAt = 100
        store.dispatch(addCard(oldCard))
        store.dispatch(setLimitTimestamp({timestamp: 500, templates: {}}))
        const limitedCard = store.getState().cards.cards['c1']
        expect(limitedCard?.limited).toBe(true)
    })

    test('setLimitTimestamp does not limit cards newer than timestamp', () => {
        const store = makeStore()
        const newCard = makeCard('c1', 'board1', 'New Card')
        newCard.updateAt = 1000
        store.dispatch(addCard(newCard))
        store.dispatch(setLimitTimestamp({timestamp: 500, templates: {}}))
        const card = store.getState().cards.cards['c1']
        expect(card?.limited).toBeFalsy()
    })

    test('extraReducer: initialReadOnlyLoad.fulfilled populates cards and templates', () => {
        const store = makeStore()
        const card = makeCard('c1', 'board1')
        const template = makeCard('tmpl1', 'board1', 'Template', true)
        const nonCard = {id: 'view1', type: 'view', parentId: 'board1', fields: {}}
        store.dispatch({
            type: 'initialReadOnlyLoad/fulfilled',
            payload: {board: {id: 'board1'}, blocks: [card, template, nonCard]},
        })
        expect(store.getState().cards.cards['c1']).toBeTruthy()
        expect(store.getState().cards.templates['tmpl1']).toBeTruthy()
    })

    test('extraReducer: loadBoardData.fulfilled populates cards', () => {
        const store = makeStore()
        const card = makeCard('c1', 'board1')
        store.dispatch({
            type: 'loadBoardData/fulfilled',
            payload: {blocks: [card]},
        })
        expect(store.getState().cards.cards['c1']).toBeTruthy()
    })

    test('extraReducer: loadBoardData.fulfilled clears previous cards', () => {
        const store = makeStore()
        const card = makeCard('c1', 'board1')
        store.dispatch(addCard(card))
        store.dispatch({type: 'loadBoardData/fulfilled', payload: {blocks: []}})
        expect(Object.keys(store.getState().cards.cards)).toHaveLength(0)
    })

    test('extraReducer: initialLoad.fulfilled sets limitTimestamp', () => {
        const store = makeStore()
        store.dispatch({
            type: 'initialLoad/fulfilled',
            payload: {
                team: {id: 't1'},
                teams: [],
                boards: [],
                boardsMemberships: [],
                boardTemplates: [],
                limits: {card_limit_timestamp: 99999, cards: 100, used_cards: 50, views: 10},
                myConfig: [],
            },
        })
        expect(store.getState().cards.limitTimestamp).toBe(99999)
    })

    test('extraReducer: refreshCards/fulfilled updates cards', () => {
        const store = makeStore()
        const card = makeCard('c1', 'board1')
        store.dispatch(addCard(card))
        const updatedCard = {...card, title: 'Updated'}
        store.dispatch({
            type: 'refreshCards/fulfilled',
            payload: [updatedCard],
        })
        expect(store.getState().cards.cards['c1']?.title).toBe('Updated')
    })
})

describe('cards selectors', () => {
    test('getCards returns empty object initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCards(state)).toEqual({})
    })

    test('getSortedCards returns cards sorted by title', () => {
        const store = makeStore()
        const c1 = makeCard('c1', 'board1', 'Zebra')
        const c2 = makeCard('c2', 'board1', 'Apple')
        store.dispatch(updateCards([c1, c2]))
        const state = store.getState() as any
        const sorted = getSortedCards(state)
        expect(sorted[0].title).toBe('Apple')
        expect(sorted[1].title).toBe('Zebra')
    })

    test('getTemplates returns empty object initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getTemplates(state)).toEqual({})
    })

    test('getSortedTemplates returns templates sorted by title', () => {
        const store = makeStore()
        const t1 = makeCard('t1', 'board1', 'Zeta Template', true)
        const t2 = makeCard('t2', 'board1', 'Alpha Template', true)
        store.dispatch(updateCards([t1, t2]))
        const state = store.getState() as any
        const sorted = getSortedTemplates(state)
        expect(sorted[0].title).toBe('Alpha Template')
    })

    test('getCard returns undefined for unknown cardId', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCard('unknown')(state)).toBeUndefined()
    })

    test('getCard returns card from cards', () => {
        const store = makeStore()
        const card = makeCard('c1', 'board1')
        store.dispatch(addCard(card))
        const state = store.getState() as any
        expect(getCard('c1')(state)?.id).toBe('c1')
    })

    test('getCard returns card from templates', () => {
        const store = makeStore()
        const tmpl = makeCard('tmpl1', 'board1', 'Template', true)
        store.dispatch(addTemplate(tmpl))
        const state = store.getState() as any
        expect(getCard('tmpl1')(state)?.id).toBe('tmpl1')
    })

    test('getCurrentBoardCards returns empty when no board set', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentBoardCards(state)).toEqual([])
    })

    test('getCurrentBoardCards returns cards for current board', () => {
        const store = makeStore()
        const board = createBoard()
        store.dispatch({type: 'boards/setCurrent', payload: board.id})
        const card = makeCard('c1', board.id)
        store.dispatch(addCard(card))
        const state = store.getState() as any
        expect(getCurrentBoardCards(state)).toHaveLength(1)
    })

    test('getCurrentBoardTemplates returns templates for current board', () => {
        const store = makeStore()
        const board = createBoard()
        store.dispatch({type: 'boards/setCurrent', payload: board.id})
        const tmpl = makeCard('tmpl1', board.id, 'Template', true)
        store.dispatch(addTemplate(tmpl))
        const state = store.getState() as any
        expect(getCurrentBoardTemplates(state)).toHaveLength(1)
    })

    test('getCurrentCard returns undefined when no card selected', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentCard(state)).toBeUndefined()
    })

    test('getCurrentCard returns selected card', () => {
        const store = makeStore()
        const card = makeCard('c1', 'board1')
        store.dispatch(addCard(card))
        store.dispatch(setCurrent('c1'))
        const state = store.getState() as any
        expect(getCurrentCard(state)?.id).toBe('c1')
    })

    test('getCardLimitTimestamp returns 0 initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCardLimitTimestamp(state)).toBe(0)
    })

    test('getCardHiddenWarning returns false initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCardHiddenWarning(state)).toBe(false)
    })

    test('getCurrentViewCardsSortedFilteredAndGrouped returns empty when no view', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentViewCardsSortedFilteredAndGrouped(state)).toEqual([])
    })

    test('getCurrentViewCardsSortedFilteredAndGroupedWithoutLimit returns empty when no view', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentViewCardsSortedFilteredAndGroupedWithoutLimit(state)).toEqual([])
    })

    test('getCurrentBoardHiddenCardsCount returns 0 when no hidden cards', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentBoardHiddenCardsCount(state)).toBe(0)
    })

    test('getCurrentBoardHiddenCardsCount counts limited cards', () => {
        const store = makeStore()
        const board = createBoard()
        store.dispatch({type: 'boards/setCurrent', payload: board.id})
        const card = makeCard('c1', board.id)
        card.limited = true
        store.dispatch(addCard(card))
        const state = store.getState() as any
        expect(getCurrentBoardHiddenCardsCount(state)).toBe(1)
    })

    test('getCurrentViewCardsSortedFilteredAndGrouped returns cards with full setup', () => {
        const store = makeStore()
        const board = createBoard()
        board.id = 'board1'
        const view = createBoardView()
        view.id = 'view1'
        view.boardId = 'board1'
        view.parentId = 'board1'
        view.fields.sortOptions = []

        store.dispatch({type: 'boards/setCurrent', payload: 'board1'})
        store.dispatch({
            type: 'initialLoad/fulfilled',
            payload: {
                team: {id: 't1'},
                teams: [],
                boards: [board],
                boardsMemberships: [],
                boardTemplates: [],
                limits: null,
                myConfig: [],
            },
        })
        store.dispatch({type: 'views/setCurrent', payload: 'view1'})
        store.dispatch({type: 'loadBoardData/fulfilled', payload: {blocks: [view]}})

        const card = makeCard('c1', 'board1', 'Card 1')
        store.dispatch(addCard(card))

        const state = store.getState() as any
        const result = getCurrentViewCardsSortedFilteredAndGrouped(state)
        expect(result).toHaveLength(1)
    })
})

describe('cards extra coverage', () => {
    test('setLimitTimestamp preserves board template cards even if old', () => {
        const store = makeStore()

        const templateCard = makeCard('template-card', 'board1', 'Template Card')
        templateCard.updateAt = 100
        templateCard.fields.properties = {prop1: 'value'}

        store.dispatch(addCard(templateCard))
        store.dispatch(setLimitTimestamp({
            timestamp: 200,
            templates: {
                'template-card': {id: 'template-card'} as any,
            },
        }))

        expect(store.getState().cards.cards['template-card'].limited).toBeFalsy()
        expect(store.getState().cards.cards['template-card'].fields.properties).toEqual({prop1: 'value'})
    })
})



describe('cards sorting and search coverage', () => {
    function setupBoardAndView(store: any, board: any, view: any) {
        store.dispatch({type: 'boards/setCurrent', payload: board.id})
        store.dispatch({
            type: 'initialLoad/fulfilled',
            payload: {
                team: {id: 't1'},
                teams: [],
                boards: [board],
                boardsMemberships: [],
                boardTemplates: [],
                limits: null,
                myConfig: [],
            },
        })
        store.dispatch({type: 'views/setCurrent', payload: view.id})
        store.dispatch({type: 'loadBoardData/fulfilled', payload: {blocks: [view]}})
    }

    test('manual sort uses cardOrder first and title fallback', () => {
        const store = makeStore()
        const board = createBoard()
        board.id = 'board-manual-extra'

        const view = createBoardView()
        view.id = 'view-manual-extra'
        view.boardId = board.id
        view.parentId = board.id
        view.fields.sortOptions = []
        view.fields.cardOrder = ['c2']

        setupBoardAndView(store, board, view)

        store.dispatch(addCard(makeCard('c1', board.id, 'Alpha')))
        store.dispatch(addCard(makeCard('c2', board.id, 'Beta')))
        store.dispatch(addCard(makeCard('c3', board.id, 'Gamma')))

        const result = getCurrentViewCardsSortedFilteredAndGrouped(store.getState() as any)

        expect(result.map((c) => c.id)).toEqual(['c1', 'c2', 'c3'])
    })

    test('sorts cards by title reversed', () => {
        const store = makeStore()
        const board = createBoard()
        board.id = 'board-title-reversed'

        const view = createBoardView()
        view.id = 'view-title-reversed'
        view.boardId = board.id
        view.parentId = board.id
        view.fields.sortOptions = [{propertyId: Constants.titleColumnId, reversed: true}]
        view.fields.cardOrder = []

        setupBoardAndView(store, board, view)

        store.dispatch(addCard(makeCard('c1', board.id, 'Alpha')))
        store.dispatch(addCard(makeCard('c2', board.id, 'Beta')))

        const result = getCurrentViewCardsSortedFilteredAndGrouped(store.getState() as any)

        expect(result.map((c) => c.id)).toEqual(['c2', 'c1'])
    })

    test('sorts cards by number property and places empty values last', () => {
        const store = makeStore()
        const board = createBoard()
        board.id = 'board-number-extra'
        board.cardProperties = [{id: 'num', name: 'Number', type: 'number', options: []} as any]

        const view = createBoardView()
        view.id = 'view-number-extra'
        view.boardId = board.id
        view.parentId = board.id
        view.fields.sortOptions = [{propertyId: 'num', reversed: false}]
        view.fields.cardOrder = []

        setupBoardAndView(store, board, view)

        const c1 = makeCard('c1', board.id, 'Card 10')
        c1.fields.properties.num = '10'
        const c2 = makeCard('c2', board.id, 'Card 2')
        c2.fields.properties.num = '2'
        const c3 = makeCard('c3', board.id, 'Empty')
        c3.fields.properties.num = ''

        store.dispatch(addCard(c1))
        store.dispatch(addCard(c2))
        store.dispatch(addCard(c3))

        const result = getCurrentViewCardsSortedFilteredAndGrouped(store.getState() as any)

        expect(result.map((c) => c.id)).toEqual(['c2', 'c1', 'c3'])
    })

    test('sorts cards by date property', () => {
        const store = makeStore()
        const board = createBoard()
        board.id = 'board-date-extra'
        board.cardProperties = [{id: 'date', name: 'Date', type: 'date', options: []} as any]

        const view = createBoardView()
        view.id = 'view-date-extra'
        view.boardId = board.id
        view.parentId = board.id
        view.fields.sortOptions = [{propertyId: 'date', reversed: false}]
        view.fields.cardOrder = []

        setupBoardAndView(store, board, view)

        const c1 = makeCard('c1', board.id, 'Later')
        c1.fields.properties.date = JSON.stringify({from: 3000})
        const c2 = makeCard('c2', board.id, 'Earlier')
        c2.fields.properties.date = JSON.stringify({from: 1000})

        store.dispatch(addCard(c1))
        store.dispatch(addCard(c2))

        const result = getCurrentViewCardsSortedFilteredAndGrouped(store.getState() as any)

        expect(result.map((c) => c.id)).toEqual(['c2', 'c1'])
    })

    test('sorts cards by select option display value', () => {
        const store = makeStore()
        const board = createBoard()
        board.id = 'board-select-extra'
        board.cardProperties = [{
            id: 'status',
            name: 'Status',
            type: 'select',
            options: [
                {id: 'todo', value: 'Todo'},
                {id: 'done', value: 'Done'},
            ],
        } as any]

        const view = createBoardView()
        view.id = 'view-select-extra'
        view.boardId = board.id
        view.parentId = board.id
        view.fields.sortOptions = [{propertyId: 'status', reversed: false}]
        view.fields.cardOrder = []

        setupBoardAndView(store, board, view)

        const c1 = makeCard('c1', board.id, 'Task A')
        c1.fields.properties.status = 'todo'
        const c2 = makeCard('c2', board.id, 'Task B')
        c2.fields.properties.status = 'done'

        store.dispatch(addCard(c1))
        store.dispatch(addCard(c2))

        const result = getCurrentViewCardsSortedFilteredAndGrouped(store.getState() as any)

        expect(result.map((c) => c.id)).toEqual(['c2', 'c1'])
    })

    test('search filters cards by title and select property value', () => {
        const store = makeStore()
        const board = createBoard()
        board.id = 'board-search-extra'
        board.cardProperties = [{
            id: 'status',
            name: 'Status',
            type: 'select',
            options: [
                {id: 'blocked', value: 'Blocked'},
                {id: 'open', value: 'Open'},
            ],
        } as any]

        const view = createBoardView()
        view.id = 'view-search-extra'
        view.boardId = board.id
        view.parentId = board.id
        view.fields.sortOptions = []
        view.fields.cardOrder = []

        setupBoardAndView(store, board, view)

        const c1 = makeCard('c1', board.id, 'Normal title')
        c1.fields.properties.status = 'blocked'
        const c2 = makeCard('c2', board.id, 'Another card')
        c2.fields.properties.status = 'open'

        store.dispatch(addCard(c1))
        store.dispatch(addCard(c2))
        store.dispatch({type: 'searchText/setSearchText', payload: 'blocked'})

        const result = getCurrentViewCardsSortedFilteredAndGrouped(store.getState() as any)

        expect(result.map((c) => c.id)).toEqual(['c1'])
    })
})
