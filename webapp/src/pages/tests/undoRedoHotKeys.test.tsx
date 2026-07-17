// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {render} from '@testing-library/react'

import {wrapIntl} from '../../testUtils'
import mutator from '../../mutator'
import {sendFlashMessage} from '../../components/flashMessages'
import UndoRedoHotKeys from '../boardPage/undoRedoHotKeys'

const registeredHotkeys: Record<string, () => void> = {}

jest.mock('react-hotkeys-hook', () => {
    return {
        useHotkeys: jest.fn((keys, callback) => {
            registeredHotkeys[keys] = callback
        }),
    }
})

jest.mock('../../mutator', () => {
    const mockMutator = {
        undo: jest.fn(() => Promise.resolve()),
        redo: jest.fn(() => Promise.resolve()),
    }
    Object.defineProperty(mockMutator, 'canUndo', {
        get: jest.fn(() => false),
        configurable: true,
    })
    Object.defineProperty(mockMutator, 'canRedo', {
        get: jest.fn(() => false),
        configurable: true,
    })
    Object.defineProperty(mockMutator, 'undoDescription', {
        get: jest.fn(() => ''),
        configurable: true,
    })
    Object.defineProperty(mockMutator, 'redoDescription', {
        get: jest.fn(() => ''),
        configurable: true,
    })
    return {
        __esModule: true,
        default: mockMutator,
    }
})

jest.mock('../../components/flashMessages', () => ({
    sendFlashMessage: jest.fn(),
}))

describe('pages/boardPage/undoRedoHotKeys', () => {
    beforeEach(() => {
        jest.clearAllMocks()
    })

    test('renders null', () => {
        const {container} = render(wrapIntl(<UndoRedoHotKeys/>))
        expect(container.firstChild).toBeNull()
    })

    test('undo triggers with Nothing to Undo when canUndo is false', () => {
        jest.spyOn(mutator, 'canUndo', 'get').mockReturnValue(false)

        render(wrapIntl(<UndoRedoHotKeys/>))
        expect(registeredHotkeys['ctrl+z,cmd+z']).toBeDefined()

        registeredHotkeys['ctrl+z,cmd+z']()

        expect(sendFlashMessage).toHaveBeenCalledWith({
            content: 'Nothing to Undo',
            severity: 'low',
        })
    })

    test('undo triggers mutator.undo when canUndo is true without description', async () => {
        jest.spyOn(mutator, 'canUndo', 'get').mockReturnValue(true)
        jest.spyOn(mutator, 'undoDescription', 'get').mockReturnValue('')

        render(wrapIntl(<UndoRedoHotKeys/>))
        registeredHotkeys['ctrl+z,cmd+z']()

        expect(mutator.undo).toHaveBeenCalled()

        // Wait for mutator.undo promise to resolve
        await new Promise(process.nextTick)

        expect(sendFlashMessage).toHaveBeenCalledWith({
            content: 'Undo',
            severity: 'low',
        })
    })

    test('undo triggers mutator.undo when canUndo is true with description', async () => {
        jest.spyOn(mutator, 'canUndo', 'get').mockReturnValue(true)
        jest.spyOn(mutator, 'undoDescription', 'get').mockReturnValue('Card deletion')

        render(wrapIntl(<UndoRedoHotKeys/>))
        registeredHotkeys['ctrl+z,cmd+z']()

        expect(mutator.undo).toHaveBeenCalled()
        await new Promise(process.nextTick)

        expect(sendFlashMessage).toHaveBeenCalledWith({
            content: 'Undo Card deletion',
            severity: 'low',
        })
    })

    test('redo triggers with Nothing to Redo when canRedo is false', () => {
        jest.spyOn(mutator, 'canRedo', 'get').mockReturnValue(false)

        render(wrapIntl(<UndoRedoHotKeys/>))
        expect(registeredHotkeys['shift+ctrl+z,shift+cmd+z']).toBeDefined()

        registeredHotkeys['shift+ctrl+z,shift+cmd+z']()

        expect(sendFlashMessage).toHaveBeenCalledWith({
            content: 'Nothing to Redo',
            severity: 'low',
        })
    })

    test('redo triggers mutator.redo when canRedo is true without description', async () => {
        jest.spyOn(mutator, 'canRedo', 'get').mockReturnValue(true)
        jest.spyOn(mutator, 'redoDescription', 'get').mockReturnValue('')

        render(wrapIntl(<UndoRedoHotKeys/>))
        registeredHotkeys['shift+ctrl+z,shift+cmd+z']()

        expect(mutator.redo).toHaveBeenCalled()
        await new Promise(process.nextTick)

        expect(sendFlashMessage).toHaveBeenCalledWith({
            content: 'Redo',
            severity: 'low',
        })
    })

    test('redo triggers mutator.redo when canRedo is true with description', async () => {
        jest.spyOn(mutator, 'canRedo', 'get').mockReturnValue(true)
        jest.spyOn(mutator, 'redoDescription', 'get').mockReturnValue('Card deletion')

        render(wrapIntl(<UndoRedoHotKeys/>))
        registeredHotkeys['shift+ctrl+z,shift+cmd+z']()

        expect(mutator.redo).toHaveBeenCalled()
        await new Promise(process.nextTick)

        expect(sendFlashMessage).toHaveBeenCalledWith({
            content: 'Redo Card deletion',
            severity: 'low',
        })
    })
})
