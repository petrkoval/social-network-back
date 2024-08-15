package main

import "github.com/petrkoval/social-network-back/internal/app"

func main() {
	sp := app.NewServiceProvider()
	sp.Init()

	sp.StartServer()
}
