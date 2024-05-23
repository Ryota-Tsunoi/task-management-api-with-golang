package main

import (
	"log"

	"github.com/Ryota-Tsunoi/task-management-api/pkg/db"
	"github.com/Ryota-Tsunoi/task-management-api/pkg/handlers"
	"github.com/Ryota-Tsunoi/task-management-api/pkg/models"
	"github.com/Ryota-Tsunoi/task-management-api/pkg/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dbInit()
	defer db.Close()

	taskRepo := repositories.NewTaskRepository(db.DB)
	taskHandler := handlers.NewTaskHandler(taskRepo)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &models.CustomValidator{Validator: validator.New()}

	e.POST("/tasks", taskHandler.CreateTask)
	e.GET("/tasks", taskHandler.GetAllTasks)
	e.GET("/tasks/:id", taskHandler.GetTaskByID)
	e.PUT("/tasks/:id", taskHandler.UpdateTask)
	e.DELETE("/tasks/:id", taskHandler.DeleteTask)

	e.Logger.Fatal(e.Start("localhost:18080"))
}

func dbInit() {
	// データベースの初期化
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	// マイグレーションの実行
	if err := db.DB.AutoMigrate(&models.Task{}); err != nil {
		log.Fatal(err)
	}
}
