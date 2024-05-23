package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Ryota-Tsunoi/task-management-api/pkg/customerrors"
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

	if err := c.Bind(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrInvalidRequest,
				Message: "Invalid input: " + err.Error(),
				Status:  http.StatusBadRequest,
			},
		})
	}

	if err := c.Validate(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrInvalidRequest,
				Message: "Validation failed: " + err.Error(),
				Status:  http.StatusBadRequest,
			},
		})
	}

	if task.Status == "" {
		task.Status = models.TaskStatusToDo
	}

	if err := h.taskRepo.Create(&task); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrInternalServerError,
				Message: "Failed to create task",
				Status:  http.StatusInternalServerError,
			},
		})
	}

	return c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetAllTasks(c echo.Context) error {
	tasks, err := h.taskRepo.FindAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrInternalServerError,
				Message: "Failed to get tasks",
				Status:  http.StatusInternalServerError,
			},
		})
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrInvalidRequest,
				Message: "Invalid task ID",
				Status:  http.StatusBadRequest,
			},
		})
	}

	task, err := h.taskRepo.FindByID(uint(id))
	fmt.Println(err)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrTaskNotFound,
				Message: "Task not found",
				Status:  http.StatusNotFound,
			},
		})
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrInvalidRequest,
				Message: "Invalid task ID",
				Status:  http.StatusBadRequest,
			},
		})
	}

	existingTask, err := h.taskRepo.FindByID(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrTaskNotFound,
				Message: "Task not found",
				Status:  http.StatusNotFound,
			},
		})
	}

	var task models.Task
	if err := c.Bind(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrInvalidRequest,
				Message: "Invalid input: " + err.Error(),
				Status:  http.StatusBadRequest,
			},
		})
	}

	if err := c.Validate(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrInvalidRequest,
				Message: "Validation failed: " + err.Error(),
				Status:  http.StatusBadRequest,
			},
		})
	}

	task.ID = existingTask.ID // Ensure the ID remains the same

	if err := h.taskRepo.Update(&task); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrInternalServerError,
				Message: "Failed to update task",
				Status:  http.StatusInternalServerError,
			},
		})
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrInvalidRequest,
				Message: "Invalid task ID",
				Status:  http.StatusBadRequest,
			},
		})
	}

	if _, err := h.taskRepo.FindByID(uint(id)); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrTaskNotFound,
				Message: "Task not found",
				Status:  http.StatusNotFound,
			},
		})
	}

	if err := h.taskRepo.Delete(uint(id)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Error: &customerrors.CustomError{
				Code:    customerrors.ErrInternalServerError,
				Message: "Failed to delete task",
				Status:  http.StatusInternalServerError,
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}
