import { useEffect, type FunctionComponent } from 'react'

import { mdiMagnify } from '@mdi/js'
import classNames from 'classnames'
import { useParams } from 'react-router-dom'

import { SyntaxHighlightedSearchQuery } from '@sourcegraph/branded'
import { useQuery } from '@sourcegraph/http-client'
import { useSettingsCascade } from '@sourcegraph/shared/src/settings/settings'
import type { TelemetryV2Props } from '@sourcegraph/shared/src/telemetry'
import { buildSearchURLQuery } from '@sourcegraph/shared/src/util/url'
import { Alert, Button, Container, H3, Icon, Link, LoadingSpinner, PageHeader } from '@sourcegraph/wildcard'

import type {
    SavedSearchFields,
    SavedSearchResult,
    SavedSearchVariables,
    SearchPatternType,
} from '../graphql-operations'
import { namespaceTelemetryMetadata } from '../namespaces/telemetry'
import { defaultPatternTypeFromSettings } from '../util/settings'

import { savedSearchQuery } from './graphql'
import { SavedSearchPage } from './Page'
import { telemetryRecordSavedSearchViewSearchResults } from './telemetry'

import styles from './DetailPage.module.scss'

/**
 * Page to show a saved search.
 */
export const DetailPage: FunctionComponent<TelemetryV2Props> = ({ telemetryRecorder }) => {
    const { id } = useParams<{ id: string }>()

    const result = useQuery<SavedSearchResult, SavedSearchVariables>(savedSearchQuery, { variables: { id: id! } })
    const savedSearch = result.data?.node?.__typename === 'SavedSearch' ? result.data.node : null

    return (
        <SavedSearchPage
            title={savedSearch ? `${savedSearch.description} - saved search` : 'Saved search'}
            actions={
                savedSearch?.viewerCanAdminister && (
                    <Button to={`${savedSearch.url}/edit`} variant="secondary" as={Link}>
                        Edit
                    </Button>
                )
            }
            breadcrumbs={savedSearch ? <PageHeader.Breadcrumb>{savedSearch.description}</PageHeader.Breadcrumb> : null}
        >
            {result.loading ? (
                <LoadingSpinner />
            ) : !savedSearch ? (
                <Alert variant="danger" as="p">
                    Saved search not found.
                </Alert>
            ) : (
                <Detail savedSearch={savedSearch} telemetryRecorder={telemetryRecorder} />
            )}
        </SavedSearchPage>
    )
}

const Detail: FunctionComponent<TelemetryV2Props & { savedSearch: SavedSearchFields }> = ({
    savedSearch,
    telemetryRecorder,
}) => {
    useEffect(() => {
        telemetryRecorder.recordEvent('savedSearches.detail', 'view', {
            metadata: namespaceTelemetryMetadata(savedSearch.owner),
        })
    }, [telemetryRecorder, savedSearch.owner])

    const defaultPatternType: SearchPatternType = defaultPatternTypeFromSettings(useSettingsCascade())
    const searchURL = `/search?${buildSearchURLQuery(savedSearch.query, defaultPatternType, false)}`
    return (
        <Container className={classNames(styles.container)}>
            <Button
                variant="primary"
                size="lg"
                to={searchURL}
                as={Link}
                onClick={() => telemetryRecordSavedSearchViewSearchResults(telemetryRecorder, savedSearch, 'Detail')}
            >
                <Icon aria-hidden={true} svgPath={mdiMagnify} className="flex-shrink-0" size="sm" /> Run search
            </Button>
            <div className="d-flex flex-column flex-gap-2 align-items-center">
                <H3>{savedSearch.description}</H3>
                <SyntaxHighlightedSearchQuery query={savedSearch.query} />
            </div>
        </Container>
    )
}
