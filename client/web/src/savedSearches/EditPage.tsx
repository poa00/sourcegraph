import { useCallback, useEffect, useState, type FormEventHandler, type FunctionComponent } from 'react'

import { mdiLink } from '@mdi/js'
import { useLocation, useNavigate, useParams } from 'react-router-dom'

import { logger } from '@sourcegraph/common'
import { useMutation, useQuery } from '@sourcegraph/http-client'
import type { TelemetryRecorder, TelemetryV2Props } from '@sourcegraph/shared/src/telemetry'
import {
    Alert,
    Button,
    ErrorAlert,
    Form,
    H3,
    Icon,
    Link,
    LoadingSpinner,
    Modal,
    PageHeader,
} from '@sourcegraph/wildcard'

import type {
    SavedSearchFields,
    SavedSearchResult,
    SavedSearchVariables,
    TransferSavedSearchOwnershipResult,
    TransferSavedSearchOwnershipVariables,
    UpdateSavedSearchResult,
    UpdateSavedSearchVariables,
} from '../graphql-operations'
import { NamespaceSelector } from '../namespaces/NamespaceSelector'
import { namespaceTelemetryMetadata } from '../namespaces/telemetry'
import { useAffiliatedNamespaces } from '../namespaces/useAffiliatedNamespaces'
import { PageRoutes } from '../routes.constants'

import { SavedSearchForm, type SavedSearchFormValue } from './Form'
import {
    deleteSavedSearchMutation,
    savedSearchQuery,
    transferSavedSearchOwnershipMutation,
    updateSavedSearchMutation,
} from './graphql'
import { SavedSearchPage } from './Page'

/**
 * Page to edit a saved search.
 */
export const EditPage: FunctionComponent<TelemetryV2Props & { isSourcegraphDotCom: boolean }> = ({
    telemetryRecorder,
    isSourcegraphDotCom,
}) => {
    const { id } = useParams<{ id: string }>()
    const result = useQuery<SavedSearchResult, SavedSearchVariables>(savedSearchQuery, { variables: { id: id! } })
    const savedSearch = result.data?.node?.__typename === 'SavedSearch' ? result.data.node : null

    return (
        <SavedSearchPage
            title={savedSearch ? `Editing: ${savedSearch.description} - saved search` : 'Edit saved search'}
            actions={
                savedSearch && (
                    <Button to={savedSearch.url} variant="secondary" as={Link}>
                        <Icon aria-hidden={true} svgPath={mdiLink} /> Permalink
                    </Button>
                )
            }
            breadcrumbs={<PageHeader.Breadcrumb>Edit</PageHeader.Breadcrumb>}
        >
            {result.loading ? (
                <LoadingSpinner />
            ) : !savedSearch ? (
                <Alert variant="danger" as="p">
                    Saved search not found.
                </Alert>
            ) : (
                <EditForm
                    savedSearch={savedSearch}
                    isSourcegraphDotCom={isSourcegraphDotCom}
                    telemetryRecorder={telemetryRecorder}
                />
            )}
        </SavedSearchPage>
    )
}

/**
 * Form to edit a saved search.
 */
const EditForm: FunctionComponent<
    TelemetryV2Props & { savedSearch: SavedSearchFields; isSourcegraphDotCom: boolean }
> = ({ savedSearch, telemetryRecorder, isSourcegraphDotCom }) => {
    useEffect(() => {
        telemetryRecorder.recordEvent('savedSearches.update', 'view', {
            metadata: namespaceTelemetryMetadata(savedSearch.owner),
        })
    }, [telemetryRecorder, savedSearch.owner])

    const [updateSavedSearch, { loading: updateLoading, error: updateError }] = useMutation<
        UpdateSavedSearchResult,
        UpdateSavedSearchVariables
    >(updateSavedSearchMutation, {})
    const [flashUpdated, setFlashUpdated] = useState(false)

    const navigate = useNavigate()
    const onSubmit = useCallback(
        async (fields: SavedSearchFormValue): Promise<void> => {
            try {
                await updateSavedSearch({
                    variables: {
                        id: savedSearch.id,
                        input: {
                            description: fields.description,
                            query: fields.query,
                        },
                    },
                })
                telemetryRecorder.recordEvent('savedSearches', 'update', {
                    metadata: namespaceTelemetryMetadata(savedSearch.owner),
                })
                setFlashUpdated(true)
                setTimeout(() => {
                    setFlashUpdated(false)
                }, 1000)
            } catch {
                // Mutation error is read in useMutation call.
            }
        },
        [savedSearch.id, savedSearch.owner, telemetryRecorder, updateSavedSearch]
    )

    const [showTransferOwnershipModal, setShowTransferOwnershipModal] = useState(false)

    const [deleteSavedSearch, { loading: deleteLoading, error: deleteError }] = useMutation(deleteSavedSearchMutation)
    const onDeleteClick = useCallback(async (): Promise<void> => {
        if (!savedSearch) {
            return
        }
        if (!window.confirm(`Delete the saved search ${JSON.stringify(savedSearch.description)}?`)) {
            return
        }
        try {
            await deleteSavedSearch({ variables: { id: savedSearch.id } })
            telemetryRecorder.recordEvent('savedSearches', 'delete', {
                metadata: namespaceTelemetryMetadata(savedSearch.owner),
            })
            navigate(PageRoutes.SavedSearches)
        } catch (error) {
            logger.error(error)
        }
    }, [savedSearch, deleteSavedSearch, telemetryRecorder, navigate])

    const location = useLocation()

    // Flash after transferring ownership.
    const justTransferredOwnership = !!location.state?.[TRANSFERRED_OWNERSHIP_LOCATION_STATE_KEY]
    useEffect(() => {
        if (justTransferredOwnership) {
            setTimeout(() => navigate({}, { state: {} }), 1000)
        }
    }, [justTransferredOwnership, navigate])

    return (
        <>
            <SavedSearchForm
                submitLabel="Save"
                onSubmit={onSubmit}
                otherButtons={
                    <>
                        <div className="flex-grow-1" />
                        {savedSearch.viewerCanAdminister && (
                            <Button
                                onClick={() => {
                                    telemetryRecorder.recordEvent('savedSearches.transferOwnership', 'openModal', {
                                        metadata: namespaceTelemetryMetadata(savedSearch.owner),
                                    })
                                    setShowTransferOwnershipModal(true)
                                }}
                                disabled={updateLoading || deleteLoading}
                                variant="secondary"
                            >
                                Transfer ownership
                            </Button>
                        )}
                        {savedSearch.viewerCanAdminister && (
                            <Button
                                onClick={onDeleteClick}
                                disabled={updateLoading || deleteLoading}
                                variant="danger"
                                outline={true}
                            >
                                Delete
                            </Button>
                        )}
                    </>
                }
                isSourcegraphDotCom={isSourcegraphDotCom}
                initialValue={savedSearch}
                loading={updateLoading || deleteLoading}
                error={updateError ?? deleteError}
                flash={flashUpdated ? 'Saved.' : justTransferredOwnership ? 'Transferred ownership.' : undefined}
                telemetryRecorder={telemetryRecorder}
                beforeFields={
                    <NamespaceSelector
                        namespaces={[savedSearch.owner]}
                        disabled={true}
                        label="Owner"
                        className="w-fit-content"
                    />
                }
            />
            {showTransferOwnershipModal && (
                <TransferOwnershipModal
                    savedSearch={savedSearch}
                    onClose={() => {
                        setShowTransferOwnershipModal(false)
                        telemetryRecorder.recordEvent('savedSearches.transferOwnership', 'closeModal', {
                            metadata: namespaceTelemetryMetadata(savedSearch.owner),
                        })
                    }}
                    telemetryRecorder={telemetryRecorder}
                />
            )}
        </>
    )
}

const TransferOwnershipModal: FunctionComponent<{
    savedSearch: Pick<SavedSearchFields, 'id' | 'owner'>
    onClose: () => void
    telemetryRecorder: TelemetryRecorder
}> = ({ savedSearch, onClose, telemetryRecorder }) => {
    const navigate = useNavigate()

    const { namespaces, loading: namespacesLoading, error: namespacesError } = useAffiliatedNamespaces()
    const validNamespaces = namespaces?.filter(ns => ns.id !== savedSearch.owner.id)
    const [selectedNamespace, setSelectedNamespace] = useState<string | undefined>()
    const selectedNamespaceOrInitial = selectedNamespace ?? validNamespaces?.at(0)?.id

    const [transferSavedSearchOwnership, { loading: transferLoading, error: transferError }] = useMutation<
        TransferSavedSearchOwnershipResult,
        TransferSavedSearchOwnershipVariables
    >(transferSavedSearchOwnershipMutation, {})
    const onSubmit = useCallback<FormEventHandler>(
        async (event): Promise<void> => {
            event.preventDefault()
            try {
                const { data } = await transferSavedSearchOwnership({
                    variables: { id: savedSearch.id, newOwner: selectedNamespaceOrInitial! },
                })
                const updated = data?.transferSavedSearchOwnership
                if (!updated) {
                    return
                }
                telemetryRecorder.recordEvent('savedSearches.transferOwnership', 'submit', {
                    metadata: {
                        [`fromNamespaceType${savedSearch.owner.__typename}`]: 1,
                        [`toNamespaceType${updated.owner.__typename}`]: 1,
                    },
                })
                navigate(`${updated.url}/edit`, { state: { [TRANSFERRED_OWNERSHIP_LOCATION_STATE_KEY]: true } })
                onClose()
            } catch (error) {
                logger.error(error)
            }
        },
        [
            transferSavedSearchOwnership,
            savedSearch.id,
            savedSearch.owner.__typename,
            selectedNamespaceOrInitial,
            telemetryRecorder,
            navigate,
            onClose,
        ]
    )

    const MODAL_LABEL_ID = 'transfer-ownership-modal-label'

    const loading = namespacesLoading || transferLoading

    return (
        <Modal aria-labelledby={MODAL_LABEL_ID} onDismiss={onClose}>
            <Form onSubmit={onSubmit} className="d-flex flex-column flex-gap-4">
                <H3 id={MODAL_LABEL_ID}>Transfer ownership of saved search</H3>
                {namespacesError ? (
                    <ErrorAlert error={namespacesError} />
                ) : loading ? (
                    <LoadingSpinner />
                ) : validNamespaces && validNamespaces.length > 0 && selectedNamespaceOrInitial ? (
                    <>
                        <NamespaceSelector
                            namespaces={validNamespaces}
                            value={selectedNamespaceOrInitial}
                            onSelect={namespace => setSelectedNamespace(namespace)}
                            disabled={transferLoading}
                            loading={namespacesLoading}
                            label="New owner"
                        />
                        <div className="d-flex flex-gap-2">
                            <Button type="submit" disabled={loading} variant="primary">
                                Transfer ownership
                            </Button>
                            <Button onClick={onClose} disabled={loading} variant="secondary" outline={true}>
                                Cancel
                            </Button>
                        </div>
                        {transferError && !loading && <ErrorAlert error={transferError} />}
                    </>
                ) : (
                    <Alert variant="warning">
                        You aren't in any organizations to which you can transfer this saved search.
                    </Alert>
                )}
            </Form>
        </Modal>
    )
}

const TRANSFERRED_OWNERSHIP_LOCATION_STATE_KEY = 'transferredOwnership'
