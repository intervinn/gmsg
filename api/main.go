package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/bwmarrin/snowflake"
	"github.com/intervinn/gmsg/api/controller"
	"github.com/intervinn/gmsg/api/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/nats-io/nats.go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	godotenv.Load()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln("failed to connect to nats: ", err)
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("POSTGRES_DSN"))))
	db := bun.NewDB(sqldb, pgdialect.New())

	err = db.Ping()
	if err != nil {
		log.Fatalln("bun: database isnt pingable:", err)
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalln("failed to create snowflake node:", err)
	}

	e := echo.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	e.Logger = logger

	e.Use(middleware.CORS("*"))
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	token := service.NewTokenService(os.Getenv("JWT_PUBLIC_KEY"), os.Getenv("JWT_PRIVATE_KEY"))

	controller.NewChannelController(e, nc, node, db, token)
	controller.NewAuthController(e, token, db, node)

	if err := e.Start(":8080"); err != nil {
		log.Fatalln("failed to start:", err)
	}
}
