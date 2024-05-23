package handlers

import "github.com/Ryota-Tsunoi/task-management-api/pkg/customerrors"

type ErrorResponse struct {
	Error *customerrors.CustomError `json:"error"`
}
