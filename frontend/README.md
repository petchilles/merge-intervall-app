# Frontend of the merge-intervall-app

Provides an intuitive user interface for entering to be merged intervals. Returns merged intervals.
Comes with unit tests for the MergeIntervals component.

## Setup with Docker:

```sh
cd into/merge-intervall-app
docker-compose build
docker-compose up
```

## Setup without Docker: Required for development
## and runnning unit tests (vitest) and e2e tests (playwright)

Note: requires a running golang backend! Please follow the instructions
of the README in the Folder "golang-backend"

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Type-Check, Compile and Minify for Production

```sh
npm run build
```

### Run Unit Tests with [Vitest](https://vitest.dev/)

```sh
npm run test:unit
npm run test:coverage
```

### Run End-to-End Tests with [Playwright](https://playwright.dev)
### (no End-to-End Tests have been implemented as of now)

```sh
# Install browsers for the first run
npx playwright install

# When testing on CI, must build the project first
npm run build

# Runs the end-to-end tests
npm run test:e2e
# Runs the tests only on Chromium
npm run test:e2e -- --project=chromium
# Runs the tests of a specific file
npm run test:e2e -- tests/example.spec.ts
# Runs the tests in debug mode
npm run test:e2e -- --debug
```

### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) to make the TypeScript language service aware of `.vue` types.

## Customize configuration

See [Vite Configuration Reference](https://vitejs.dev/config/).