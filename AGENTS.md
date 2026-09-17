# AI instructions: applications built with goserver

This file applies to the entire repository. Follow it when creating or changing a website, an API, or an application combining both. Deliver a complete, understandable, secure, and fast product that uses goserver capabilities while preserving the simplicity of ordinary Go.

## 1. Philosophy and decision criteria

- Use goserver as the application foundation: routing, middleware, lifecycle, responses, web runtime, and required integrations. Check built-in capabilities before selecting an additional library.
- Stay close to standard Go and `net/http`. Do not introduce a custom HTTP context, a second router, a parallel middleware stack, or a generic service/repository hierarchy without a concrete need.
- GOSH is the native path for web applications: SSR in Go, file-based pages, components, and progressively added client interactivity. Do not replace it with an SPA framework unless requested or required by the existing project.
- Web and API share one listener. The optional web layer attaches to the existing server; it does not need a separate port.
- Write application code and necessary overrides only. Do not copy the framework into the project, modify the Go module cache, or create empty packages for hypothetical future use.
- Prefer small, clearly named functions and packages, bounded concurrency, and understandable dependencies. An abstraction should solve an existing problem.
- A working feature includes error handling, access checks, and user-facing states. Do not leave nonfunctional buttons, API stubs, or demonstration content in the finished product.

## 2. Establish the facts first

Before making changes, read `README.md`, `go.mod`, `taskfile.yml`, `cmd/main.go`, and the relevant source and configuration files. Check `git status` and preserve the user's existing changes.

The API source of truth is **the goserver version selected by this project**. Use `go list -m -json github.com/myelophone/goserver` to locate the dependency; read the documentation and source in its reported `Dir`. Verify actual signatures and method availability instead of guessing from names or Nuxt/Vue conventions.

A neighboring goserver checkout and the main-branch README may be newer than the selected dependency. Changes there do not update the application. Do not add a permanent `replace` or `go.work`, or upgrade the dependency without a task-related need. If a new capability is required, select a suitable version explicitly, update `go.mod` and `go.sum` together, and verify compatibility.

This template initially disables `runtime.enabled`. Its Task commands invoke web tooling even in API mode; Go, Node.js, and Task are required for the standard run/build workflow. Do not copy commands from the framework's Taskfile without checking the application's Taskfile.

## 3. Where code belongs

| Responsibility                          | Location and approach                                                         |
| --------------------------------------- | ----------------------------------------------------------------------------- |
| Routes, service registration, and hooks | New files in `app/`, registered through `goserver.RegisterHook`               |
| Business logic                          | `app/` for simple features; domain packages in `internal/` as needed          |
| Application startup                     | `cmd/main.go`: server assembly and startup only                               |
| Pages and layouts                       | `web/pages/*.gosh`, `web/layouts/*.gosh`                                      |
| Reusable UI                             | `web/components/`                                                             |
| Go logic for pages and components       | `web/logic/`, importing `github.com/myelophone/goserver/web/runtime` directly |
| File-based HTTP endpoints               | Actual `web/server/**/*.method.go` handlers                                   |
| Shared styling                          | `web/css/default.css`; component styles in `<style scoped>`                   |
| Public client state and hooks           | `web/stores/`, `web/plugins/`                                                 |
| Markdown content                        | `web/content/`                                                                |
| Public files                            | `assets/`                                                                     |
| Non-secret web settings                 | `websettings.json` and environment overrides                                  |
| Environment and secrets                 | Environment variables; never commit the local `.env`                          |

Preserve the template's bootstrap sequence: `NewServer` → `Defaults` → configuration and application hooks → i18n → `EnableWeb` → `Run`. Do not duplicate `Defaults()` or start another HTTP server. If the task requires a different service or middleware setup order, verify it against the selected version's source.

`app/app.go` is marked as a template file that updates may overwrite. Add application code in neighboring files. Local `web/modules` are not automatically discovered: explicitly import the Go packages you use so their registration executes.

When renaming the project, update its module path and application imports, including the `app` import in `cmd/main.go`. Keep `github.com/myelophone/goserver` imports pointing to the framework.

## 4. API: ordinary Go, explicit contracts

For a standalone API, prefer standard handlers `func(http.ResponseWriter, *http.Request)` and `s.GET`, `s.POST`, other HTTP methods, and `s.Group`. Endpoints belonging to the file-based web layer may use `web/server` and `runtime.Event`. Choose one clear registration mechanism for each endpoint.

Minimal registration in a new `app/routes.go` file:

```go
package app

import (
    "net/http"

    "github.com/myelophone/goserver"
)

func init() {
    goserver.RegisterHook(func(s *goserver.Server) {
        s.GET("/api/hello", func(w http.ResponseWriter, r *http.Request) {
            s.RespondJSON(w, r, map[string]string{"message": "Hello"})
        })
    })
}
```

This example assumes an empty `API_PREFIX`. With `API_PREFIX=/api`, register `/hello`: the external prefix is stripped before routing. Do not repeat it in registered routes. For a website exposing APIs under `/api/*`, normally leave the global prefix empty and use an `/api` route group: the global prefix is not a setting scoped only to that API group. Verify external page, endpoint, and operational URLs whenever the prefix changes.

- Explicit routes take priority over the file-based web fallback. Avoid accidentally registering `/` or a catch-all that shadows the website.
- Use `APP_ERROR_MODE=json` for an API-only service. In a combined application, render API errors explicitly through `RenderErrorJSON` while retaining appropriate HTML errors for pages.
- Define methods, request/response shapes, status codes, validation, and authorization rules. Return consistent errors; never expose internal SQL, stack traces, or secrets.
- Use `goserver.ParseRequest` and validation helpers where appropriate. Check errors and typed-getter results; strict DTOs may use a standard JSON decoder with an explicit unknown-field policy.
- Bound request bodies, file sizes, item counts, and pagination. Client validation supplements server validation.
- Check authentication and object-level permissions on every read or mutation. Do not trust user IDs, tenant IDs, roles, or prices supplied by the request.
- Select meaningful statuses for creation, missing resources, invalid input, denied access, conflicts, and overload. Do not return `200` for every outcome.
- Propagate request context and deadlines to database, outbound HTTP, and other operations. Do not continue after request cancellation without an explicitly justified background workflow.
- Never perform mutations through GET. For retryable writes, use opt-in idempotency after authentication, verifying key scope and storage behavior in the selected version.
- Document public contracts and request examples. If the project uses OpenAPI, keep the specification synchronized with changes.

## 5. Web: inherit the framework layer

goserver supplies base pages, components, layouts, CSS, stores, plugins, content, and system runtime. Layers merge by relative file path:

- A new file adds a capability.
- A matching path completely replaces the base file; contents are not merged.
- Removing an override restores the base file.
- A missing or empty local directory continues inheriting the base layer.
- Go packages are not merged. A local handler must compile independently and explicitly import its dependencies.

In development, an existing `web/playground` is the highest-priority local overlay: playground > project > goserver. Its paths mirror `web/`; `web/playground/assets/` overrides/adds public assets. Production builds exclude the playground and resolve project > goserver into the distribution. Do not disable inherited resources when a local directory is empty or missing.

Start by inspecting available components and their props and handlers. Reuse suitable `Ui*`, grid, navigation, image, form, and consent components. Create application-specific components and compositions rather than manually duplicating built-in behavior.

Replace the inherited welcome page with a local `web/pages/index.gosh` when delivering a finished website. Set real branding, SEO, language, favicon, and social image; verify the rendered result as well as the configuration.

Do not create a complete local `web/system` tree. Override an individual system file for a concrete need after reading the original; remember that it replaces the entire file and needs attention during framework updates. When replacing the default layout, retain required consent, shortcuts, and other enabled capabilities, or deliberately configure alternatives.

### GOSH and server logic

- A renderable `.gosh` file is a UTF-8 SFC with exactly one `<template>`. HTML belongs inside it; `<head>` is supported separately.
- Place directives in separate HTML comments at the beginning: `<!-- @layout default -->`, `<!-- @server ./web/logic/catalog.go#CatalogPage -->`. `@layout none` disables the layout; omitting the directive uses the configured default.
- Use actual GOSH and existing-component syntax, including `m-if` and bindings. Do not automatically transfer Vue/Nuxt directives or composables.
- `index.gosh` maps to its directory, `[id].gosh` to a parameter, and `[...all].gosh` to a catch-all. For file-based endpoints, dynamic segments belong in the filename, not a Go import directory.
- In `web/logic`, use `package logic` and the `runtime.Handler` contract: `Render`/`Action`; `runtime.Noop` can supply unused methods. Verify the export referenced by `@server` and regenerate bindings.
- Ordinary `<script>` blocks are forbidden in `.server.gosh`; `.client.gosh` denotes a client component. Use client islands and lazy hydration for concrete needs.
- `web/global/head`, `scripts`, and `styles` contain injection fragments rather than ordinary SFCs. Do not automatically wrap them in `<template>`.
- Make public pages useful on initial SSR. Do not fetch server-available content again in the browser solely to display the initial page.

### Navigation and language switching

- Preserve GOSH runtime navigation for internal same-origin pages, including language-switching links. Use real `<a href="...">` links with valid destination URLs so direct access and native navigation also work.
- **Do not add `data-runtime-off` to ordinary internal links or language links.** Its presence tells the runtime to skip interception and leaves navigation to the browser, causing a full document load. Even `data-runtime-off="false"` opts out because the runtime checks attribute presence, not its value.
- Use `data-runtime-off` only when a flow explicitly requires native browser navigation or submission. Explain the reason near the code. Do not add it as a workaround for stale translations, missing handlers, or lifecycle bugs; fix the underlying problem.
- Do not use `window.location.href`, `location.assign`, `location.replace`, or `location.reload()` for ordinary internal navigation or locale changes. For programmatic navigation, verify and use the documented runtime API in the selected version.
- Inspect built-in `UiLangLink`, `UiLanguageSelect`, and the `ui` store before composing a locale switcher. Do not assume a locale-aware link switches language, or that updating store language navigates to and rerenders a translated page; verify each component's actual behavior.
- Generate locale URLs from the configured locales and default locale. Preserve the equivalent page, relevant query parameters, and fragment where appropriate; avoid duplicate locale prefixes and hard-coded default-language assumptions.
- Verify that locale navigation updates page content, `<html lang>`, title/meta/canonical, and the switcher's active state. Updating only client state or `<html lang>` is insufficient when server-rendered content depends on locale.
- For an unexpected reload, inspect the rendered link/form attributes, click handlers, destination origin and target, runtime initialization, console errors, redirects, and runtime fallback behavior. A successful destination page alone does not prove SPA navigation worked.
- Keep intentional browser behaviors working: downloads, external links, new-tab/modifier clicks, and same-page anchors. Do not force every link through custom click handlers.

### State, forms, and client lifecycle

- SSR props are not serialized to the client: client-module `ctx.props` is empty by design. Expose intentionally public state through `runtime.UseStoreState` and stores; keep secrets and authoritative state on the server.
- Runtime continuation tokens continue rendering; they do not replace authentication or authorization. Check access and submitted input in every server action.
- GOSH action forms use `data-gosh-form="event"` without an HTML `action`. Read fields through `props.Form()`, separately from SSR props. Choose an ordinary HTML workflow explicitly for native forms targeting an API endpoint.
- Send files to a separate protected upload endpoint: runtime action forms do not transmit file inputs.
- Use documented `_gosh` hooks, forms, queries, and action APIs. Do not modify framework-owned `/_gosh/*`, opaque tokens, or internal `data-gosh-*` markup.
- During SPA navigation, release listeners, observers, timers, and subscriptions according to runtime lifecycle. Verify page re-entry and absence of duplicate registration.
- Async UI needs pending, empty, error, and success states. Optimistic updates must roll back on failure and reconcile with the server result.

### Visual quality, accessibility, and SEO

- Build a coherent design: typography, spacing rhythm, clear hierarchy, and consistent colors and components. Use existing theme tokens, including `--ui-*`, and verify supported themes.
- Put shared CSS in `web/css/default.css` and local styles in `<style scoped>`. Place required Tailwind directives in the designated `web/system/tailwind/global.css`; account for complete replacement of the base file.
- Use semantic HTML, labels, real buttons and links, keyboard navigation, visible focus, sufficient contrast, and understandable errors. Check modal focus behavior and reduced motion.
- Verify mobile and desktop layouts, long text, empty data, errors, and images. Do not conceal overflow problems with arbitrary global CSS.
- Give images appropriate alt text, dimensions, and responsive delivery. Do not lazily load critical above-the-fold images without a reason; defer noncritical resources.
- Configure real page SEO through web settings and `runtime.UseSeo`, canonical URLs, and indexing policy. Keep private pages out of the public search index and sitemap.
- Use canonical paths without trailing slashes; apply `goserver.CanonicalURL` when generating links in Go.
- Integrate third-party trackers and embeds according to the project's consent policy, using built-in consent components. Declining consent must preserve core website functionality.

## 6. Security and operations

- Preserve protective defaults. Do not globally disable CSRF, security headers, rate limits, or body limits for one test or endpoint. Make targeted policy changes and explain the reason.
- Distinguish browser cookie authentication from API bearer authentication; configure precise trusted origins and a separate CORS policy when needed. `CSRF_TRUSTED_ORIGINS` is not a CORS setting. Do not use wildcard origins with credentials.
- Sessions and JWT do not constitute a complete user system. Implement token/session verification, expiration, and required revocation rules; hash passwords with an appropriate password-hashing library.
- Before production, replace example secrets with externally supplied, random, stable values. Never expose secrets, bearer tokens, cookies, or personal data in logs or responses.
- Check the selected version's limitations around `SetCookie`, WebSocket origin policy, forwarding headers, and redirect Host handling. With a reverse proxy, validate Host and remove untrusted forwarding headers at the edge.
- WebSockets and SSE require authentication, origin policy, resource bounds, cancellation, and compatible timeouts. Do not assume defaults support an indefinite stream.
- Use public `/healthz` for health checks. Restrict detailed diagnostics, `/metricz`, and pprof using the built-in secret and network policy.
- Use request IDs and the built-in logger; provide enough diagnostic information while preserving privacy. Do not expose detailed internal errors to users.
- For PostgreSQL, use parameterized queries, context, and transactions for related changes. Plan migrations explicitly. Calculate the connection-pool budget across all replicas.
- Add DB/Redis only when needed. In-memory sessions, cache, and runtime state do not suit every deployment: with multiple replicas, verify shared storage or the required sticky-session strategy and interface support in the selected version.
- Background work needs bounded concurrency, error handling, and shutdown. Use framework lifecycle/hooks and tracked jobs instead of uncontrolled goroutines.

## 7. Caching and performance

- Establish correctness first and measure bottlenecks. Do not enable arbitrary cache TTLs, prefetch, preload, or performance flags without checking their effect.
- Shared response/render caches are for public data. Do not cache personalized or tenant-dependent responses under a common key; account for identity, locale, tenant, and parameters where needed.
- `routeRules` with `public: true` are allowed only for demonstrably public pages without sessions, identity, or per-request runtime state. Do not apply this to the entire website by default.
- Do not rely solely on automatic cache bypass for cookies/Authorization: verify the actual policy and every other personalization source.
- Use cache tags and revalidation after successfully persisting changes. Invalidation must cover the server/browser caches in use; a failed mutation must not appear successful.
- Do not confuse local caching with distributed coordination. Verify TTL/SWR, idempotency, and invalidation behavior for the selected storage adapter.
- Reduce JS, repeated requests, and unnecessary hydration. Use supported Server-Timing/profile settings for measurements; enable diagnostic headers deliberately.

## 8. Generation, builds, and delivery

- `task run` generates bindings and runs in development; `task dev` generates bindings and starts Air; `task preview` runs source with production settings rather than validating the built distribution.
- `task web:generate` updates bindings; `task build` creates the production distribution; `task server` starts the built executable.
- Do not edit or commit `internal/goservergen/`, `cmd/web_import_gen.go`, `dist/`, or temporary reports. Regenerate after adding/removing Go handlers or changing the module path.
- Always compile the entire `./cmd` package rather than the individual `cmd/main.go` file, otherwise the generated importer is omitted. The CLI uses `-tags webcli`; do not use this tag for the runtime binary.
- Task loads `.env` and selects the environment. Direct Go commands do neither; set the environment explicitly instead of assuming production/development behavior.
- The production distribution contains the executable **and `dist/assets/`**. Ordinary images, fonts, and downloads are not embedded in the binary. The web build embeds resolved production web configuration and `robots.txt` when present.
- Add relative paths for dynamically selected public assets to `build.includeFiles`. Verify actual files in dist; an asset audit does not prove dynamic resources are unused.
- Production public assets resolve beside the executable, independent of the process working directory. Keep the executable and merged `assets/` directory together. `task server` runs from `dist`; also test direct executable startup from the repository and an unrelated directory. Check real browser image loading (`naturalWidth > 0`) and rendered resource URLs, not just files existing on disk or the page returning `200`.
- Production requires neither Go/Node nor the source web tree. Preserve the non-root container, health checks, and controlled port binding. Do not add a toolchain to the final image without a need.
- Production binaries must exclude development and build tooling at compile time: no embedded Yarn distributions, tooling package manifests/lockfiles, installers, Node/esbuild/Tailwind/PostCSS invocation, development source-tree fallback, playground, code generation, or self-build CLI. Setting `APP_ENV=prod` or hiding CLI commands is insufficient. Use the production build tag (`myelophone_prod`) and verify the selected dependency actually excludes these capabilities. The builder prepares the complete distribution; the deployed server only consumes its prepared resources. Do not retain disabled compiler branches or compiler stubs in production.
- This separation must work for every importing Go module, including goserver-template. Development/build tooling comes from the dependency without copying development files or a `web/system` tree into the consumer. Verify both inherited defaults and local overrides, then run the resulting distribution without Go/Node or access to module/toolchain caches.
- Do not automatically prune or delete source assets based on an audit report. Do not update all dependencies through `task update` as an incidental step in an ordinary task.

## 9. Completing a task

Choose validation appropriate to the change; a documentation-only edit does not need a full build. For application changes:

1. Regenerate bindings when GOSH/server handlers or imports change.
2. Format changed Go code; run `task vet`, `task test`, and `task build` for substantial code/runtime changes. If external infrastructure prevents validation, report the specific reason and the checks you could complete.
3. Verify changed HTTP flows: success, invalid input, denied access, not found, and required headers/statuses. Test through the standard middleware stack.
4. For web changes, check initial SSR, direct URLs, refresh, runtime navigation, back/forward, forms, and assets in a browser. For internal links and language switching, use the Network panel to confirm normal clicks do not issue a new top-level Document request when runtime navigation is expected; verify translated SSR content, metadata, and history after switching and returning. Inspect mobile/desktop layouts and console/network errors. Check logs: an isolated component error can be hidden behind HTTP `200`.
5. For deployment changes, run the built distribution in the appropriate environment and verify health, pages/API, and disk assets.
6. Add meaningful tests for changed contracts, permissions, and regressions; avoid tests that merely mirror implementation. Update documentation and `.env.example` for new public settings without including secrets.
7. Review the diff for overwritten user changes, secrets, generated files, accidental dependencies, and unfinished stubs.

In the final response, briefly state what works, which checks were completed, and any remaining limitations. Do not claim passing tests, visual verification, or production readiness unless the relevant checks were performed.

**A result follows goserver's philosophy when the application remains understandable ordinary Go, the web uses SSR and necessary GOSH interactivity, the framework is inherited as a dependency, and security, errors, and delivery are verified alongside the main feature.**
