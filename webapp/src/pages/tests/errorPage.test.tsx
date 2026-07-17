// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {render, screen, fireEvent, waitFor} from '@testing-library/react'
import {Router} from 'react-router-dom'
import {createMemoryHistory} from 'history'

import {wrapIntl} from '../../testUtils'
import ErrorPage from '../errorPage'
import {ErrorId} from '../../errors'

describe('pages/errorPage', () => {
    test('renders sorry something went wrong by default', () => {
        const history = createMemoryHistory()
        history.location.search = ''
        render(
            wrapIntl(
                <Router history={history}>
                    <ErrorPage/>
                </Router>,
            ),
        )
        expect(screen.getByText('Sorry, something went wrong')).toBeDefined()
    })

    test('renders team-undefined error page and button click pushes redirect', async () => {
        const history = createMemoryHistory()
        history.push = jest.fn()
        history.location.search = `?id=${ErrorId.TeamUndefined}`

        const originalLocation = window.location
        delete (window as any).location;
        (window as any).location = {origin: 'http://test.com', href: ''}

        render(
            wrapIntl(
                <Router history={history}>
                    <ErrorPage/>
                </Router>,
            ),
        )

        expect(screen.getByText('Not a valid team.')).toBeDefined()
        const button = screen.getByText('Back to Home')
        expect(button).toBeDefined()
        fireEvent.click(button)

        await waitFor(() => {
            expect(window.location.href).toBe('http://test.com')
        })

        window.location = originalLocation
    })

    test('renders board-not-found error page and handles button redirect', () => {
        const history = createMemoryHistory()
        history.push = jest.fn()
        history.location.search = `?id=${ErrorId.BoardNotFound}`

        render(
            wrapIntl(
                <Router history={history}>
                    <ErrorPage/>
                </Router>,
            ),
        )

        expect(screen.getByText('Board not found.')).toBeDefined()
        const button = screen.getByText('Back to team')
        expect(button).toBeDefined()
        fireEvent.click(button)

        expect(history.push).toHaveBeenCalledWith('/')
    })

    test('handles path redirect correctly if redirect is a direct string', async () => {
        const history = createMemoryHistory()
        history.location.search = `?id=${ErrorId.InvalidReadOnlyBoard}`

        const originalLocation = window.location
        delete (window as any).location;
        (window as any).location = {origin: 'http://test.com', href: ''}

        render(
            wrapIntl(
                <Router history={history}>
                    <ErrorPage/>
                </Router>,
            ),
        )

        expect(screen.getByText("You don't have access to this board. Log in to access Boards.")).toBeDefined()
        const button = screen.getByText('Log in')
        expect(button).toBeDefined()
        fireEvent.click(button)

        await waitFor(() => {
            expect(window.location.href).toBe('http://test.com')
        })

        window.location = originalLocation
    })

    test('handles not-logged-in error and does immediate redirect', () => {
        const history = createMemoryHistory()
        history.push = jest.fn()
        history.location.search = `?id=${ErrorId.NotLoggedIn}`

        render(
            wrapIntl(
                <Router history={history}>
                    <ErrorPage/>
                </Router>,
            ),
        )

        expect(history.push).toHaveBeenCalledWith('/login')
    })
})
