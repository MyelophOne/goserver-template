# Project instructions

## Core rules

- Build with the repository's installed `github.com/myelophone/goserver` version. First inspect `git status`, relevant files, `go.mod`, and the version-matched goserver README/source; never guess APIs.
- Preserve existing code, architecture, content, and behavior unless the task requires changes. Prefer documented goserver capabilities over new packages or infrastructure; never patch, copy, vendor, or modify goserver or its module cache. Report missing extension points.
- Keep `cmd/main.go` for assembly, add application Go code in `app/*.go`, never edit `app/app.go`, and reserve `internal/` for generated code.
- Never hand-edit generated, compiled, cached, or tool-owned files: `internal/goservergen/**`, `cmd/web_import_gen.go`, `tmp/goserver/build/**`, `dist/**`, or `web/system/**`. Use repository tasks to regenerate outputs.
- Deliver complete working features: no placeholder copy, fake buttons, TODO flows, fabricated facts, framework-code duplication, or inherited welcome page in a finished product.
- Keep edits focused. Comment non-obvious logic only. Format changed code, review the diff, and report the checks actually performed without excessive logs.
- Prevent resource leaks: bound long-lived caches and collections, release unneeded references, close owned resources, and clean up timers, goroutines, listeners, and subscriptions. Investigate suspected leaks with profiling.

## Choose an application mode

- **REST/API:** keep web disabled; use goserver Go routes in `app/*.go`. Do not introduce GOSH or a client UI.
- **Website:** enable web; use GOSH pages and add Go routes when necessary.
- **Combined:** enable web; use GOSH pages and explicit Go API routes on the same listener, normally under `/api`. Avoid catch-all routes that shadow pages.
- Configure `MYELOPHONE_WEB_ENABLED=true/false` or `runtime.enabled=true/false` in `websettings.json`; the environment variable takes precedence. Build and runtime tasks use the resolved setting.

## Frontend design skill — mandatory when applicable

- For **new websites, landing pages, UI creation, substantial visual changes, and requested redesigns**, read and follow `.agents/skills/frontend-design/SKILL.md` before implementation. Treat its design process and real-browser review as acceptance criteria.
- For **small UI bug fixes**, apply only the skill's relevant responsive/accessibility/verification rules; do not redesign unrelated surfaces.
- For **REST/API-only or backend-only work**, do not load or apply the design skill.
- User-supplied references and explicit design requirements override default aesthetic choices; do not overwrite an existing visual identity unless a redesign is requested.
- A successful compilation is not proof of good UX, responsive layout, or visual quality. Never claim browser verification unless rendered pages were inspected.

## Goserver: backend and API

- Verify the installed goserver version and source before using any native feature. Prefer its routing, middleware, parsing, validation, responses, errors, authentication, sessions, cache, database adapters, integrations, observability, and lifecycle APIs over custom frameworks.
- Use method routes (`s.GET`, `s.POST`, etc.), `s.Group`, and `group.Use`; honor middleware ordering. With `API_PREFIX`, do not repeat the prefix in route paths. Avoid catch-all routes that shadow more specific endpoints.
- Use documented `ParseRequest` and `ValidationRule`/`Validate` for JSON, URL-encoded, and multipart requests. Use explicit request/response DTOs, correct HTTP methods and statuses, consistent error formats, and safe public error messages.
- Authenticate protected requests and authorize every protected action and object, including tenant-scoped access. Revalidate server-side even when client validation exists. Never expose secrets or personal information in URLs, errors, logs, or responses.
- Preserve documented middleware and security defaults: request IDs, logging, recovery, security headers, CSRF, origin checks, body limits, rate limiting, timeouts, load shedding, and maintenance behavior. Do not bypass the configured stack.
- Keep cookie/session security separate from bearer-token rules; CSRF trusted origins are not CORS. Enforce appropriate cookie flags, token handling, and access checks as supported by installed APIs; read secrets from environment/configuration, not source.
- Bound input sizes, file uploads, pagination, concurrency, queues, streaming, retries, outbound requests, and work duration. Propagate request cancellation and deadlines. Parameterize queries, use transactions for multi-step writes, and close owned DB rows, streams, bodies, and connections.
- Reuse goserver cache, idempotency, sessions, JWT/encryption, PostgreSQL helpers, and Redis adapters where applicable. Key cached responses by all identity, authorization, tenant, locale, and request dimensions that affect the result; set explicit TTL/eviction bounds. Use `TenantStore` for tenant-aware storage.
- Follow documented APIs for SSE, NDJSON, HTML/JSON streaming, and WebSockets. Handle disconnects, backpressure, write failures, and cleanup; do not leave goroutines or subscriptions behind after clients exit.
- Prefer native resilient HTTP clients, proxies/Webshare adapters, HTML/soft-404 helpers, jobs/cron, SMTP, and Telegram integrations when required. Set bounded retries, timeouts, and concurrency; avoid duplicating goserver subsystems.
- Use the documented lifecycle for startup, tracked background jobs, shutdown, and graceful resource cleanup. Keep work cancellation-aware; stop timers and goroutines; close resources and cap long-lived caches/collections. When leaks are suspected, inspect pprof, memory retention, and goroutine counts rather than guessing.
- Keep observability intact: structured logs/request IDs, health endpoints, `/metricz`, profiling/pprof when configured, and relevant server timings. Do not log credentials, tokens, or sensitive payloads.
- Register translations with `RegisterI18nFS` before `NewI18n`; API languages use `I18N_DEFAULT_LANGUAGE`/`I18N_LANGUAGES`. Preserve locale-aware response and cache behavior.
- For API changes, verify routes, methods, validation errors, unauthorized/forbidden handling, not-found cases, cancellation, limits, and representative success paths with relevant tests. Report checks actually run.

## GOSH website architecture

- Keep `web/pages` routes thin; shared chrome belongs in `web/layouts`, reusable sections in `web/components`, Go data/actions in `web/logic`, optional file endpoints in `web/server`, and public files in `assets/`.
- Inherited files can be overridden by matching project files. Layer order: development `playground → project → goserver`; production `project → goserver`. Never create or override `web/system/**`.
- Pages/layouts/components are `.gosh` files with exactly one `<template>`. Supported blocks include optional `<head>`, `<style scoped>`, intentional global `<style>`, `<script setup>`, and client `<script>`; follow version-specific `.client.gosh` and `.server.gosh` restrictions. Never assume Vue/Nuxt syntax.
- Use documented file routing, directives, layouts, SSR, `_gosh` navigation/actions/hooks/stores/loading, and SEO. `web/logic` uses `package logic` and `github.com/myelophone/goserver/web/runtime`. Regenerate bindings after handler/import changes.
- Render initial content with SSR. Internal navigation uses real links. Keep authority server-side; add client-side behavior only when needed. Async UI must provide pending/empty/error/success states and clean up listeners/subscriptions.
- Prefer native grid/UI, consent, SEO, search, images, themes, and stores. Reuse accessible `UiInput`, `UiTextarea`, `UiCheckbox`, `UiButton` where appropriate, without forcing layouts into stock arrangements.
- GOSH action forms use `data-gosh-form` with server `Action`/`props.Form()`; use `_gosh.useForm` only for client validation/reactive form state. Do not add an HTML `action` or custom transport. Revalidate/authorize server-side; binary uploads use protected endpoints.
- For new marketing pages, use documented `data-gosh-reveal` with restrained opacity/small-transform motion unless prohibited; support reduced motion, visible SSR content without JS, and client navigation. Never leak observers or lifecycle resources. Preserve motion conventions on existing sites unless requested otherwise.
- Maintain semantic headings (one visible descriptive `<h1>` per page, `<h2>` sections, `<h3>` subsections), accessible controls, focus, contrast, alternative text, responsive images, and SEO. Use `UiHeading` only when it fits; semantic levels and visual sizes are independent. Avoid clipping full-viewport sections; use `dvh` when relevant.
- Web locales use `locales`/`defaultLocale` in `websettings.json`; defer non-essential external scripts via `cookieScripts` when consent control is enabled. Private pages must not be indexed.

## Tailwind v4 integration

- Use goserver's managed Tailwind v4/PostCSS pipeline. No CDN, second import, v3-style configuration, or changes to framework-owned config. Project styles belong in `web/css/default.css` and GOSH styles.
- Respect this project's restrictions on `@apply`, `@theme`, `@utility`, `@variant`, and `!important`. These are project integration constraints, not generic Tailwind limitations.
- Use utilities for layout/states, CSS variables for shared design tokens, and scoped CSS for specialized components. A CSS custom property does not automatically create a named utility; reference it explicitly when necessary (e.g. `bg-[var(--brand)]`).
- Classes must be statically discoverable; choose complete strings rather than interpolating class fragments (never `bg-${color}-600`). Verify classes originating in Go/client code appear in the generated CSS; use validated CSS variables for truly dynamic values.
- Base utilities apply on mobile; `sm:`/`md:` apply upward. Choose breakpoints for layout needs; use container queries for reusable components when helpful.
- Avoid conflicting utilities; class-attribute order does not resolve cascade conflicts. Don't fix conflicts with `!important`. Pair hover with `focus-visible` and actual disabled/selected/expanded/invalid states; avoid indiscriminate `transition-all`.

## Verification and completion

- Run relevant formatting, tests, and builds for modified areas; do not repeat checks that have already passed without changes.
- For new websites or production-affecting visual changes, build via `task build` with web enabled, serve via `task server`, and inspect the built site. `task preview` is not the built distribution.
- For website work, use the design skill's browser, responsive, accessibility, interaction, animation, and fragment-navigation checks. Verify both direct URL loads and client navigation.
- Review the final diff for secrets, generated files, accidental user-code changes, unnecessary dependencies, fake controls, and placeholders.
- State precisely what was tested and what remains unverified. If browser tools are unavailable, report visual/interaction acceptance as pending rather than calling the site finished.
