import type { FunctionComponent, ReactNode } from 'react'

import { mdiMagnify } from '@mdi/js'
import { useLocation } from 'react-router-dom'

import { PageHeader } from '@sourcegraph/wildcard'

import { Page } from '../components/Page'
import { PageTitle } from '../components/PageTitle'
import { PageRoutes } from '../routes.constants'

import { SavedSearchIcon } from './SavedSearchIcon'

/**
 * The template for a saved search page.
 */
export const SavedSearchPage: FunctionComponent<{
    title?: string
    actions?: ReactNode
    breadcrumbs?: ReactNode
    children: ReactNode
    ['data-testid']?: string
}> = ({ title, actions, breadcrumbs, children, ['data-testid']: dataTestId }) => {
    const location = useLocation()
    const isRootPage = location.pathname === PageRoutes.SavedSearches

    return (
        <Page className="w-100">
            {title && <PageTitle title={title} />}
            <PageHeader actions={actions} className="mb-3" data-testid={dataTestId}>
                <PageHeader.Heading as="h2" styleAs="h1" className="mb-1">
                    <PageHeader.Breadcrumb icon={mdiMagnify} to="/search" aria-label="Code Search" />
                    <PageHeader.Breadcrumb
                        icon={SavedSearchIcon}
                        to={isRootPage ? undefined : PageRoutes.SavedSearches}
                    >
                        Saved Searches
                    </PageHeader.Breadcrumb>
                    {breadcrumbs}
                </PageHeader.Heading>
            </PageHeader>
            {children}
        </Page>
    )
}
