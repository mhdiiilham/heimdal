package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mhdiiilham/heimdall/auth"
	"github.com/mhdiiilham/heimdall/config"
)

func main() {
	err := config.Init()
	if err != nil {
		panic("failed init config " + err.Error())
	}

	app := fiber.New()
	app.Get("/login", auth.SpotifyLogin)
	app.Get("/callback", auth.SpotifyCallback)
	app.Get("/spotify-access-token", auth.GetSpotifyAccessToken)

	log.Fatalf("%v", app.Listen(":8991"))
}
