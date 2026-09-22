# Project instructions

- Build the requested application with goserver. Inspect `git status`, relevant files, and the version in `go.mod`; verify APIs in that version's README/source before use.
- Keep changes focused and preserve user code. Use goserver's documented subsystems before adding dependencies or infrastructure. Never patch, copy, vendor, or modify goserver or its module cache; report missing extension points.
- Keep `cmd/main.go` for assembly. Add application code in new `app/*.go` files; do not edit `app/app.go`. Reserve `internal/` for generated code.
- Do not edit generated, compiled, cached, or tool-owned files, including `internal/goservergen/**`, `cmd/web_import_gen.go`, `tmp/goserver/build/**`, and `dist/**`. Regenerate with repository tasks.
- No placeholders, fake controls, TODO flows, invented facts, copied framework code, or inherited welcome page in a finished site. Comment only non-obvious code. Keep command output and responses concise.
- Prevent memory and resource leaks: bound long-lived collections and caches, release references to objects no longer needed, close owned resources, stop timers, and remove event handlers and subscriptions when their lifecycle ends. When a leak is suspected, use profiling to investigate retained memory and goroutine growth.

## Select the mode

- **REST/API:** leave web disabled. Build with goserver Go routes in `app/*.go`; do not add GOSH pages or client UI.
- **Website:** enable web. Use GOSH for pages and add Go routes when needed.
- **Combined:** enable web. Use GOSH pages plus explicit Go API routes on the same listener; normally group APIs under `/api` and avoid catch-all routes that shadow pages.
- Set `MYELOPHONE_WEB_ENABLED=true/false` or `runtime.enabled=true/false` in `websettings.json`. The environment variable takes precedence when present. Runtime and build tasks use the resolved value.

## Goserver capabilities

Verify the selected version before using a feature.

| Area | Native capabilities |
| --- | --- |
| HTTP | Method routes, params/wildcards, groups/middleware, standard handlers, request formats, validation, responses/errors, redirects/rewrites, templates/static files |
| Middleware | Defaults/custom stack, hooks, logging/request IDs, recovery, headers, CSRF, limits, rate limiting, timeouts, load shedding, maintenance |
| Data/state | Cookies, sessions, JWT/encryption, cache/idempotency, PostgreSQL helpers, Redis cache/sessions, tenants |
| Realtime | HTML/JSON streaming, SSE, NDJSON, WebSockets |
| Integrations | Resilient HTTP client, browser TLS, proxies/Webshare, HTML/soft-404 helpers, cron/jobs, SMTP, Telegram |
| Operations | Config/utilities, lifecycle/tracked jobs/shutdown, health, `/metricz`, pprof, logging, Docker, build/assets |
| GOSH, when enabled | SSR/file routes, layers/layouts/components, actions/forms, islands/hydration, navigation/hooks/query/loading/errors |
| Web platform, when enabled | Grid/UI/consent/SEO/i18n, Tailwind v4, themes/stores, images/content/search/commands/reveal/transitions, route rules/SWR/public-static, cache tags/optimistic actions, preload/prefetch, Early Hints/Server-Timing/Web Vitals |

## Server work: all modes

- Use native routes (`s.GET`, `s.POST`, etc.), `s.Group`, `group.Use`, parsers, validators, responders, errors, lifecycle, clients, cache, and adapters. With `API_PREFIX`, do not repeat the prefix in route paths.
- Use `ParseRequest` and `ValidationRule`/`Validate` for JSON, URL-encoded, and multipart input; keep field validation and response/error handling consistent across endpoints.
- Define correct methods, DTOs, validation, status codes, and safe errors. Authenticate and authorize protected actions and objects. Bound input, uploads, pagination, concurrency, streams, and outbound work; propagate cancellation and use parameterized queries/transactions.
- Keep CSRF, headers, origin checks, and rate/body limits. Separate cookie and bearer policies; CSRF trusted origins are not CORS. Read secrets from the environment; never log or return secrets or personal data.
- Cache by complete locale/tenant/identity key. Use `TenantStore` for multi-tenancy. Use goserver i18n: `RegisterI18nFS` before `NewI18n`; API locales use `I18N_DEFAULT_LANGUAGE`/`I18N_LANGUAGES`.

## Website work: only when requested

- Keep routes thin in `web/pages`; put shared chrome in `layouts`, reusable sections in `components`, Go data/actions in `logic`, and optional file endpoints in `server`. Use other `web/` directories only for their named features; put public files in `assets/`.
- Inherited files remain available; matching project files replace them. Development layers are `playground → project → goserver`; production layers are `project → goserver`. Do not create or override `web/system/**`.
- Web pages, layouts, and components are `.gosh` files with exactly one `<template>`. Use optional `<head>`, `<style scoped>` for local CSS, plain `<style>` for intentional global CSS, `<script setup>` for setup code, and `<script>` for client code. Follow documented `.client.gosh`/`.server.gosh` rules; server components cannot contain ordinary `<script>` blocks.
- Use documented GOSH file routes, layouts, and directives; do not assume Vue/Nuxt syntax. `web/logic` uses `package logic` and `github.com/myelophone/goserver/web/runtime`. Regenerate bindings after handler/import changes.
- SSR initial content. Use documented `_gosh` navigation, actions/forms, hooks, stores, SEO, and loading. Use client behavior only when needed; do not alter runtime internals. Internal navigation uses real links.
- Reuse native grid, UI, consent, SEO, search, image optimization, themes, and stores. Keep authority server-side. Async UI handles pending/empty/error/success and cleans up subscriptions/listeners.
- Keep forms consistent with inherited `UiInput`, `UiTextarea`, `UiCheckbox`, `UiButton`, labels, field errors, and pending/success states. GOSH action forms use `data-gosh-form` and server `Action`/`props.Form()`; use `_gosh.useForm` only for client validation or reactive form state. Do not add an HTML `action` or custom submission transport for these forms. Validate and authorize again on the server; upload binaries through a protected endpoint.
- Use goserver-managed Tailwind v4/PostCSS, existing utilities/tokens, `web/css/default.css`, and scoped component styles. Do not add a parallel CSS pipeline or use `@apply`, `@theme`, `@utility`, `@variant`, `!important`, or new breakpoints. Use `data-gosh-reveal` for entrance animation.
- Use semantic headings, accessible controls, focus states, contrast, reduced motion, responsive images, and correct alt text. Fit one-screen sections to the real viewport with `dvh`, not fixed `100vh` or clipping. Set SEO metadata; exclude private pages from indexing.
- Give each page one descriptive, visible `<h1>` near the start of `<main>`, consistent with its `<title>`; shared headers/layouts must not add another. Use `<h2>` for main sections and `<h3>` onward for subsections, without skipping levels or choosing levels for font size.
- Prefer inherited `UiHeading` and goserver section components. Set their semantic `level` for the page outline and `size` for appearance; verify the final SSR heading order after composing components.
- Web locales use `locales`/`defaultLocale` in `websettings.json`. Load non-essential external scripts through `cookieScripts` after consent when cookie control is enabled.

## Website design

- Before coding, define the audience, purpose, primary action, page structure, and visual direction. Infer sensible defaults from the brief.
- Structure pages around visitors' questions and tasks. Every section must add useful information or enable an action; do not default to a generic hero, feature-card grid, testimonials, and CTA.
- Write specific, informative copy. Avoid filler, repeated claims, and invented statistics, clients, reviews, or credentials.
- Establish a coherent visual identity with deliberate typography, spacing, alignment, contrast, imagery, surfaces, and icon style. Avoid interchangeable templates, gratuitous gradients, glass, glow, blobs, and repetitive cards.
- Use references for their useful design principles without copying branding or content. Keep the composition specific to the subject while using goserver's components and styling rules.
- Design mobile layouts deliberately: prioritize content, adapt navigation, and simplify compositions rather than merely stacking desktop blocks.
- Inspect rendered desktop and mobile pages. Refine weak hierarchy, awkward spacing, poor wrapping, unbalanced sections, and visual repetition before finishing.

## Verify

- Format changed code and run relevant checks only. Build/run production and inspect browser output when needed for the changed behavior. Avoid repeated passing checks and long logs.
- Review the diff for user changes, secrets, generated files, copied framework code, unnecessary dependencies, and placeholders. Report meaningful verification results briefly.
