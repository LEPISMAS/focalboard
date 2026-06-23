// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {render} from '@testing-library/react'
import BackwardCompatibilityQueryParamsRedirect from '../boardPage/backwardCompatibilityQueryParamsRedirect'

describe('pages/boardPage/backwardCompatibilityQueryParamsRedirect', () => {
    test('renders null', () => {
        const {container} = render(<BackwardCompatibilityQueryParamsRedirect />)
        expect(container.firstChild).toBeNull()
    })
})
