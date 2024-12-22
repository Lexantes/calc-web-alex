package main

import (
	"github.com/Lexantes/calc-web-alex/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
