package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Ryota-Tsunoi/task-management-api/pkg/customerrors"
	"github.com/Ryota-Tsunoi/task-management-api/pkg/models"
	"github.com/Ryota-Tsunoi/task-management-api/pkg/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type HandlerTestSuite struct {
	suite.Suite
	db          *gorm.DB
	taskRepo    *repositories.TaskRepository
	e           *echo.Echo
	taskHandler *TaskHandler
}

func (suite *HandlerTestSuite) SetupSuite() {
	var err error
	suite.db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = suite.db.AutoMigrate(&models.Task{})
	if err != nil {
		panic("failed to migrate database")
	}

	suite.taskRepo = repositories.NewTaskRepository(suite.db)
	suite.e = echo.New()
	suite.e.Validator = &models.CustomValidator{Validator: validator.New()}
	suite.taskHandler = NewTaskHandler(suite.taskRepo)
}

func (suite *HandlerTestSuite) TearDownSuite() {
	suite.db.Exec("DROP TABLE tasks")
}

func (suite *HandlerTestSuite) BeforeTest(suiteName, testName string) {
	suite.db.Exec("DELETE FROM tasks")
}

func (suite *HandlerTestSuite) TestCreateTask() {
	taskJSON := `{"title":"Test Task","description":"Test Description","status":"ToDo"}`

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(taskJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)

	if assert.NoError(suite.T(), suite.taskHandler.CreateTask(c)) {
		assert.Equal(suite.T(), http.StatusCreated, rec.Code)

		var task models.Task
		err := json.Unmarshal(rec.Body.Bytes(), &task)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "Test Task", task.Title)
		assert.Equal(suite.T(), "Test Description", task.Description)
		assert.Equal(suite.T(), models.TaskStatusToDo, task.Status)
	}
}

func (suite *HandlerTestSuite) TestCreateTaskInvalidPayload() {
	taskJSON := `{"title":"","description":"Test Description","status":"ToDo"}`

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(taskJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)

	err := suite.taskHandler.CreateTask(c)
	if assert.Error(suite.T(), err) {
		httpErr, ok := err.(*echo.HTTPError)
		if assert.True(suite.T(), ok) {
			assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
		}
	}
}

func (suite *HandlerTestSuite) TestGetAllTasks() {
	task1 := models.Task{Title: "Task 1", Description: "Description 1", Status: models.TaskStatusToDo}
	task2 := models.Task{Title: "Task 2", Description: "Description 2", Status: models.TaskStatusInProgress}
	suite.db.Create(&task1)
	suite.db.Create(&task2)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)

	if assert.NoError(suite.T(), suite.taskHandler.GetAllTasks(c)) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)

		var tasks []models.Task
		err := json.Unmarshal(rec.Body.Bytes(), &tasks)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), tasks, 2)
	}
}

func (suite *HandlerTestSuite) TestGetTaskByID() {
	task := models.Task{Title: "Task 1", Description: "Description 1", Status: models.TaskStatusToDo}
	suite.db.Create(&task)

	req := httptest.NewRequest(http.MethodGet, "/tasks/"+strconv.FormatUint(uint64(task.ID), 10), nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)
	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatUint(uint64(task.ID), 10))

	if assert.NoError(suite.T(), suite.taskHandler.GetTaskByID(c)) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)

		var gotTask models.Task
		err := json.Unmarshal(rec.Body.Bytes(), &gotTask)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), task.Title, gotTask.Title)
		assert.Equal(suite.T(), task.Description, gotTask.Description)
		assert.Equal(suite.T(), task.Status, gotTask.Status)
	}
}

func (suite *HandlerTestSuite) TestGetTaskByIDNotFound() {
	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)
	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := suite.taskHandler.GetTaskByID(c)
	if assert.Error(suite.T(), err) {
		httpErr, ok := err.(*echo.HTTPError)
		if assert.True(suite.T(), ok) {
			assert.Equal(suite.T(), http.StatusNotFound, httpErr.Code)

			// エラーレスポンスの内容を確認
			resp, ok := httpErr.Message.(*ErrorResponse)
			if assert.True(suite.T(), ok) {
				assert.Equal(suite.T(), customerrors.ErrTaskNotFound, resp.Error.Code)
				assert.Equal(suite.T(), "Task not found", resp.Error.Message)
				assert.Equal(suite.T(), http.StatusNotFound, resp.Error.Status)
			}
		}
	}
}
func (suite *HandlerTestSuite) TestUpdateTask() {
	task := models.Task{Title: "Task 1", Description: "Description 1", Status: models.TaskStatusToDo}
	suite.db.Create(&task)

	updatedTaskJSON := `{"title":"Updated Task","description":"Updated Description","status":"Done"}`

	req := httptest.NewRequest(http.MethodPut, "/tasks/"+strconv.FormatUint(uint64(task.ID), 10), strings.NewReader(updatedTaskJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)
	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatUint(uint64(task.ID), 10))

	if assert.NoError(suite.T(), suite.taskHandler.UpdateTask(c)) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)

		var updatedTask models.Task
		err := json.Unmarshal(rec.Body.Bytes(), &updatedTask)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "Updated Task", updatedTask.Title)
		assert.Equal(suite.T(), "Updated Description", updatedTask.Description)
		assert.Equal(suite.T(), models.TaskStatusDone, updatedTask.Status)
	}
}

func (suite *HandlerTestSuite) TestUpdateTaskNotFound() {
	updatedTaskJSON := `{"title":"Updated Task","description":"Updated Description","status":"Done"}`

	req := httptest.NewRequest(http.MethodPut, "/tasks/999", strings.NewReader(updatedTaskJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)
	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := suite.taskHandler.UpdateTask(c)
	if assert.Error(suite.T(), err) {
		httpErr, ok := err.(*echo.HTTPError)
		if assert.True(suite.T(), ok) {
			assert.Equal(suite.T(), http.StatusNotFound, httpErr.Code)
		}
	}
}

func (suite *HandlerTestSuite) TestUpdateTaskInvalidPayload() {
	task := models.Task{Title: "Valid Task", Description: "Valid Description", Status: models.TaskStatusToDo}
	suite.db.Create(&task)

	invalidTaskJSON := `{"title":"","description":"Updated Description","status":"Done"}`

	req := httptest.NewRequest(http.MethodPut, "/tasks/"+strconv.FormatUint(uint64(task.ID), 10), strings.NewReader(invalidTaskJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)
	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatUint(uint64(task.ID), 10))

	err := suite.taskHandler.UpdateTask(c)
	if assert.Error(suite.T(), err) {
		httpErr, ok := err.(*echo.HTTPError)
		if assert.True(suite.T(), ok) {
			assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
		}
	}
}

func (suite *HandlerTestSuite) TestDeleteTask() {
	task := models.Task{Title: "Task to delete", Description: "Description to delete", Status: models.TaskStatusToDo}
	suite.db.Create(&task)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/"+strconv.FormatUint(uint64(task.ID), 10), nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)
	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatUint(uint64(task.ID), 10))

	if assert.NoError(suite.T(), suite.taskHandler.DeleteTask(c)) {
		assert.Equal(suite.T(), http.StatusNoContent, rec.Code)
	}

	var count int64
	suite.db.Model(&models.Task{}).Where("id = ?", task.ID).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *HandlerTestSuite) TestDeleteTaskNotFound() {
	req := httptest.NewRequest(http.MethodDelete, "/tasks/999", nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)
	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := suite.taskHandler.DeleteTask(c)
	if assert.Error(suite.T(), err) {
		httpErr, ok := err.(*echo.HTTPError)
		if assert.True(suite.T(), ok) {
			assert.Equal(suite.T(), http.StatusNotFound, httpErr.Code)
		}
	}
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
