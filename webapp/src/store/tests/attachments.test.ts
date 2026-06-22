// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {reducer, updateAttachments, updateUploadPrecent, getCardAttachments, getUploadPercent} from '../attachments'
import {createAttachmentBlock} from '../../blocks/attachmentBlock'

function makeStore() {
    return configureStore({reducer: {attachments: reducer}})
}

function makeAttachment(id: string, parentId: string, deleteAt = 0, createAt = 1000) {
    const att = createAttachmentBlock()
    att.id = id
    att.parentId = parentId
    att.deleteAt = deleteAt
    att.createAt = createAt
    att.uploadingPercent = 0
    return att
}

describe('attachments reducer', () => {
    test('initial state is empty', () => {
        const store = makeStore()
        expect(store.getState().attachments.attachments).toEqual({})
        expect(store.getState().attachments.attachmentsByCard).toEqual({})
    })

    test('updateAttachments adds new attachment (new parent)', () => {
        const store = makeStore()
        const attachment = makeAttachment('att1', 'card1')
        store.dispatch(updateAttachments([attachment]))
        expect(store.getState().attachments.attachments.att1).toBeTruthy()
        expect(store.getState().attachments.attachmentsByCard.card1).toHaveLength(1)
    })

    test('updateAttachments adds to existing parent', () => {
        const store = makeStore()
        const att1 = makeAttachment('att1', 'card1', 0, 1000)
        const att2 = makeAttachment('att2', 'card1', 0, 1001)
        store.dispatch(updateAttachments([att1]))
        store.dispatch(updateAttachments([att2]))
        expect(store.getState().attachments.attachmentsByCard.card1).toHaveLength(2)
    })

    test('updateAttachments does not duplicate if same id exists', () => {
        const store = makeStore()
        const att1 = makeAttachment('att1', 'card1')
        store.dispatch(updateAttachments([att1]))
        store.dispatch(updateAttachments([att1]))
        expect(store.getState().attachments.attachmentsByCard.card1).toHaveLength(1)
    })

    test('updateAttachments removes attachment when deleteAt != 0', () => {
        const store = makeStore()
        const att1 = makeAttachment('att1', 'card1')
        store.dispatch(updateAttachments([att1]))
        const deletedAtt = makeAttachment('att1', 'card1', 999)
        store.dispatch(updateAttachments([deletedAtt]))
        expect(store.getState().attachments.attachments.att1).toBeUndefined()
        expect(store.getState().attachments.attachmentsByCard.card1).toHaveLength(0)
    })

    test('updateAttachments delete with no existing parentId in byCard deletes from attachments', () => {
        const store = makeStore()

        // Attachment doesn't exist in store yet, so parentId lookup returns undefined
        const deletedAtt = makeAttachment('att99', 'card99', 999)
        store.dispatch(updateAttachments([deletedAtt]))
        expect(store.getState().attachments.attachments.att99).toBeUndefined()
    })

    test('updateUploadPrecent updates uploadingPercent on an attachment', () => {
        const store = makeStore()
        const att1 = makeAttachment('att1', 'card1')
        store.dispatch(updateAttachments([att1]))
        store.dispatch(updateUploadPrecent({blockId: 'att1', uploadPercent: 75}))
        expect(store.getState().attachments.attachments.att1.uploadingPercent).toBe(75)
    })

    test('extraReducer: initialReadOnlyLoad.fulfilled populates and sorts attachments', () => {
        const store = makeStore()
        const att1 = makeAttachment('att1', 'card1', 0, 2000)
        const att2 = makeAttachment('att2', 'card1', 0, 1000)
        store.dispatch({
            type: 'initialReadOnlyLoad/fulfilled',
            payload: {
                board: {id: 'board1', isTemplate: false},
                blocks: [att1, att2],
            },
        })
        expect(store.getState().attachments.attachments.att1).toBeTruthy()
        expect(store.getState().attachments.attachments.att2).toBeTruthy()

        // Should be sorted by createAt (att2=1000 first, att1=2000 second)
        const byCard = store.getState().attachments.attachmentsByCard.card1
        expect(byCard[0].id).toBe('att2')
        expect(byCard[1].id).toBe('att1')
    })

    test('extraReducer: initialReadOnlyLoad.fulfilled ignores non-attachment blocks', () => {
        const store = makeStore()
        const nonAttachment = {id: 'card1', type: 'card', parentId: 'board1'}
        store.dispatch({
            type: 'initialReadOnlyLoad/fulfilled',
            payload: {board: {id: 'board1'}, blocks: [nonAttachment]},
        })
        expect(Object.keys(store.getState().attachments.attachments)).toHaveLength(0)
    })

    test('extraReducer: loadBoardData.fulfilled populates attachments', () => {
        const store = makeStore()
        const att1 = makeAttachment('att1', 'card1', 0, 1000)
        store.dispatch({
            type: 'loadBoardData/fulfilled',
            payload: {blocks: [att1]},
        })
        expect(store.getState().attachments.attachments.att1).toBeTruthy()
    })

    test('extraReducer: loadBoardData.fulfilled clears previous attachments', () => {
        const store = makeStore()
        const att1 = makeAttachment('att1', 'card1')
        store.dispatch(updateAttachments([att1]))
        store.dispatch({
            type: 'loadBoardData/fulfilled',
            payload: {blocks: []},
        })
        expect(Object.keys(store.getState().attachments.attachments)).toHaveLength(0)
    })

    test('extraReducer: multiple attachments for same card are sorted', () => {
        const store = makeStore()
        const att1 = makeAttachment('att1', 'card1', 0, 3000)
        const att2 = makeAttachment('att2', 'card1', 0, 1000)
        const att3 = makeAttachment('att3', 'card1', 0, 2000)
        store.dispatch({
            type: 'loadBoardData/fulfilled',
            payload: {blocks: [att1, att2, att3]},
        })
        const byCard = store.getState().attachments.attachmentsByCard.card1
        expect(byCard[0].id).toBe('att2')
        expect(byCard[1].id).toBe('att3')
        expect(byCard[2].id).toBe('att1')
    })
})

describe('getCardAttachments selector', () => {
    test('returns empty array for unknown cardId', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCardAttachments('card-none')(state)).toEqual([])
    })

    test('returns attachments for known cardId', () => {
        const store = makeStore()
        const att1 = makeAttachment('att1', 'card1')
        store.dispatch(updateAttachments([att1]))
        const state = store.getState() as any
        expect(getCardAttachments('card1')(state)).toHaveLength(1)
    })
})

describe('getUploadPercent selector', () => {
    test('returns the upload percent for an attachment', () => {
        const store = makeStore()
        const att1 = makeAttachment('att1', 'card1')
        store.dispatch(updateAttachments([att1]))
        store.dispatch(updateUploadPrecent({blockId: 'att1', uploadPercent: 50}))
        const state = store.getState() as any
        expect(getUploadPercent('att1')(state)).toBe(50)
    })
})
