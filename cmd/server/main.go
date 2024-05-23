package main

import (
	"log"

	"github.com/Ryota-Tsunoi/task-management-api/pkg/db"
	"github.com/Ryota-Tsunoi/task-management-api/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dbInit()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.Logger.Fatal(e.Start("localhost:18080"))
}

func dbInit() {
	// データベースの初期化
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// マイグレーションの実行
	if err := db.DB.AutoMigrate(&models.Task{}); err != nil {
		log.Fatal(err)
	}

}
