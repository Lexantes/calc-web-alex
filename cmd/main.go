package main

import (
	"calc-web-alex/internal/application"
	// "github.com/Lexantes/calc-web-alex/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
