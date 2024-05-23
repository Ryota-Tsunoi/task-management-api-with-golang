package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ryota-Tsunoi/task-management-api/pkg/models"
	"github.com/Ryota-Tsunoi/task-management-api/pkg/repositories"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	taskRepo *repositories.TaskRepository
}

func NewTaskHandler(taskRepo *repositories.TaskRepository) *TaskHandler {
	return &TaskHandler{taskRepo: taskRepo}
}

func (h *TaskHandler) CreateTask(c echo.Context) error {
	var task models.Task

	// リクエストボディをtaskにバインド
	if err := c.Bind(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	// バリデーション
	if err := c.Validate(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	// TaskStatusが空の場合のデフォルト設定
	if task.Status == "" {
		task.Status = models.TaskStatusToDo
	}

	// タスクの作成
	if err := h.taskRepo.Create(&task); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetAllTasks(c echo.Context) error {
	tasks, err := h.taskRepo.FindAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}

	task, err := h.taskRepo.FindByID(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Task not found")
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}

	existingTask, err := h.taskRepo.FindByID(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Task not found")
	}

	var task models.Task
	if err := c.Bind(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	if err := c.Validate(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	task.ID = existingTask.ID // Ensure the ID remains the same

	if err := h.taskRepo.Update(&task); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}

	if _, err := h.taskRepo.FindByID(uint(id)); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Task not found")
	}

	if err := h.taskRepo.Delete(uint(id)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
