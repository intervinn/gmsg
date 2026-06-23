package main

import (
	"log"
	"os"

	"github.com/intervinn/gmsg/gateway/service"
	"github.com/intervinn/gmsg/gateway/ws"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/nats-io/nats.go"
)

func main() {
	godotenv.Load()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln("failed to connect to nats:", err)
	}

	e := echo.New()

	e.Use(middleware.CORS("*"))
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	token := service.NewTokenService(os.Getenv("JWT_PUBLIC_KEY"), "") // unable to generate tokens

	ws.New(e, nc, token)

	if err := e.Start(":8088"); err != nil {
		log.Fatalln("failed to start server:", err)
	}
}
