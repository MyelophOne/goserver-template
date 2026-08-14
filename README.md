# @myelophone/goserver-template

A production-oriented Go backend boilerplate powered by [`@myelophone/goserver`](https://github.com/myelophone/goserver). It provides a ready application entry point, an extensible `app` package, routing and middleware defaults, HTML templates, static assets, environment configuration, live reload, Docker deployment, health checks, profiling, testing, linting, security scans, load testing, and release tasks.

**[Create a new repository from this template](https://github.com/myelophone/goserver-template/generate)** · **[Open @myelophone/goserver documentation](https://github.com/myelophone/goserver)**

> **Building a full-stack application?**
> `@myelophone/goserver-template` works seamlessly with [`@myelophone/nuxt`](https://github.com/myelophone/nuxt), the companion frontend foundation. Start the frontend from [`@myelophone/nuxt-template`](https://github.com/myelophone/nuxt-template), which provides a ready boilerplate for quickly creating the client side of an `@myelophone/goserver` application.
>
> **[Open @myelophone/nuxt](https://github.com/myelophone/nuxt)** · **[Open @myelophone/nuxt-template](https://github.com/myelophone/nuxt-template)** · **[Create a frontend repository](https://github.com/myelophone/nuxt-template/generate)**

Use this repository when you want to start building application routes and business logic immediately instead of assembling the server lifecycle, development tooling, directory layout, and deployment files yourself.

## Contents

- [What is already prepared](#what-is-already-prepared)
- [Frontend integration with @myelophone/nuxt](#frontend-integration-with-myelophonenuxt)
- [Quick start](#quick-start)
- [Project structure](#project-structure)
- [How the template boots](#how-the-template-boots)
- [Where application code belongs](#where-application-code-belongs)
- [Complete quick-start application](#complete-quick-start-application)
- [Routing and middleware](#routing-and-middleware)
- [Unified form and request parsing](#unified-form-and-request-parsing)
- [@myelophone/goserver capabilities available to the template](#myelophonegoserver-capabilities-available-to-the-template)
- [HTML templates](#html-templates)
- [Static assets](#static-assets)
- [Configuration](#configuration)
- [Development commands](#development-commands)
- [Docker and deployment](#docker-and-deployment)
- [Updating @myelophone/goserver](#updating-myelophonegoserver)
- [License](#license)

## What is already prepared

The template includes:

- a minimal executable in `cmd/main.go`;
- the opinionated `goserver.Defaults()` production middleware stack;
- an `app` package loaded for side effects and connected through server hooks;
- automatic HTML page routing from `templates/pages`;
- layouts, partials, global assets, and page-specific assets;
- static files and `robots.txt` from `assets`;
- `.env` loading through Task;
- development, preview, optimized build, test, lint, vulnerability, analysis, and release tasks;
- live reload with Air on Windows and Unix-like systems;
- a multi-stage, non-root distroless Docker image;
- Docker Compose health checks, restart policy, and bounded container logs;
- built-in `/healthz`, `/metricz`, and protected pprof support from `@myelophone/goserver`;
- a Dev Container definition for an isolated Go development environment.

The template currently pins `github.com/myelophone/goserver v0.5.1` in `go.mod`. Update that dependency deliberately and verify the application before adopting a newer release.

## Frontend integration with @myelophone/nuxt

[`@myelophone/goserver-template`](https://github.com/myelophone/goserver-template) and [`@myelophone/nuxt`](https://github.com/myelophone/nuxt) form a ready full-stack foundation:

- `@myelophone/goserver` owns API routes, validation, authentication, sessions, caching, databases, background work, WebSockets, and operational endpoints;
- `@myelophone/nuxt` owns the browser application and consumes the backend through HTTP and WebSocket APIs;
- [`myelophone/nuxt-template`](https://github.com/myelophone/nuxt-template) provides the prepared frontend project, just as this repository provides the prepared backend project.

For a clean same-origin deployment, expose `@myelophone/goserver` below `/api`:

```dotenv
API_PREFIX=/api
```

Application code still registers prefix-free routes:

```go
s.GET("/users", func(w http.ResponseWriter, r *http.Request) {
	s.RespondJSON(w, r, []map[string]any{
		{"id": 1, "name": "Ada"},
	})
})
```

The frontend then calls `/api/users`. Keeping the public prefix in configuration lets the same backend run behind a Nuxt server, reverse proxy, ingress, or separate development port without adding `/api` to every Go route.

When the frontend and backend use different origins, configure the trusted origins and cookie/proxy policy for the actual deployment instead of weakening origin checks globally. The same-origin `/api` arrangement is the simplest choice for cookie-backed sessions and CSRF protection.

## Quick start

### 1. Create or clone the project

The fastest option is **[Use this template](https://github.com/myelophone/goserver-template/generate)** on GitHub. To clone it directly:

```bash
git clone https://github.com/myelophone/goserver-template.git my-backend
cd my-backend
```

### 2. Prepare the environment

macOS/Linux:

```bash
cp .env.example .env
```

PowerShell:

```powershell
Copy-Item .env.example .env
```

The example secrets are suitable only for local startup. Generate independent, high-entropy `SESSION_KEY`, `JWT_SECRET`, `WS_TOKEN_KEY`, and `METRICS_SECRET` values before deploying the service.

### 3. Run immediately

Without additional tools (`@myelophone/goserver` defaults are used; `go run` does not load `.env` automatically):

```bash
go mod download
go run ./cmd/main.go
```

Or with Task:

```bash
task run
```

Open <http://localhost:8080/healthz> or test it from a terminal:

```bash
curl -A 'Mozilla/5.0' http://localhost:8080/healthz
```

`server.Defaults()` rejects empty user agents and several scanner/automation signatures, including curl's default user agent. The explicit browser-like user agent above is intentional.

### 4. Rename the Go module

When the repository becomes a real project, change its module path:

```bash
go mod edit -module github.com/acme/orders-api
go mod tidy
```

Then update the blank import in `cmd/main.go`:

```go
_ "github.com/acme/orders-api/app"
```

That import is what executes the `init` functions in the application package and registers its hooks.

## Project structure

```text
goserver-template/
├── .devcontainer/          # containerized Go development environment
├── app/
│   └── app.go              # base application hook registration
├── assets/
│   └── robots.txt          # public static file served by @myelophone/goserver
├── cmd/
│   ├── healthcheck/        # static health probe used by the container
│   ├── main.go             # application composition root
│   ├── dev.sh              # Air build helper for Unix-like systems
│   └── dev.bat             # Air build helper for Windows
├── templates/
│   ├── head/               # global and per-page head fragments
│   ├── layouts/            # optional custom layouts
│   ├── pages/              # automatically routed HTML pages
│   ├── partials/           # reusable template definitions
│   ├── scripts/            # global and per-page scripts
│   ├── styles/             # global and per-page styles
│   └── _examples.md        # source template examples
├── .air.toml               # live-reload configuration
├── .env.example            # documented environment baseline
├── docker-compose.yml      # local/container orchestration
├── Dockerfile              # optimized production image
├── go.mod / go.sum         # Go module and locked dependencies
└── taskfile.yml            # development and maintenance commands
```

## How the template boots

`cmd/main.go` performs a small, explicit sequence:

```text
import app package
        ↓
create goserver.Server
        ↓
install Defaults middleware
        ↓
load templates and automatic page routes
        ↓
apply application hooks
        ↓
add operational routes and run
```

The actual composition root is intentionally short:

```go
server := goserver.NewServer(":" + goserver.GetEnv("HTTP_PORT", "8080"))
server.Defaults()

tm := goserver.NewTemplateManager()
server.TemplatesMiddleware(tm)

server.ApplyHooks()
server.Run()
```

This order has useful consequences:

- framework defaults apply to application routes;
- filesystem pages are registered before application hooks;
- a hook can replace an automatically generated static page with a dynamic handler registered at the same path;
- all hook-based routes exist before `Run` compiles the final middleware chain;
- `/healthz`, `/metricz`, server context, optional gzip, and runtime rate limiting are installed by `Run`.

## Where application code belongs

Keep `cmd/main.go` as the composition root. Add application code as new `.go` files inside `app` or in domain packages imported by `app`.

The bundled `app/app.go` is template-managed and explicitly marked as a file that should not be modified. Add files such as:

```text
app/
├── app.go
├── routes.go
├── middleware.go
├── users.go
└── billing.go
```

Every file uses `package app`. Register server-level functionality from `init`:

```go
package app

import (
	"net/http"

	"github.com/myelophone/goserver"
)

func init() {
	goserver.RegisterHook(func(s *goserver.Server) {
		s.GET("/ready", func(w http.ResponseWriter, r *http.Request) {
			s.RespondJSON(w, r, map[string]string{"status": "ready"})
		})
	})
}
```

Multiple files may register hooks. Use `RegisterHookWithPriority` when installation order matters; smaller priority values run first. `cmd/main.go` calls `server.ApplyHooks()` once after templates have loaded.

## Complete quick-start application

Create `app/routes.go` with the following code. It demonstrates routes, path parameters, unified JSON/form parsing, validation and type normalization, sessions, cache-aside generation, request IDs, and JSON responses without modifying the template bootstrap.

```go
package app

import (
	"context"
	"net/http"
	"time"

	"github.com/myelophone/goserver"
)

type greeting struct {
	Message   string `json:"message"`
	Generated string `json:"generated"`
}

func init() {
	goserver.RegisterHook(func(s *goserver.Server) {
		s.Cache = goserver.NewCache(1_000, "")
		s.SetSessionStorage(goserver.NewInMemorySessionStorage(1_000))

		s.GET("/", func(w http.ResponseWriter, r *http.Request) {
			s.RespondJSON(w, r, map[string]any{
				"message":    "@myelophone/goserver-template is running",
				"request_id": goserver.GetRequestID(r.Context()),
				"try": []string{
					"GET /hello/Ada",
					"POST /users",
					"POST /session",
					"GET /session",
					"GET /cached",
					"GET /healthz",
				},
			})
		})

		s.GET("/hello/:name", func(w http.ResponseWriter, r *http.Request) {
			name := s.GetParams(r).Get("name")
			s.RespondJSON(w, r, map[string]string{"message": "Hello, " + name + "!"})
		})

		s.POST("/users", func(w http.ResponseWriter, r *http.Request) {
			input, err := goserver.ParseRequest(r)
			if err != nil {
				s.RenderErrorJSON(w, r, http.StatusBadRequest, err.Error())
				return
			}

			validationErrors := input.Validate([]goserver.ValidationRule{
				{Name: "email", Type: "string", Required: true, Regex: goserver.RegexEmail},
				{Name: "age", Type: "int", Required: true, Min: goserver.FloatPtr(18)},
				{Name: "role", Type: "string", Default: "reader", Enum: []any{"reader", "editor"}},
			})
			if len(validationErrors) != 0 {
				s.RenderErrorJSON(w, r, http.StatusUnprocessableEntity, validationErrors[0].Error())
				return
			}

			age, ok := input.GetInt("age")
			if !ok {
				s.RenderErrorJSON(w, r, http.StatusUnprocessableEntity, "age must be an integer")
				return
			}

			s.RespondJSON(w, r, map[string]any{
				"email": input.GetString("email"),
				"age":   age,
				"role":  input.GetString("role"),
			})
		})

		s.POST("/session", func(w http.ResponseWriter, r *http.Request) {
			id := s.SessionStart(&w, r)
			id.SetSessionValue("started_at", time.Now().Format(time.RFC3339))
			s.RespondJSON(w, r, map[string]bool{"session_started": true})
		})

		s.GET("/session", func(w http.ResponseWriter, r *http.Request) {
			id := s.SessionStart(&w, r)
			s.RespondJSON(w, r, id.GetAllSessionValues())
		})

		s.GET("/cached", func(w http.ResponseWriter, r *http.Request) {
			value, err := goserver.Fetch(r.Context(), s.Cache, "demo:greeting", 30*time.Second,
				func(ctx context.Context) (greeting, error) {
					return greeting{
						Message:   "generated once, then served from cache",
						Generated: time.Now().Format(time.RFC3339Nano),
					}, nil
				},
			)
			if err != nil {
				s.RenderErrorJSON(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			s.RespondJSON(w, r, value)
		})
	})
}
```

Run it:

```bash
task run
```

Exercise the endpoints:

```bash
curl -A 'Mozilla/5.0' http://localhost:8080/
curl -A 'Mozilla/5.0' http://localhost:8080/hello/Ada

# JSON input
curl -A 'Mozilla/5.0' -H 'Content-Type: application/json' \
  -d '{"email":"ada@example.com","age":37,"role":"editor"}' \
  http://localhost:8080/users

# The same handler with a browser-style form
curl -A 'Mozilla/5.0' -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'email=ada@example.com' --data-urlencode 'age=37' \
  http://localhost:8080/users

# Preserve the encrypted session cookie between requests
curl -A 'Mozilla/5.0' -c cookies.txt -X POST http://localhost:8080/session
curl -A 'Mozilla/5.0' -b cookies.txt http://localhost:8080/session

# The Generated value remains stable for 30 seconds
curl -A 'Mozilla/5.0' http://localhost:8080/cached
curl -A 'Mozilla/5.0' http://localhost:8080/cached
```

Browser form controls arrive as strings. `GetInt` can parse a numeric string, and `ValidationRule{Type: "int"}` also replaces the stored field with a normalized Go `int` before business logic reads it.

## Routing and middleware

Routes registered from `app` support `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, `OPTIONS`, and `ANY`. Every registration method accepts the path followed by a standard `http.HandlerFunc`.

Named parameters and trailing wildcards are read through `s.GetParams(r)`:

```go
s.GET("/files/*path", func(w http.ResponseWriter, r *http.Request) {
	s.RespondText(w, s.GetParams(r).Get("path"))
})
```

A custom middleware uses the standard Go signature:

```go
func requestHeader(name, value string) goserver.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set(name, value)
			next.ServeHTTP(w, r)
		})
	}
}
```

Group related endpoints and install middleware only for that group:

```go
api := s.Group("/api/v1")
api.Use(requestHeader("X-API-Scope", "v1"))
api.GET("/profile", func(w http.ResponseWriter, r *http.Request) {
	s.RespondJSON(w, r, map[string]string{"scope": r.Header.Get("X-API-Scope")})
})
```

`server.Defaults()` already installs request IDs, environment-aware access logging, slow-request reporting, load shedding, panic recovery, rate limiting, favicon/static files, canonical redirects, browser and malicious-request filtering, blocked-path protection, URL sanitization, CSRF checks, security headers, timeouts, body limits, bot/AI classification, and idempotency.

Add custom middleware from an application hook with `s.Use(...)`. The first registered middleware is the outermost wrapper.

## Unified form and request parsing

`goserver.ParseRequest` gives JSON APIs and browser forms the same handler-facing model:

| Request format                      | Result                                                                       |
| ----------------------------------- | ---------------------------------------------------------------------------- |
| `application/json`                  | Object properties in `ParsedRequest.Fields`; nested data remains structured. |
| `application/x-www-form-urlencoded` | Single controls become strings; repeated controls become `[]string`.         |
| `multipart/form-data`               | Text controls use `Fields`; uploads are grouped in `ParsedRequest.Files`.    |
| Other content types                 | Original bytes are available in `ParsedRequest.Raw`.                         |

Typed getters—`GetString`, `GetInt`, `GetFloat`, `GetBool`, `GetSlice`, `GetMap`, and `GetFiles`—let business logic work the same way regardless of how the client submitted the form.

`ValidationRule` supports required/default values, strings, integers, floats, booleans, nested objects, slices, enums, regular expressions, minimum/maximum values, length limits, file size limits, and MIME allowlists. Built-in regular expressions cover email, phone, URL, and UUID values.

For uploads:

```go
validationErrors := input.Validate([]goserver.ValidationRule{
	{Name: "avatar", Type: "file", Required: true,
		AllowedMime: []string{"image/jpeg", "image/png"},
		MaxFileBytes: 5 << 20},
})
if len(validationErrors) != 0 {
	s.RenderErrorJSON(w, r, http.StatusUnprocessableEntity, validationErrors[0].Error())
	return
}

avatars := input.GetFiles("avatar")
```

## @myelophone/goserver capabilities available to the template

The template exposes the complete public `@myelophone/goserver` API. The most important subsystems are summarized here; see the [@myelophone/goserver repository](https://github.com/myelophone/goserver) for the full API documentation and scenario examples.

### Responses and errors

- JSON, raw JSON, secure-prefixed JSON, ASCII JSON, XML, HTML, text, and raw responses;
- custom status codes and headers;
- HTML or JSON errors selected by `APP_ERROR_MODE`;
- custom not-found and internal-error handlers;
- file serving and forced downloads;
- redirects, cache headers, and preload headers;
- error notifications with request context and stack-enriched errors.

### Cookies, sessions, JWT, and encryption

- cookie read/write/remove helpers;
- AES-GCM encrypted `session_id` cookies;
- bounded in-memory session storage for one process;
- Redis session storage for multiple instances;
- typed `GetOrSetSessionValue` generation;
- HS256 JWT creation and validation with expiry;
- signed short-lived WebSocket tokens;
- standalone AES-GCM string encryption/decryption.

### Caching and idempotency

- bounded in-memory LRU cache;
- optional disk-backed cache;
- Redis cache adapter;
- generic `Fetch[T]` cache-aside loading;
- `FetchSWR[T]` stale-while-revalidate loading;
- singleflight miss coalescing;
- idempotent non-GET endpoints through `Idempotency-Key` and `s.Cache`.

### PostgreSQL and Redis

Import `github.com/myelophone/goserver/db` for pgx v5 pools, typed struct queries, execution helpers, batches, transaction retries, health checks, query deadlines, logging, and pool monitoring.

Import `github.com/myelophone/goserver/redis` (package name `goredis`) for distributed cache and session implementations.

### Outbound HTTP and proxies

`goserver.NewHttpClient` provides timeouts, retries with backoff, per-host circuit breakers, browser-like session headers, a cookie jar, DNS caching, metrics callbacks, proxy selection, and SSRF protection for internal/private addresses.

`ProxyManager` rotates providers and temporarily suppresses failed proxies. `WebshareProvider` directly integrates with the Webshare.io proxy-list API; implement `ProxyProvider` for another vendor.

### WebSockets

`WebSocketHub` supplies named channels, bounded queues, ping/pong handling, broadcasts, and an optional `WSBroker` interface for cross-instance fan-out. The first client message selects the channel; later messages are broadcast to subscribers.

### Background work and integrations

- interval and daily cron jobs;
- graceful tracking of application background tasks with `RunAsync`;
- SMTP mail with text/HTML bodies, attachments, inline files, and a bounded async queue;
- Telegram Bot API messages, photos, galleries, and documents;
- multi-tenant resolution through headers, query values, exact domains, or wildcard domains;
- HTML parsing/query/manipulation and soft-404 title detection;
- bot and AI client classification;
- hooks for setup, before/after requests, errors, and shutdown.

### Observability and protection

- request IDs and production/development access logs;
- connection and request statistics;
- `/healthz` status and protected detailed runtime metrics;
- Prometheus-compatible `/metricz` output;
- token-protected `/debug/pprof/*` routes;
- graceful SIGINT/SIGTERM shutdown and SIGHUP reload;
- maximum connections, concurrent request shedding, rate limits, deadlines, body/header/URL limits;
- CSRF and cross-origin protection, security headers, blocked malicious paths, URL sanitization, and panic recovery.

## HTML templates

The template directory is already prepared for layouts, pages, partials, and assets:

```text
templates/
├── layouts/      # base.html, alt.html, and other layouts
├── pages/        # index.html, contact.html, about@alt.html
├── partials/     # footer.html, header.html
├── head/         # global.html and page-specific head fragments
├── styles/       # global.html and page-specific style fragments
└── scripts/      # global.html and page-specific script fragments
```

`server.TemplatesMiddleware(tm)` automatically maps `pages/index.html` to `/`, `pages/contact.html` to `/contact`, and other page names to their extensionless path. A page such as `about@alt.html` renders through `layouts/alt.html` and remains available at `/about`.

### Page

Every page rendered with the bundled base layout defines `content`:

```html
{{define "content"}}
<h1>Contact us</h1>
<p>Phone: {{var . "phone"}}</p>

<ul>
 {{range (var . "emails")}}
 <li>{{.}}</li>
 {{end}}
</ul>

<ul>
 {{range (var . "contacts")}}
 <li>{{index . "name"}} — {{index . "email"}}</li>
 {{end}}
</ul>
{{end}}
```

`var` reads `PageData.Variables`. `default` supplies a value for nil, empty, or whitespace-only data:

```html
<p>{{default "Phone is not specified" (var . "phone")}}</p>
```

Variables can also be accessed directly:

```html
<p>{{index .Variables "phone"}}</p>
```

### Alternate layout

For `templates/layouts/alt.html`, the defined name must match the filename without `.html`:

```html
{{define "alt"}}
<!doctype html>
<html lang="en" class="nojs">
 <head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>{{.Title}}</title>
  <meta name="description" content="{{.Description}}" />
  {{.GlobalHead}} {{.PageHead}} {{.GlobalStyles}} {{.PageStyles}}
 </head>
 <body>
  <main>{{template "content" .}}</main>
  {{template "footer.html" .}} {{.GlobalScripts}} {{.PageScripts}}
 </body>
</html>
{{end}}
```

If `templates/layouts/base.html` is absent, `@myelophone/goserver` uses its embedded base layout. Use a page filename with `@` when the page needs another layout.

### Partial

`templates/partials/footer.html`:

```html
{{define "footer.html"}}
<footer><p>© 2026 MyelophOne</p></footer>
{{end}}
```

Include it with `{{template "footer.html" .}}`.

### Global and page-specific assets

| File                   | Field            | Scope                              |
| ---------------------- | ---------------- | ---------------------------------- |
| `head/global.html`     | `.GlobalHead`    | Every automatically rendered page. |
| `styles/global.html`   | `.GlobalStyles`  | Every automatically rendered page. |
| `scripts/global.html`  | `.GlobalScripts` | Every automatically rendered page. |
| `head/contact.html`    | `.PageHead`      | Only `contact.html`.               |
| `styles/contact.html`  | `.PageStyles`    | Only `contact.html`.               |
| `scripts/contact.html` | `.PageScripts`   | Only `contact.html`.               |

The asset files must contain their own `<style>` or `<script>` tags. They are inserted as trusted `template.HTML`, so only application-controlled markup belongs in these files.

The embedded base script changes the `nojs` class on `<html>` to `js`, which allows progressive-enhancement styling for both states.

### Dynamic rendering

An application handler can render a page with metadata and variables:

```go
s.GET("/contact", func(w http.ResponseWriter, r *http.Request) {
	data := goserver.PageData{
		Title:         "Contact",
		Description:   "Ways to contact the team",
		GlobalHead:    tm.GlobalHead,
		GlobalStyles:  tm.GlobalStyles,
		GlobalScripts: tm.GlobalScripts,
		Variables: map[string]any{
			"phone":  "+1-234-567-890",
			"emails": []string{"info@example.com", "support@example.com"},
			"contacts": []map[string]string{
				{"name": "Alice", "email": "alice@example.com"},
				{"name": "Bob", "email": "bob@example.com"},
			},
		},
	}

	htmlContent, err := tm.RenderHTML("contact.html", data)
	if err != nil {
		s.RenderError(w, r, http.StatusInternalServerError, "template rendering failed")
		return
	}

	s.RespondHTML(w, htmlContent)
})
```

Register this handler after `TemplatesMiddleware` to replace the automatic `/contact` endpoint. Additional source examples are in [`templates/_examples.md`](./templates/_examples.md).

## Static assets

Files under `assets` are served from `/assets/*` by `StaticAssetsMiddleware`, which is included in `Defaults`.

```text
assets/app.css       → /assets/app.css
assets/app.js        → /assets/app.js
assets/logo.svg      → /assets/logo.svg
```

The middleware assigns longer cache lifetimes to fonts and images, a one-week lifetime to CSS/JavaScript/font files, and shorter values to HTML/text/XML.

`cmd/main.go` also exposes `/robots.txt` through the same public-files handler. Edit `assets/robots.txt` for the deployed site's crawler policy.

## Configuration

Task loads `.env` automatically. `goserver.NewServer` reads most configuration during startup, so set variables before creating the server.

Durations accept values such as `500ms`, `15s`, and `2m`; a plain integer means seconds. Byte limits accept integers, `KB`, `MB`, `GB`, and expressions such as `1<<20`.

### Variables included in `.env.example`

| Variable                 | Example/default         | Purpose                                                                                                 |
| ------------------------ | ----------------------- | ------------------------------------------------------------------------------------------------------- |
| `DOCKER_PORT_BINDING`    | `127.0.0.1:0:8080`      | Compose host IP, host port, and container port. Host port `0` requests an automatically allocated port. |
| `HTTP_PORT`              | `8080`                  | Port used by `cmd/main.go`.                                                                             |
| `APP_ENV`                | `prod`                  | Selects production/development behavior. `task run` and `task dev` override it to `dev`.                |
| `API_PREFIX`             | empty                   | Optional common external route prefix such as `/api`.                                                   |
| `APP_ERROR_MODE`         | `html`                  | Default errors as `html` or `json`.                                                                     |
| `ENABLE_GZIP`            | `false` in the template | Enables response compression in `Run`.                                                                  |
| `CSRF_TRUSTED_ORIGINS`   | empty                   | Comma-separated exact or wildcard trusted origins.                                                      |
| `SMTP_HOST`, `SMTP_PORT` | empty / `25` fallback   | SMTP server and port.                                                                                   |
| `SMTP_USER`, `SMTP_PASS` | empty                   | SMTP credentials.                                                                                       |
| `SMTP_FROM`              | empty                   | Default sender address.                                                                                 |
| `SESSION_KEY`            | local example value     | AES-GCM session-cookie secret. Must be stable and private in production.                                |
| `TZ`                     | `Europe/Warsaw`         | Application timezone value.                                                                             |
| `JWT_SECRET`             | local example value     | HS256 JWT signing secret.                                                                               |
| `METRICS_SECRET`         | local example value     | Token for detailed health, metrics, and pprof access.                                                   |
| `METRICS_ENABLED`        | `true`                  | Enables `/metricz` and protected pprof routes.                                                          |

### Additional @myelophone/goserver variables

| Variable                    | Default            | Purpose                                                             |
| --------------------------- | ------------------ | ------------------------------------------------------------------- |
| `MAX_URL_LENGTH`            | `2048`             | Maximum request URL length.                                         |
| `MAX_HEADERS`               | `100`              | Maximum request header keys.                                        |
| `MAX_CONNECTIONS`           | `10000`            | Maximum active connections.                                         |
| `CONCURRENCY_LIMIT`         | `100`              | Maximum requests admitted concurrently by load shedding.            |
| `MAX_BODY_SIZE`             | `1MB`              | Maximum request body size.                                          |
| `MAX_HEADER_BYTES`          | `65536`            | Runtime HTTP header byte limit.                                     |
| `READ_TIMEOUT`              | `15s`              | Server read and default handler timeout.                            |
| `WRITE_TIMEOUT`             | `15s`              | Server write timeout.                                               |
| `IDLE_TIMEOUT`              | `90s`              | Keep-alive idle timeout.                                            |
| `READ_HEADER_TIMEOUT`       | `500ms`            | Header read/slow-client deadline.                                   |
| `PING_TIMEOUT`              | `15s`              | HTTP/2 ping timeout.                                                |
| `RELOAD_SHUTDOWN_TIMEOUT`   | `30s`              | Graceful shutdown/reload deadline.                                  |
| `ENABLE_SLOWLORIS_CHECK`    | `false`            | Compatibility flag; header deadlines are configured independently.  |
| `RATE_LIMIT_SIZE`           | `10000`            | LRU capacity for tracked client addresses.                          |
| `RATE_LIMIT_RATE`           | `360`              | Requests allowed per window.                                        |
| `RATE_LIMIT_WINDOW`         | `1m`               | Rate-limit window duration.                                         |
| `RATE_LIMIT_SKIP_LOCALHOST` | `true`             | Exempts loopback clients from rate limiting.                        |
| `WS_TOKEN_KEY`              | random per process | WebSocket token secret; set a stable value for production.          |
| `DATABASE_URL`              | empty              | Complete PostgreSQL pgx DSN.                                        |
| `POSTGRES_HOST`             | empty              | PostgreSQL host when no DSN is supplied.                            |
| `POSTGRES_USER`             | `postgres`         | PostgreSQL user.                                                    |
| `POSTGRES_PASSWORD`         | empty              | PostgreSQL password.                                                |
| `POSTGRES_DB`               | `postgres`         | PostgreSQL database.                                                |
| `DB_EXEC_MODE`              | empty              | pgx default query execution mode.                                   |
| `DB_MAX_CONNS`              | `8 × CPU`          | Pool maximum connection override.                                   |
| `DB_LOG_MODE`               | `sanitized`        | `off`, `blind`, `full`, or SQL-without-arguments logging.           |
| `SMTP_QUEUE_SIZE`           | `20`               | Bounded async mail queue capacity.                                  |
| `SMTP_WORKERS`              | `1`                | Number of SMTP delivery workers. See the pinned-version note below. |

The template currently pins `@myelophone/goserver` `v0.5.1`. That release reads `SMTP_WORKERS` from the wrong environment key; upgrade to the first `@myelophone/goserver` release containing the fix before relying on this override. Without the override, the default worker count remains one.

When `API_PREFIX=/api`, an application route registered as `/users` is exposed as `/api/users`. Keep route definitions prefix-free.

When `APP_ERROR_MODE=json`, default 404, validation, overload, and internal errors use JSON responses. Handler code may still choose `RenderError` or `RenderErrorJSON` explicitly.

## Development commands

The repository uses [Task](https://taskfile.dev/) and `taskfile.yml`. Install Task before using these commands:

```bash
go install github.com/go-task/task/v3/cmd/task@latest

# Windows alternative
winget install Task.Task

# macOS/Linux alternative
brew install go-task/tap/go-task
```

Air is required for live reload:

```bash
go install github.com/air-verse/air@latest
```

### All Taskfile commands

| Command               | Purpose                                                                                           | Requirement or side effect                         |
| --------------------- | ------------------------------------------------------------------------------------------------- | -------------------------------------------------- |
| `task run`            | Runs `cmd/main.go` with `APP_ENV=dev`.                                                            | Go toolchain.                                      |
| `task dev`            | Live reload with Air, development environment, and metrics enabled.                               | Requires Air.                                      |
| `task preview`        | Runs the application with `APP_ENV=prod`.                                                         | Go toolchain.                                      |
| `task build`          | Produces a stripped, trimpath, PGO-enabled binary in `tmp` with prod environment and Git version. | Go and Git; uses POSIX environment syntax.         |
| `task view -- [args]` | Runs the previously built binary with optional arguments.                                         | Run `task build` first.                            |
| `task lint`           | Runs golangci-lint.                                                                               | Requires golangci-lint.                            |
| `task format`         | Runs `go fmt ./...`.                                                                              | Modifies Go formatting.                            |
| `task vet`            | Runs `go vet ./...`.                                                                              | Go toolchain.                                      |
| `task vulncheck`      | Runs Go vulnerability analysis.                                                                   | Requires govulncheck.                              |
| `task test`           | Runs the verbose project test suite.                                                              | Go toolchain.                                      |
| `task clean`          | Clears the Go module cache.                                                                       | The next build downloads dependencies again.       |
| `task cover`          | Creates `coverage.out` and `coverage.html`, then attempts to open the report.                     | Uses `xdg-open` or `open` when available.          |
| `task check-updates`  | Lists available Go module updates without changing files.                                         | Network access.                                    |
| `task update`         | Updates Go dependencies and runs `go mod tidy`.                                                   | Modifies `go.mod` and `go.sum`.                    |
| `task ba`             | Reports struct alignment opportunities.                                                           | Requires betteralign; does not fail on findings.   |
| `task bafix`          | Applies betteralign changes.                                                                      | Modifies Go source.                                |
| `task critic`         | Runs go-critic.                                                                                   | Requires go-critic.                                |
| `task optimize`       | Runs format, lint, vet, vulnerability scan, go-critic, and betteralign fixes.                     | Requires all analysis tools and may modify source. |
| `task loadtest`       | Runs three Bombardier stages against `localhost:8080`.                                            | Start the server first; requires Bombardier.       |
| `task analyze-web`    | Builds and opens Go Size Analyzer in web mode.                                                    | Requires `gsa`.                                    |
| `task analyze-tui`    | Builds and opens Go Size Analyzer in terminal mode.                                               | Requires `gsa`.                                    |
| `task analyze`        | Builds and opens the default Go Size Analyzer view.                                               | Requires `gsa`.                                    |
| `task tag:patch`      | Creates and pushes the next patch tag.                                                            | Publishing operation and POSIX shell.              |
| `task tag:minor`      | Creates and pushes the next minor tag.                                                            | Publishing operation and POSIX shell.              |
| `task tag:major`      | Creates and pushes the next major tag.                                                            | Publishing operation and POSIX shell.              |

Install the optional tools used by those tasks:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
go install github.com/codesenberg/bombardier@latest
go install github.com/go-critic/go-critic/cmd/go-critic@latest
go install github.com/dkorunic/betteralign/cmd/betteralign@latest
go install github.com/Zxilly/go-size-analyzer/cmd/gsa@v1.12.2
go install honnef.co/go/tools/cmd/staticcheck@latest
```

Suggested maintenance workflows use only existing tasks:

```bash
# Normal code change
task format
task lint
task vet
task test

# Inspect and apply dependency updates
task check-updates
task update
git diff -- go.mod go.sum
task vet
task test

# Full static/security pass before release
task optimize
task test
task build
```

Some Taskfile commands use POSIX shell syntax. Use WSL or Git Bash on Windows when the active Task shell cannot interpret them.

## Docker and deployment

### Docker Compose

```bash
docker compose up --build
```

The default `DOCKER_PORT_BINDING=127.0.0.1:0:8080` asks Docker to choose a free host port. Discover it with:

```bash
docker compose port goserver 8080
```

For a fixed local port, set this in `.env`:

```dotenv
DOCKER_PORT_BINDING=127.0.0.1:8080:8080
HTTP_PORT=8080
```

Compose provides:

- `/healthz` health checks every 30 seconds through the bundled static `healthcheck` binary;
- restart-on-failure behavior;
- a five-attempt deployment restart policy;
- bounded JSON logs: three compressed files of up to 50 MB;
- optional `.env` loading.

The probe reads `HTTP_PORT`, requests `http://127.0.0.1:<port>/healthz` with a three-second timeout, and exits non-zero when the endpoint is unavailable or returns a non-200 status. It uses an explicit internal `User-Agent`, so the request is accepted by the protection middleware installed by `Defaults()`. No shell, curl, or wget is added to the distroless runtime image.

### Production image

The Dockerfile:

1. downloads modules in a cached builder layer;
2. creates a static Linux/amd64 binary with trimpath, PGO, stripped symbols, prod environment, and Git revision;
3. creates a separate static healthcheck binary and copies only both binaries, assets, templates, and license;
4. runs on `gcr.io/distroless/static:nonroot` from `/home/nonroot/app`;
5. exposes port 8080.

Keep `HTTP_PORT=8080` when using the supplied container mapping, or update both the container port mapping and health check together.

### Dev Container

`.devcontainer` defines an Alpine Go environment, mounts the repository at `/app`, forwards port 8080, and is configured to install Git, curl, Bash, build tools, Task, Air, golangci-lint, and govulncheck. Open the project with a Dev Containers-compatible editor to keep local tooling isolated.

> **Current template limitation:** the final cleanup command in `.devcontainer/Dockerfile` uses `go env/GOPATH` instead of `go env GOPATH`, so the image build can fail at that layer. Correct that command before depending on the Dev Container workflow.

## Updating @myelophone/goserver

Inspect available versions:

```bash
go list -m -versions github.com/myelophone/goserver
task check-updates
```

Update deliberately:

```bash
go get github.com/myelophone/goserver@latest
go mod tidy
task vet
task test
task build
```

Review the upstream [@myelophone/goserver README](https://github.com/myelophone/goserver) and release changes before updating. Commit `go.mod` and `go.sum` together.

## License

Copyright © 2026 Aliaksandr Ivanou ([aleksivanov.me](https://aleksivanov.me), [GitHub](https://github.com/aleksivanou)). All rights reserved.

This project is licensed under the **PolyForm Noncommercial License 1.0.0**. See [`LICENSE`](./LICENSE) for the full text. Commercial use is not permitted without separate permission from the copyright holder.
