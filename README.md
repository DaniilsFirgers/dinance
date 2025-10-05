# Dinance

---

## Monorepo short guideline

### Root tsconfig.json

- Create `tsconfig.json` at the root of the monorepo project. It should contain `references` array like this:

```
    "references": [
        { "path": "./packages/shared" },
        { "path": "./packages/api" },
    ]
```

- It tells "my root project depends on these two subprojects". Typescript will know how to **build them in correct order** (`shared` first, then `api`).
- Each `{ "path": "..." }` is the entry point to another **tsconfig.json** file.

_How it works:_

- When `build` is being done, Typescript goes to `./packages/shared/tsconfig.json` and builds it first.
- Then builds `./packages/api/`, which might depend on it.
- It brings **faster builds**, **type safe imports** and is perfect for **monorepos**.

### Root package.json

- Should include `"files": []`, so that TypeScript does not try to compile everything in the root directory.
- Should include `"private": true`, so that the package cannot be accidentally published to the npm registry, meaning, signals that root package is a **workspace container**.
- Should contain `build` command that builds all workspace packages, typically

```
"scripts": {
  "build": "pnpm -r build"
}

```

- `pnpm -r` means "run recursively" through all _workspace packages_.

- Should include `"workspaces": []` which tells us that all folders inside `/packages` are part of this monorepo - treat them as linked local packages.

Then if `@dinance/api` depends on `@dinance/shared`, `package.json` of api will contain:

```
"dependencies": {
  "@dinance/shared": "workspace:*"
}

```

- It says “Use the local package @dinance/shared from my workspace — not the one from npm.
  And match any version.”. If we specify a version and it is not found locally, it might get fetched from npm registry.
- Pnpm automatically **links** that local dependency instead of fetching from npm.

### Root pnpm-workspace.yaml

- Tells pnpm **exactly where to look for workspace packages**
