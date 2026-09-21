# AGENTS.md

Build the requested production application with **goserver**. GOSH is optional and is used only when a web UI is required.

## 1. Rules

- Inspect `git status`, repository docs/config, relevant code, and the goserver version selected by `go.mod`. Locate that dependency with `go list -m -json github.com/myelophone/goserver`; verify APIs in its README/source.
- **Goserver native first.** Use its existing subsystem before writing application infrastructure or adding a dependency. Never add another router, server, SPA/state/CSS/i18n/theme/form/image/cache system.
- Use only documented extension points. Never patch/copy/vendor goserver, modify its module cache, replace a native subsystem, or hide a limitation behind a workaround. Report a missing extension point.
- Never edit generated, compiled, cached, vendored, or tool-owned files. Regenerate them with their owning task. This includes any generated file regardless of path, plus `internal/goservergen/**`, `cmd/web_import_gen.go`, `tmp/goserver/build/**`, and `dist/**`.
- Preserve user changes and the selected dependency. Do not edit `app/app.go`; add neighboring `app/*.go` files.
- No placeholders, dead controls, fake forms, TODO flows, invented facts, copied framework code, or inherited welcome page in a finished site.
- Do not comment obvious code. Comments are only for genuinely non-obvious constraints or decisions.

## 2. Built-in capability index

Verify the selected version before use.

| Area | Goserver capabilities |
| --- | --- |
| HTTP | method routes, params/wildcards, groups/group middleware, standard handlers, all request formats, validation, responses/errors, redirects/rewrites, templates/static assets |
| Middleware | defaults/custom stack, hooks, logging/request IDs, recovery, headers, CSRF, limits, rate limiting, timeouts, load shedding, maintenance |
| Web (optional) | GOSH SSR/file routes, layers/layouts/components, server actions/forms, islands/hydration, navigation/hooks/query/loading/errors/WebSockets |
| Web platform (optional) | grid/UI/consent/SEO/i18n, Tailwind v4, themes/stores, images/content/search/commands/reveal/transitions, route rules/SWR/public-static, cache tags/optimistic actions, preload/prefetch, Early Hints/Server-Timing/Web Vitals |
| Data/state | cookies, sessions, JWT/encryption, cache/idempotency, PostgreSQL helpers, Redis cache/sessions, tenants |
| Realtime | HTML/JSON streaming, SSE, NDJSON, WebSockets |
| Integrations | resilient HTTP client, browser TLS, proxies/Webshare, HTML/soft-404 helpers, cron/jobs, SMTP, Telegram |
| Operations | config/utility API, lifecycle/tracked jobs/shutdown, health, `/metricz`, pprof, logging, Docker, build and asset delivery |

Search goserver's public Go API, examples, `web/components`, `web/stores`, and `web/plugins` before implementing.

## 3. Select the application mode

- **REST/API:** keep the web framework disabled. Do not create/use GOSH pages, layouts, components, stores, plugins, content, or client assets. Register goserver routes/groups and `http.HandlerFunc` handlers in new `app/*.go` files.
- **Web:** enable the goserver web framework and follow section 4.
- **Combined:** use GOSH pages plus explicit goserver API routes on the same listener. Explicit routes have priority over the web fallback.

Keep `cmd/main.go` for assembly. `internal/` is reserved for goserver-generated output, never application code. Startup is `NewServer → Defaults → hooks/config → i18n if used → EnableWeb if configured → Run`.

## 4. Web framework rules

### Project structure and layers

```text
web/pages/       thin route compositions
web/layouts/     shared shell
web/components/ reusable UI and page sections
web/logic/       Go render data/actions
web/server/      optional file endpoints
web/stores/      new public cross-component state
web/plugins/     runtime extensions
web/content/     Markdown content
web/css/         all application global CSS/tokens
web/global/      application head/script inserts
web/modules/     explicitly imported modules
web/teleport/    application teleports
web/tenants/     tenant page/content overrides
assets/          public images/fonts/icons/downloads
```

- Goserver is the base layer. Resolution in development is `playground → project → goserver`; production is `project → goserver` and excludes `web/playground/**`.
- A matching project file fully replaces the lower-layer file; missing files stay inherited.
- `web/system/**` is immutable. Never create, copy, modify, override, or extend any system runtime, CSS, Tailwind/PostCSS, client, logic, template, generated file, or plugin.
- Pages compose section components; layouts own shared chrome. Split reusable or independently styled/interactive sections.

### GOSH and runtime

- Enable `runtime.enabled` only for web applications. Retain enabled consent/search/commands/preloader/shortcuts when replacing a layout.
- A renderable `.gosh` file has exactly one `<template>`. Use verified GOSH syntax, not Vue/Nuxt syntax.
- Omit `@layout` for the configured default. Use `@layout name` only for another layout and `@layout none` only to disable layouts.
- File routes use `index.gosh`, `[id].gosh`, and `[...all].gosh`. Server logic uses `package logic` and `github.com/myelophone/goserver/web/runtime`; regenerate bindings after handler/import changes.
- SSR initial content. Use client components, hydration, `ClientOnly`, and `useQuery` only for actual client behavior.
- Use documented `_gosh` navigation, hooks, actions/forms, query, stores, SEO, and loading. Never touch `/_gosh/*`, runtime tokens, or internal markup.
- Internal navigation uses real links. Do not use `data-runtime-off`, `window.location`, or reloads as fixes.
- Async UI has pending/empty/error/success states and cleans up listeners, observers, timers, and subscriptions.

### Native components, themes, stores, forms

- Inspect/reuse inherited `Grid*`, `Ui*`, `View*`, `Cookie*`, `Consent*`, and `Seo*` components.
- Use `GridContainer` and `--layout-container-*`. Header, navigation, sections, and footer share one grid.
- Use `--ui-*`, `data-theme`, and the `preferences` store. New UI works in every enabled theme. Do not duplicate theme state/storage.
- Use the inherited `ui` store for its modal/search/command/language state and `cookies` for consent. New `web/stores/*.js` are only for genuinely new public shared state; hydrate public SSR state with `runtime.UseStoreState`. Keep secrets/authority server-side.
- Do not use `data-gosh-store-*` bindings. Use an inherited component or a proper client component with the documented store API.
- Forms must be visually and behaviorally consistent through goserver's inherited form controls, validation/error patterns, and GOSH form/action contract. Do not build parallel controls, custom submission transports, or per-form loading/error conventions. Validate and authorize again on the server; binary uploads use a protected upload endpoint.
- Use goserver image/picture optimization, search, command palette, loading/preloader, consent, and SEO instead of substitutes.
- Entrance animation uses `data-gosh-reveal` and its documented speed/step/repeat options. Do not add reveal libraries, scroll listeners, or another `IntersectionObserver`.

### Tailwind v4 and CSS

- Tailwind CSS v4 utilities are the default styling method. Use the goserver-managed Tailwind/PostCSS toolchain; do not install/configure another Tailwind, PostCSS, package, lockfile, or CSS build.
- Application CSS belongs only in `web/css/default.css`, its local `web/css/**` imports, or component `<style scoped>`.
- Application files must not use `@apply`, `@theme`, `@utility`, `@variant`, Tailwind configuration directives, PostCSS plugins, CSS `!important`, or Tailwind important modifiers.
- Prefer inherited components, standard utilities, and existing tokens over arbitrary values or duplicated component CSS.
- Use inherited Tailwind v4 breakpoints (`sm`, `md`, `lg`, `xl`, `2xl`) and goserver container tokens; do not invent breakpoints/containers.

### Visual and responsive quality

- Define one coherent palette, type system, spacing rhythm, surfaces, and icon family. Use intentional fonts with fallbacks; self-host licensed files and load only used weights.
- Avoid generic AI styling: gratuitous gradients/glass/glow/blobs, pills everywhere, empty oversized heroes, and repetitive card grids.
- Long text uses the text container and a readable measure. Use semantic landmarks/headings, real links/buttons, labels, visible focus, keyboard support, contrast, reduced motion, and correct alt text.
- A one-screen section fits the available visual viewport in both axes, including content, header, padding, safe areas, and mobile browser chrome. Use applicable goserver viewport/snap components and `dvh`, never fixed `100vh` or guessed heights. Adapt or split content; never clip/mask overflow.
- Images have dimensions/aspect ratios, responsive sizing, deliberate crops, alt text, and verified production loading. Do not lazy-load critical above-fold media.
- Set title, description, canonical, favicon, social image, and indexing policy. Private pages stay out of search/sitemaps.

### i18n and consent

- Use goserver i18n, never custom maps or duplicated language page trees. Dictionaries live at `i18n/<language>/<namespace>.json` and are registered with `RegisterI18nFS` before `NewI18n`.
- With web enabled, configure locale routes through `locales`/`defaultLocale` in `websettings.json`; without web use `I18N_DEFAULT_LANGUAGE`/`I18N_LANGUAGES`.
- Web logic uses `ctx.T`; handlers use goserver `L`/`Lf`/plural helpers. Verify localized SSR, metadata, `html[lang]`, canonical/hreflang, direct loads, refresh, runtime navigation, and history.
- When `cookieControl.enabled` is true, every non-essential external script belongs in `cookieScripts` under the correct `analytics`, `marketing`, or `functional` category and loads only after consent. Do not insert it directly into global scripts, layouts, or components. Use consent components/wrappers for external embeds; classify as `necessary` only when strictly required for core operation.

## 5. Backend and security

- Use goserver routing, middleware, parsers, validators, responders, errors, lifecycle/jobs, HTTP client, cache, and adapters.
- Register APIs with goserver's `s.GET`/`POST`/`PUT`/`PATCH`/`DELETE`/`HEAD`/`OPTIONS` methods or `s.Group`; attach shared middleware with `group.Use`. Use params/wildcards from the native router. With `API_PREFIX`, register routes without repeating the external prefix. In combined applications normally keep the global prefix empty and put APIs in an `/api` group; never add a root/catch-all route that shadows the web fallback.
- Define methods, DTOs, validation, statuses, authorization, and safe errors. Never mutate through GET or return `200` for every result.
- Authenticate/authorize every protected operation, including object access. Never trust client identity, tenant, role, or price.
- Bound bodies, uploads, pagination, concurrency, streams, and outbound work. Propagate context/cancellation; use parameterized queries and transactions.
- Keep CSRF, security headers, origin checks, and rate/body limits. Cookie and bearer auth use separate policies; CSRF trusted origins are not CORS.
- Secrets come from environment. Never log/return secrets, tokens, cookies, SQL, stacks, or personal data.
- Cache only output safe for its complete locale/tenant/identity key. Public/static routes must be identical for every visitor.
- Multi-tenancy uses `TenantStore` middleware and `web/tenants/<id>/{pages,content}`; do not build another tenant resolver or duplicate site trees.

## 6. Validation

- Use repository tasks for the selected mode. API-only work does not generate web artifacts. Web work regenerates bindings through the standard web run/build tasks.
- Format and run relevant tests, vet, production build, and the built server. Exercise success, validation, error, authorization, and asset paths.
- Browser-check web work across every used goserver breakpoint and enabled theme: direct/refresh/history navigation, forms, stores, i18n, overflow, viewport-height sections, focus, reduced motion, console/network errors, and loaded images.
- Review the diff for user changes, secrets, generated files, framework copies, redundant implementations, dependencies, and placeholders.
- Report compilation, tests, production execution, and visual verification separately.
