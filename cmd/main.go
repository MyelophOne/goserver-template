//go:build !webcli

package main

import (
	goserver "github.com/myelophone/goserver"
	_ "github.com/myelophone/goserver-template/app"
	logic "github.com/myelophone/goserver/web/runtime"
)

func main() {
	httpPort := goserver.GetEnv("HTTP_PORT", "8080")

	server := goserver.NewServer(httpPort)

	server.Defaults()

	tm := goserver.NewTemplateManager()
	server.TemplatesMiddleware(tm)

	server.ApplyHooks()

	defaultLanguage := server.Config.I18nDefaultLanguage
	languages := server.Config.I18nLanguages
	webConfig, err := logic.UseRuntimeConfig()
	if err != nil {
		server.Logger.Fatal(err)
	}
	if webConfig.Runtime.Enabled {
		if webConfig.DefaultLocale != "" {
			defaultLanguage = webConfig.DefaultLocale
		}
		if len(webConfig.Locales) > 0 {
			languages = webConfig.Locales
		}
	}
	if _, err := server.NewI18n(defaultLanguage, languages); err != nil {
		server.Logger.Fatal(err)
	}

	if err := server.EnableWeb(); err != nil {
		server.Logger.Fatal(err)
	}

	server.Run()
}
