package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusToDo       TaskStatus = "ToDo"
	TaskStatusInProgress TaskStatus = "InProgress"
	TaskStatusDone       TaskStatus = "Done"
)

func (s TaskStatus) String() string {
	return string(s)
}

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusToDo, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

type Task struct {
	gorm.Model
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	DueDate     time.Time  `json:"dueDate"`
	Status      TaskStatus `json:"status" gorm:"default:ToDo"`
}
