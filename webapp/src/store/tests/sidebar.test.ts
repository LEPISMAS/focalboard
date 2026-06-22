import {configureStore} from '@reduxjs/toolkit'

import {
    reducer as sidebarReducer,
    updateCategories,
    updateBoardCategories,
    updateCategoryOrder,
    updateCategoryBoardsOrder,
    fetchSidebarCategories,
    getSidebarCategories,
    getHiddenBoardIDs,
    getCategoryOfBoard,
} from '../sidebar'

jest.mock('../../octoClient', () => ({
    default: {
        getSidebarCategories: jest.fn(),
    },
}))

jest.mock('../../utils', () => ({
    Utils: {
        logError: jest.fn(),
    },
}))

function makeStore() {
    return configureStore({
        reducer: {
            sidebar: sidebarReducer,
        },
    })
}

function makeCategory(id: string, name = 'Category', deleteAt = 0) {
    return {
        id,
        name,
        userID: 'user1',
        teamID: 'team1',
        createAt: 1,
        updateAt: 1,
        deleteAt,
        collapsed: false,
        sortOrder: 0,
        type: 'custom' as const,
        isNew: false,
        boardMetadata: [],
    }
}

describe('sidebar reducer', () => {
    test('initial state is empty', () => {
        const store = makeStore()
        expect(store.getState().sidebar.categoryAttributes).toEqual([])
        expect(store.getState().sidebar.hiddenBoardIDs).toEqual([])
    })

    test('updateCategories adds new category', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1', 'Category 1')]))

        expect(store.getState().sidebar.categoryAttributes).toHaveLength(1)
        expect(store.getState().sidebar.categoryAttributes[0].id).toBe('cat1')
        expect(store.getState().sidebar.categoryAttributes[0].isNew).toBe(true)
    })

    test('updateCategories updates existing category', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1', 'Old')]))
        store.dispatch(updateCategories([{...makeCategory('cat1', 'New'), updateAt: 99}]))

        expect(store.getState().sidebar.categoryAttributes[0].name).toBe('New')
        expect(store.getState().sidebar.categoryAttributes[0].updateAt).toBe(99)
        expect(store.getState().sidebar.categoryAttributes[0].isNew).toBe(false)
    })

    test('updateCategories removes deleted category', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1')]))
        store.dispatch(updateCategories([makeCategory('cat1', 'Category', 123)]))

        expect(store.getState().sidebar.categoryAttributes).toHaveLength(0)
    })

    test('updateBoardCategories adds board to category', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1')]))
        store.dispatch(updateBoardCategories([{boardID: 'board1', categoryID: 'cat1', hidden: false}]))

        expect(store.getState().sidebar.categoryAttributes[0].boardMetadata[0]).toEqual({
            boardID: 'board1',
            hidden: false,
        })
    })

    test('updateBoardCategories updates hidden board ids', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1')]))
        store.dispatch(updateBoardCategories([{boardID: 'board1', categoryID: 'cat1', hidden: true}]))

        expect(store.getState().sidebar.hiddenBoardIDs).toEqual(['board1'])

        store.dispatch(updateBoardCategories([{boardID: 'board1', categoryID: 'cat1', hidden: false}]))

        expect(store.getState().sidebar.hiddenBoardIDs).toEqual([])
    })

    test('updateBoardCategories moves board between categories', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1'), makeCategory('cat2')]))

        store.dispatch(updateBoardCategories([{boardID: 'board1', categoryID: 'cat1', hidden: false}]))
        store.dispatch(updateBoardCategories([{boardID: 'board1', categoryID: 'cat2', hidden: false}]))

        const cat1 = store.getState().sidebar.categoryAttributes.find((c) => c.id === 'cat1')
        const cat2 = store.getState().sidebar.categoryAttributes.find((c) => c.id === 'cat2')

        expect(cat1?.boardMetadata.find((m) => m.boardID === 'board1')).toBeUndefined()
        expect(cat2?.boardMetadata.find((m) => m.boardID === 'board1')).toBeTruthy()
    })

    test('updateCategoryOrder does nothing with empty payload', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1')]))

        store.dispatch(updateCategoryOrder([]))

        expect(store.getState().sidebar.categoryAttributes).toHaveLength(1)
    })

    test('updateCategoryOrder reorders categories', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1'), makeCategory('cat2')]))

        store.dispatch(updateCategoryOrder(['cat1', 'cat2']))

        expect(store.getState().sidebar.categoryAttributes.map((c) => c.id)).toEqual(['cat1', 'cat2'])
    })

    test('updateCategoryOrder ignores unknown category ids', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1')]))

        store.dispatch(updateCategoryOrder(['missing']))

        expect(store.getState().sidebar.categoryAttributes).toEqual([])
    })

    test('updateCategoryBoardsOrder does nothing with empty metadata', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1')]))

        store.dispatch(updateCategoryBoardsOrder({
            categoryID: 'cat1',
            boardsMetadata: [],
        }))

        expect(store.getState().sidebar.categoryAttributes[0].boardMetadata).toEqual([])
    })

    test('updateCategoryBoardsOrder updates board metadata order', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1')]))

        store.dispatch(updateCategoryBoardsOrder({
            categoryID: 'cat1',
            boardsMetadata: [
                {boardID: 'board2', hidden: false},
                {boardID: 'board1', hidden: true},
            ],
        }))

        expect(store.getState().sidebar.categoryAttributes[0].boardMetadata.map((m) => m.boardID)).toEqual(['board2', 'board1'])
        expect(store.getState().sidebar.categoryAttributes[0].isNew).toBe(false)
    })

    test('updateCategoryBoardsOrder ignores unknown category', () => {
        const store = makeStore()

        store.dispatch(updateCategoryBoardsOrder({
            categoryID: 'missing',
            boardsMetadata: [{boardID: 'board1', hidden: false}],
        }))

        expect(store.getState().sidebar.categoryAttributes).toEqual([])
    })

    test('fetchSidebarCategories.fulfilled loads categories and hidden boards', () => {
        const store = makeStore()

        const categories = [
            {
                ...makeCategory('cat1'),
                boardMetadata: [
                    {boardID: 'board1', hidden: true},
                    {boardID: 'board2', hidden: false},
                ],
            },
            {
                ...makeCategory('cat2'),
                boardMetadata: [
                    {boardID: 'board3', hidden: true},
                ],
            },
        ]

        store.dispatch(fetchSidebarCategories.fulfilled(categories, 'req1', 'team1'))

        expect(store.getState().sidebar.categoryAttributes).toEqual(categories)
        expect(store.getState().sidebar.hiddenBoardIDs).toEqual(['board1', 'board3'])
    })

    test('fetchSidebarCategories.fulfilled with null payload clears state', () => {
        const store = makeStore()

        store.dispatch(fetchSidebarCategories.fulfilled(null as any, 'req1', 'team1'))

        expect(store.getState().sidebar.categoryAttributes).toEqual([])
        expect(store.getState().sidebar.hiddenBoardIDs).toEqual([])
    })
})

describe('sidebar selectors', () => {
    test('getSidebarCategories returns categories', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1')]))

        expect(getSidebarCategories(store.getState() as any)[0].id).toBe('cat1')
    })

    test('getHiddenBoardIDs returns hidden board ids', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1')]))
        store.dispatch(updateBoardCategories([{boardID: 'board1', categoryID: 'cat1', hidden: true}]))

        expect(getHiddenBoardIDs(store.getState() as any)).toEqual(['board1'])
    })

    test('getCategoryOfBoard returns category containing board', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1')]))
        store.dispatch(updateBoardCategories([{boardID: 'board1', categoryID: 'cat1', hidden: false}]))

        expect(getCategoryOfBoard('board1')(store.getState() as any)?.id).toBe('cat1')
    })

    test('getCategoryOfBoard returns undefined for missing board', () => {
        const store = makeStore()
        store.dispatch(updateCategories([makeCategory('cat1')]))

        expect(getCategoryOfBoard('missing')(store.getState() as any)).toBeUndefined()
    })
})
