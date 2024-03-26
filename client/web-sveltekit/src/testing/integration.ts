import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import path from 'path'

import { faker } from '@faker-js/faker'
import { test as base, type Page, type Locator } from '@playwright/test'
import glob from 'glob'
import { buildSchema } from 'graphql'
import * as mime from 'mime-types'

import { GraphQLMockServer } from './graphql-mocking'
import type { TypeMocks, ObjectMock, UserMock, OperationMocks } from './graphql-type-mocks'

export { expect, defineConfig, type Locator, type Page, type Browser } from '@playwright/test'

const defaultMocks: TypeMocks = {
    Query: () => ({
        // null means not signed in
        currentUser: null,
    }),
    Person: () => {
        const firstName = faker.person.firstName()
        const lastName = faker.person.lastName()
        return {
            name: `${firstName} ${lastName}`,
            email: faker.internet.email({ firstName, lastName }),
            displayName: faker.internet.userName({ firstName, lastName }),
            avatarURL: null,
        }
    },
    User: () => ({
        avatarURL: null,
    }),
    SettingsCascade: () => ({
        // Ensure this is valid JSON
        final: '{}',
    }),
    TemporarySettings: () => ({
        // Ensure this is valid JSON
        contents: '{}',
    }),
    GitBlob: () => ({
        highlight: {
            // Ensure this is valid JSON
            lsif: '{}',
        },
    }),
    GitRef: () => ({
        url: faker.internet.url(),
    }),
    Signature: () => ({
        date: faker.date.past().toISOString(),
    }),
    GitObjectID: () => faker.git.commitSha(),
    GitCommit: () => ({
        abbreviatedOID: faker.git.commitSha({ length: 7 }),
        subject: faker.git.commitMessage(),
    }),
    JSONCString: () => '{}',
}

const SCHEMA_DIR = path.resolve(
    path.join(path.dirname(fileURLToPath(import.meta.url)), '../../../../cmd/frontend/graphqlbackend')
)
const typeDefs = glob
    .sync('**/*.graphql', { cwd: SCHEMA_DIR })
    .map(file => readFileSync(path.join(SCHEMA_DIR, file), 'utf8'))
    .join('\n')

class Sourcegraph {
    private debugMode = false
    constructor(private readonly page: Page, private readonly graphqlMock: GraphQLMockServer) {}

    async setup(): Promise<void> {
        // All assets are mocked and served from the filesystem. If you do want to use
        // a local preview server or even backend, you can set this env var
        if (!parseBool(process.env.DISABLE_APP_ASSETS_MOCKING)) {
            // routes in playwright are tested in reverse registration order
            // so in order to make this the fallback we register it first
            // all unmatched routes are treated as routes within the application
            // and so only route to the manifest
            await this.page.route('/**/*', route => {
                route.fulfill({
                    status: 200,
                    contentType: 'text/html',
                    body: readFileSync('./build/index.html'),
                })
            })

            // Intercept any asset calls and replace them with static files
            await this.page.route('_app/**/*', route => {
                const asset = new URL(route.request().url()).pathname
                const contentType = mime.contentType(path.basename(asset)) || undefined
                route.fulfill({
                    status: 200,
                    contentType,
                    body: readFileSync(path.join('./build/', new URL(route.request().url()).pathname)),
                    headers: {
                        'cache-control': 'public, max-age=31536000, immutable',
                    },
                })
            })
        }
        // mock graphql calls
        await this.page.route(/\.api\/graphql/, route => {
            const { query, variables, operationName } = JSON.parse(route.request().postData() ?? '')
            const result = this.graphqlMock.query(
                query,
                variables,
                operationName,
                this.debugMode
                    ? {
                          logGraphQLErrors: true,
                          warnOnMissingOperationMocks: true,
                      }
                    : undefined
            )
            route.fulfill({ json: result })
        })
    }

    public debug() {
        this.debugMode = true
    }

    public mockTypes(mocks: TypeMocks): void {
        this.graphqlMock.addTypeMocks(mocks)
    }

    public mockOperations(mocks: OperationMocks): void {
        this.graphqlMock.addOperationMocks(mocks)
    }

    public fixture(fixtures: (ObjectMock & { __typename: NonNullable<ObjectMock['__typename']> })[]): void {
        this.graphqlMock.addFixtures(fixtures)
    }

    public signIn(userMock: UserMock = {}): void {
        this.mockTypes({
            Query: () => ({
                currentUser: {
                    avatarURL: null,
                    ...userMock,
                },
            }),
        })
    }

    public signOut(): void {
        this.mockTypes({
            Query: () => ({
                currentUser: null,
            }),
        })
    }

    public teardown(): void {
        this.graphqlMock.reset()
    }
}

interface Utils {
    scrollYAt(locator: Locator, distance: number): Promise<void>
}

export const test = base.extend<{ sg: Sourcegraph; utils: Utils }, { graphqlMock: GraphQLMockServer }>({
    utils: async ({ page }, use) => {
        use({
            async scrollYAt(locator: Locator, distance: number): Promise<void> {
                // Position mouse over target that wheel events will scrolls the container
                // that contains the target
                const { x, y } = (await locator.boundingBox()) ?? { x: 0, y: 0 }
                await page.mouse.move(x, y)

                // Scroll list, which should load next page
                await page.mouse.wheel(0, distance)
            },
        })
    },
    sg: [
        async ({ page, graphqlMock }, use) => {
            const sg = new Sourcegraph(page, graphqlMock)
            await sg.setup()
            await use(sg)
            sg.teardown()
        },
        { auto: true },
    ],
    graphqlMock: [
        async ({}, use) => {
            const graphqlMock = new GraphQLMockServer({
                schema: buildSchema(typeDefs),
                mocks: defaultMocks,
                typePolicies: {
                    GitBlob: {
                        keyField: 'canonicalURL',
                    },
                    GitTree: {
                        keyField: 'canonicalURL',
                    },
                },
            })
            await use(graphqlMock)
        },
        { scope: 'worker' },
    ],
})

function parseBool(s: string | undefined): boolean {
    if (s === undefined) {
        return false
    }
    return s.toLowerCase() === 'true'
}
