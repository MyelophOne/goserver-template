package main

import (
	_ "github.com/myelophone/goserver-template/app"

	goserver "github.com/myelophone/goserver"
)

func main() {
	httpPort := goserver.GetEnv("HTTP_PORT", "8080")

	server := goserver.NewServer(httpPort)

	server.Defaults()

	tm := goserver.NewTemplateManager()
	server.TemplatesMiddleware(tm)

	server.ApplyHooks()

	server.Run()
}
