import { useEffect, type FunctionComponent } from 'react'

import { useNavigate } from 'react-router-dom'

import { lazyComponent } from '@sourcegraph/shared/src/util/lazyComponent'

import { urlToSavedSearchesList } from '../savedSearches/ListPage'

import type { NamespaceProps } from '.'
import type { NamespaceAreaRoute } from './NamespaceArea'

const WorkflowArea = lazyComponent(() => import('../workflows/WorkflowArea'), 'WorkflowArea')

export const namespaceAreaRoutes: readonly NamespaceAreaRoute[] = [
    {
        path: 'searches/*',
        render: props => <SavedSearchesRedirect {...props} />,
        condition: () => window.context?.codeSearchEnabledOnInstance,
    },
    {
        path: 'workflows/*',
        render: props => <WorkflowArea {...props} />,
        condition: () => window.context?.codyEnabledForCurrentUser,
    },
]

/**
 * Redirect from `/users/USER/searches` and `/orgs/ORG/searches` to the new global URL path
 * `/searches?owner=OWNER`, for backcompat.
 */
const SavedSearchesRedirect: FunctionComponent<NamespaceProps> = ({ namespace }) => {
    const navigate = useNavigate()
    useEffect(() => {
        navigate(urlToSavedSearchesList(namespace.id))
    }, [navigate, namespace.id])
    return null
}
