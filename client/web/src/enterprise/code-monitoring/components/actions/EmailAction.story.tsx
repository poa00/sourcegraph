import { storiesOf } from '@storybook/react'
import React from 'react'
import sinon from 'sinon'

import { WebStory } from '../../../../components/WebStory'
import { mockAuthenticatedUser } from '../../testing/util'

import { EmailAction, EmailActionProps } from './EmailAction'

const { add } = storiesOf('web/enterprise/code-monitoring/actions/EmailAction', module)

const defaultProps: EmailActionProps = {
    action: undefined,
    setAction: sinon.fake(),
    actionCompleted: false,
    setActionCompleted: sinon.fake(),
    disabled: false,
    authenticatedUser: mockAuthenticatedUser,
    monitorName: 'Example code monitor',
}

const action: EmailActionProps['action'] = {
    __typename: 'MonitorEmail',
    id: 'id1',
    recipients: { nodes: [{ id: 'userID' }] },
    enabled: true,
}

add('disabled', () => <WebStory>{() => <EmailAction {...defaultProps} disabled={true} />}</WebStory>)

add('closed, not populated', () => <WebStory>{() => <EmailAction {...defaultProps} />}</WebStory>)

add('open, not populated', () => <WebStory>{() => <EmailAction {...defaultProps} _testStartOpen={true} />}</WebStory>)

add('closed, populated, enabled', () => (
    <WebStory>{() => <EmailAction {...defaultProps} action={action} actionCompleted={true} />}</WebStory>
))

add('open, populated, enabled', () => (
    <WebStory>
        {() => <EmailAction {...defaultProps} _testStartOpen={true} action={action} actionCompleted={true} />}
    </WebStory>
))

add('closed, populated, disabled', () => (
    <WebStory>
        {() => <EmailAction {...defaultProps} action={{ ...action, enabled: false }} actionCompleted={true} />}
    </WebStory>
))

add('open, populated, disabled', () => (
    <WebStory>
        {() => (
            <EmailAction
                {...defaultProps}
                _testStartOpen={true}
                action={{ ...action, enabled: false }}
                actionCompleted={true}
            />
        )}
    </WebStory>
))
