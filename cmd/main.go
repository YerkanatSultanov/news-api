package main

import (
	_ "news-api/docs"
	"news-api/internal/app"
)

// @title News API
// @version 1.0
// @description Simple REST API for managing news with authentication
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	app.NewApp().Run()
}
