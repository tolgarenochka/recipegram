package main

import "guiapi/internal/app"

// @title Recipegram Swagger API
// @version 1.0
// @description Swagger API for Golang Project Recipegram
// @host guiapi:8080
// @basePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	app.Run()
}
