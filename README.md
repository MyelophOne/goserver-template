# @myelophone/goserver-template

A starting point for your own Go application, powered by [goserver](https://github.com/myelophone/goserver). This repository owns your application code, configuration and deployment. The server framework and the base GOSH web layer come from the goserver dependency; you do not need to copy them into your project.

[Create a repository from this template](https://github.com/myelophone/goserver-template/generate) · [goserver README and API documentation](https://github.com/myelophone/goserver/blob/main/README.md)

This README explains how to use the template. For routing, middleware, sessions, databases, caching, WebSockets and the full web-framework reference, use the upstream documentation.

## Quick start

You need Go matching `go.mod`, Node.js for web tooling, and Task. Air is optional for live reload. PostgreSQL and Redis are only needed if your application uses them.

Create a repository using the template button, or clone it:

```bash
git clone https://github.com/myelophone/goserver-template.git my-app
cd my-app
```

Prepare your environment on Linux/macOS:

```bash
cp .env.example .env
```

Or in PowerShell:

```powershell
Copy-Item .env.example .env
```

Review `.env`, especially `HTTP_PORT`, trusted origins and secret values. Replace the example `SESSION_KEY`, `JWT_SECRET` and `METRICS_SECRET` before deployment. Do not commit secrets.

For GOSH pages, set `runtime.enabled` to `true` in `websettings.json`, or add this to `.env`:

```dotenv
MYELOPHONE_WEB_ENABLED=true
```

The checked-in web settings have web rendering disabled. Enable it explicitly if you want the inherited home page and other web pages.

Start the application:

```bash
task run
```

Open `http://localhost:8080`, or the port configured in `.env`. Web-enabled projects inherit the base goserver home page even when local `web/pages` is absent or empty. `/healthz` is available for health checks.

`task run` generates handler bindings before starting. Its setup dependency installs the web builders automatically on first use; `task setup` can also be run separately. The template's current run/build workflows invoke web tooling even when rendering is disabled, so they still require Node.js.

## Make it your application

Change the module path in `go.mod` to your repository path and update imports referring to `github.com/myelophone/goserver-template`, including the `app` import in `cmd/main.go`. Keep imports of `github.com/myelophone/goserver` unchanged. Regenerate bindings after renaming; generated files should not be edited manually.

The relevant directories are:

```text
app/                         application routes, hooks and business logic
cmd/main.go                  application startup
cmd/web_cli.go               build-tagged generation/build CLI
cmd/healthcheck/              container health-check helper
websettings.json             web rendering, SEO and feature configuration
websettings.Development.json optional development overrides
websettings.Production.json  optional production overrides
web/                         only your local additions and overrides
assets/                      only your local public files and overrides
templates/                   optional Go html/template pages and layouts
.env                         local environment and secrets, not committed
taskfile.yml                 development and build commands
Dockerfile                   production container build
docker-compose.yml           container environment, ports and health checks
internal/goservergen/         generated bindings/resources, not committed
cmd/web_import_gen.go         generated importer, not committed
dist/                        production distribution, not committed
tmp/                         temporary build inputs and reports
```

The optional directories need not contain package stubs or duplicate base files.

## Add backend routes

`cmd/main.go` creates the server, installs defaults, loads the `app` package, applies application hooks, initializes i18n, attaches web rendering when enabled, and starts the server. The application and web layer share one listener.

Add a new file such as `app/routes.go` rather than putting business logic into the entrypoint:

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

Explicit goserver routes take priority over file-based web pages. If you configure `API_PREFIX=/api`, register `/hello` instead of `/api/hello` to avoid repeating the prefix.

See [goserver routing](https://github.com/myelophone/goserver/blob/main/README.md#routing-and-route-groups) and [middleware](https://github.com/myelophone/goserver/blob/main/README.md#middleware) for the full API and security configuration.

## Web layers: add only what you customize

goserver supplies the base pages, components, layouts, CSS, stores, plugins, content, teleports, system runtime/templates and file-based server handlers. Your files are merged with that layer by relative path:

- A new path adds a file.
- The same relative path replaces the base file completely.
- Removing a local replacement restores the base file.
- An empty or missing local directory does not disable inheritance.

For example, `web/pages/index.gosh` replaces the base home page, while `web/pages/about.gosh` adds `/about`. Removing the local `index.gosh` restores the inherited home page. Base 404/error pages are inherited too. The welcome page/component are demonstration defaults, not mandatory application code.

Create `web/pages/index.gosh` for a minimal custom home page:

```html
<!-- @layout none -->
<template>
 <main><h1>My application</h1></main>
</template>
```

Omit the layout directive to use the configured default layout. File-based routing supports `index.gosh`, `[id].gosh` and `[...all].gosh`.

You do not need a local `web/system` directory. It is inherited; create individual files there only when intentionally overriding framework resources. For application styling, prefer `web/css/default.css`, which follows the system CSS in the cascade.

### Local Go handlers

Go packages are not merged like source files. Referenced handler files resolve to your local package when the referenced file exists in your project, otherwise to the base goserver package. Local code must compile independently and explicitly import any base functionality it needs.

For page/component handlers, use `package logic` in your own `web/logic/*.go` files and import `github.com/myelophone/goserver/web/runtime` directly. The alias-only `web/logic/logic.go` is unnecessary unless your code relies on its aliases.

File-based endpoints need real handler files, not a `server.go` placeholder. For example, `web/server/hello.get.go`:

```go
package server

import runtime "github.com/myelophone/goserver/web/runtime"

func Hello(event *runtime.Event) error {
	return event.JSON(map[string]string{"message": "Hello from a web handler"})
}
```

This adds `GET /hello`. Same-path local files override base handlers. `task run` / `task build` regenerate their registrations.

Local `web/modules` packages are not automatically discovered: import them explicitly from application code so their registration runs. An unimported empty `modules.go` is unnecessary.

See [goserver web documentation](https://github.com/myelophone/goserver/blob/main/README.md#nuxt-like-web-application) and [layer rules](https://github.com/myelophone/goserver/blob/main/README.md#base-layer-and-consumer-overrides) for server directives, components, layouts, client code and extension contracts.

## Public assets and configuration

Place public files in `assets/`. An `assets/icon.svg` file overrides the same base asset; a new `assets/images/product.png` adds `/assets/images/product.png`. Existing public files are also accessible at root paths such as `/icon.svg`.

Development reads the resolved goserver dependency's base assets from disk, with your local files taking priority. Production merges the layers into `dist/assets`. Images, favicons, downloads and other public files are not embedded in the executable. Only the resolved `robots.txt`, if present, is embedded by the web production build. The names/paths selected in web settings are not limited to the base demonstration assets.

Production can omit unused files. For dynamically selected assets without detectable static references, add their paths relative to `assets/` to `build.includeFiles`:

```json
{
 "build": {
  "includeFiles": ["downloads/catalog.pdf", "images/product.png"]
 }
}
```

`.env` controls listener/environment settings and secrets. `websettings.json` controls non-secret web features; optional `websettings.Development.json` and `websettings.Production.json` provide environment overrides. Environment variables have final precedence. Production embeds resolved web configuration, so configuration JSON files need not accompany the executable; deployment secrets remain external.

See [goserver configuration](https://github.com/myelophone/goserver/blob/main/README.md#configuration) for available environment variables and the [web reference](https://github.com/myelophone/goserver/blob/main/README.md#nuxt-like-web-application) for web settings. Legacy `templates/` uses Go `html/template`; it is separate from GOSH `web/pages`.

## Development commands and generated files

| Command                                | Purpose                                                            |
| -------------------------------------- | ------------------------------------------------------------------ |
| `task setup`                           | Install/reuse the cached esbuild, Tailwind and PostCSS builders.   |
| `task run`                             | Generate bindings and run with development settings.               |
| `task dev`                             | Generate bindings and start Air live reload; requires Air.         |
| `task preview`                         | Run source with production settings; not the built distribution.   |
| `task build`                           | Create the production distribution in `dist/`.                     |
| `task server -- [args]`                | Run the built executable, forwarding arguments.                    |
| `task web:generate`                    | Regenerate page/component and server-handler bindings.             |
| `task web:asset-report`                | Write the compiled asset report to `tmp/web-assets.json`.          |
| `task web:audit`                       | Report public assets without detectable static references.         |
| `task web:clean`                       | Remove generated distribution, bindings/imports and audit reports. |
| `task profile`                         | Build with web `Server-Timing` profiling enabled.                  |
| `task format`, `task vet`, `task test` | Format, analyze and test your Go code.                             |

`setup:client` and `setup:tailwind` are compatibility aliases for `task setup`. Builders are installed into the user cache, not into a new local `web/system`. Set `MYELOPHONE_TOOLCHAIN_CACHE` to change the cache root. Local tool configuration overrides get their own cache fingerprint.

Do not commit `internal/goservergen/` or `cmd/web_import_gen.go`; both are ignored and recreated automatically. No committed generated-package stub is needed. Always compile the entire `./cmd` package, not the single `cmd/main.go` file, so the generated importer is included. The CLI bootstrap uses the `webcli` tag and deliberately excludes that importer.

For manual workflows, set your environment explicitly; direct Go commands do not get Task's `.env` loading or environment selection:

```bash
go run -tags webcli ./cmd generate
go run ./cmd
```

## Production and Docker

The production builder compiles the application and its goserver dependency with `-tags myelophone_prod`. The plain-server build target uses the same tag. With a goserver version containing this separation, production excludes development sources/playground, embedded Yarn and tooling configuration, installers, CSS/JS compilers, code generation, and the web CLI. Development/build tools are provided by the dependency without copying a `web/system` tree into this project. An ordinary untagged `go build` with `APP_ENV=prod` is not equivalent to this production build.

Build and test the distribution:

```bash
task build
task server
```

The result is `dist/goserver` (`dist/goserver.exe` on Windows) plus `dist/assets/`. Deploy both together with production secrets/environment supplied externally. The production runtime resolves relative public asset roots beside the executable, so images and other disk assets do not depend on the process working directory. Standalone execution needs no Go, Node.js, Tailwind or project `web/` source tree. Removing an ordinary asset from disk makes it unavailable; there is no embedded image fallback.

`task server` launches the executable with `dist` as its working directory, matching deployment. You can also start the executable directly:

```bash
cd dist
./goserver
```

In PowerShell use `Set-Location dist` and `.\goserver.exe`.

For a container deployment:

```bash
docker compose up --build
```

The builder uses Go and Node to compile the web distribution. The final container is non-root Alpine with the executable and disk assets; it has no Go/Node toolchain. Compose reads `.env` and configures health checks, restart behavior and bounded logs. Review `DOCKER_PORT_BINDING`: the example binds a random local host port, so choose an explicit mapping when needed, for example `127.0.0.1:8080:8080`.

## Update the framework, not a copied web tree

goserver is selected by `go.mod`. Use a published release/revision containing the layer implementation described above. Editing a separate goserver working tree does not update the dependency used by this project. No module-cache copying, permanent `replace` or workspace is required for normal use.

To update deliberately:

```bash
go get github.com/myelophone/goserver@<version>
go mod tidy
task web:generate
task vet
task test
task build
```

Replace `<version>` with the release you selected, review upstream changes, and commit `go.mod` / `go.sum` together. The welcome component displays the goserver dependency version, not your application's Git version.

For a separate Nuxt frontend, use [nuxt-template](https://github.com/myelophone/nuxt-template) and expose application APIs through HTTP/WebSockets; GOSH rendering is optional. Configure trusted origins and cookie/CSRF policy for your actual deployment rather than weakening them globally.

## License

[PolyForm Noncommercial License 1.0.0](./LICENSE). Review the license before using this template commercially.
