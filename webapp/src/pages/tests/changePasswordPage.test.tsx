// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {render, screen, fireEvent, waitFor} from '@testing-library/react'
import {Router} from 'react-router-dom'
import {createMemoryHistory} from 'history'
import {Provider as ReduxProvider} from 'react-redux'
import configureStore from 'redux-mock-store'
import thunk from 'redux-thunk'
import {mocked} from 'jest-mock'

import {wrapIntl} from '../../testUtils'
import client from '../../octoClient'
import ChangePasswordPage from '../changePasswordPage'

jest.mock('../../octoClient')
const mockedOctoClient = mocked(client, true)

describe('pages/changePasswordPage', () => {
    const mockStore = configureStore([thunk])
    let history = createMemoryHistory()

    beforeEach(() => {
        jest.resetAllMocks()
        history = createMemoryHistory()
    })

    test('renders log in first message if no user is logged in', () => {
        const store = mockStore({
            users: {
                me: null,
            },
        })

        render(
            <ReduxProvider store={store}>
                {wrapIntl(
                    <Router history={history}>
                        <ChangePasswordPage/>
                    </Router>,
                )}
            </ReduxProvider>,
        )

        expect(screen.getByText('Change Password')).toBeDefined()
        expect(screen.getByText('Log in first')).toBeDefined()
    })

    test('renders change password form when user is logged in', () => {
        const store = mockStore({
            users: {
                me: {
                    id: 'user_id_1',
                    username: 'testuser',
                },
            },
        })

        render(
            <ReduxProvider store={store}>
                {wrapIntl(
                    <Router history={history}>
                        <ChangePasswordPage/>
                    </Router>,
                )}
            </ReduxProvider>,
        )

        expect(screen.getByText('Change Password')).toBeDefined()
        expect(screen.getByPlaceholderText('Enter current password')).toBeDefined()
        expect(screen.getByPlaceholderText('Enter new password')).toBeDefined()
        expect(screen.getByText('Change password')).toBeDefined()
        expect(screen.getByText('Cancel')).toBeDefined()
    })

    test('handles input changes', () => {
        const store = mockStore({
            users: {
                me: {
                    id: 'user_id_1',
                    username: 'testuser',
                },
            },
        })

        render(
            <ReduxProvider store={store}>
                {wrapIntl(
                    <Router history={history}>
                        <ChangePasswordPage/>
                    </Router>,
                )}
            </ReduxProvider>,
        )

        const oldPwdInput = screen.getByPlaceholderText('Enter current password') as HTMLInputElement
        const newPwdInput = screen.getByPlaceholderText('Enter new password') as HTMLInputElement

        fireEvent.change(oldPwdInput, {target: {value: 'old_pass'}})
        fireEvent.change(newPwdInput, {target: {value: 'new_pass'}})

        expect(oldPwdInput.value).toBe('old_pass')
        expect(newPwdInput.value).toBe('new_pass')
    })

    test('handles successful password change submission', async () => {
        const store = mockStore({
            users: {
                me: {
                    id: 'user_id_1',
                    username: 'testuser',
                },
            },
        })

        mockedOctoClient.changePassword.mockResolvedValue({
            code: 200,
            json: {},
        })

        render(
            <ReduxProvider store={store}>
                {wrapIntl(
                    <Router history={history}>
                        <ChangePasswordPage/>
                    </Router>,
                )}
            </ReduxProvider>,
        )

        const oldPwdInput = screen.getByPlaceholderText('Enter current password')
        const newPwdInput = screen.getByPlaceholderText('Enter new password')
        const submitBtn = screen.getByText('Change password')

        fireEvent.change(oldPwdInput, {target: {value: 'old_pass'}})
        fireEvent.change(newPwdInput, {target: {value: 'new_pass'}})
        fireEvent.click(submitBtn)

        await waitFor(() => {
            expect(mockedOctoClient.changePassword).toHaveBeenCalledWith('user_id_1', 'old_pass', 'new_pass')
            expect(screen.getByText('Password changed, click to continue.')).toBeDefined()
        })
    })

    test('handles failed password change submission', async () => {
        const store = mockStore({
            users: {
                me: {
                    id: 'user_id_1',
                    username: 'testuser',
                },
            },
        })

        mockedOctoClient.changePassword.mockResolvedValue({
            code: 400,
            json: {error: 'Incorrect current password'},
        })

        render(
            <ReduxProvider store={store}>
                {wrapIntl(
                    <Router history={history}>
                        <ChangePasswordPage/>
                    </Router>,
                )}
            </ReduxProvider>,
        )

        const oldPwdInput = screen.getByPlaceholderText('Enter current password')
        const newPwdInput = screen.getByPlaceholderText('Enter new password')
        const submitBtn = screen.getByText('Change password')

        fireEvent.change(oldPwdInput, {target: {value: 'old_pass'}})
        fireEvent.change(newPwdInput, {target: {value: 'new_pass'}})
        fireEvent.click(submitBtn)

        await waitFor(() => {
            expect(mockedOctoClient.changePassword).toHaveBeenCalledWith('user_id_1', 'old_pass', 'new_pass')
            expect(screen.getByText('Change password failed: Incorrect current password')).toBeDefined()
        })

        // typing clears error message
        fireEvent.change(oldPwdInput, {target: {value: 'another'}})
        expect(screen.queryByText('Change password failed: Incorrect current password')).toBeNull()
    })
})
