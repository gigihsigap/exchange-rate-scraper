package main

import (
	"backend/app"
)

func main() {
	var server app.Routes
	server.StartGin()
}
