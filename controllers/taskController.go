package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"gotham/models"
	"gotham/requests"
	"gotham/services"
	"gotham/viewModels"
)

type TaskController struct {
	TaskService services.ITaskService

	// TaskPolicy policies.ITaskPolicy
}

// Index godoc
// @Summary create task
// @Description
// @Tags Task
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param reqBody body models.TaskListCreate true "reqBody"
// @Success 200 {object} viewModels.Message{}
// @Failure 400 {object} viewModels.Message{}
// @Failure 401 {object} viewModels.Message{}
// @Failure 500 {object} viewModels.Message{}
// @Router /v1/restricted/task [post]
func (t TaskController) CreateTaskList(c echo.Context) (err error) {
	auth := models.ConvertUser(c.Get("auth"))

	// Request Bind And Validation
	request := models.TaskList{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &request); err != nil {
		return err
	}

	if request.Status != "ongoing" && request.Status != "done" {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"status": "status must be ongoing or done",
		}))
	}
	// Parse the time with the correct time zone (UTC)
	request.DueDate, err = models.ParseTimeWithLocation(request.DueDate, "Asia/Jakarta")
	if err != nil {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"due_date": "due date must be filled",
		}))
	}
	request.CreatedBy = auth.ID

	err = t.TaskService.CreateTaskList(&request)
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	// Response
	return c.JSON(http.StatusOK, viewModels.SuccessResponse(map[string]string{
		"message": "task created",
	}))
}

// TaskPagination godoc
// @Summary get task with pagination
// @Description
// @Tags Task
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} viewModels.Paginator{data=[]models.TaskList}
// @Failure 400 {object} viewModels.Message{}
// @Failure 401 {object} viewModels.Message{}
// @Failure 500 {object} viewModels.Message{}
// @Router /v1/restricted/task [get]
func (t TaskController) GetTaskList(c echo.Context) (err error) {
	auth := models.ConvertUser(c.Get("auth"))

	// Request Bind And Validation
	request := new(requests.TaskIndexRequest)
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, &request.QueryParams); err != nil {
		return err
	}

	var count int64
	var task []models.TaskList
	task, count, err = t.TaskService.GetTaskListPagination(auth.ID, &request.QueryParams.Find, &request.QueryParams.Pagination, &request.QueryParams.Order)
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	// Response
	return c.JSON(http.StatusOK, viewModels.SuccessResponse(viewModels.Paginator{
		TotalRecord: count,
		Records:     task,
		Limit:       request.QueryParams.Pagination.GetLimit(),
		Page:        request.QueryParams.Pagination.GetPage(),
	}))

}

// TaskById godoc
// @Summary get task by id
// @Tags Task
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "id"
// @Success 200 {object} viewModels.HTTPSuccessResponse{data=models.TaskList}
// @Failure 400 {object} viewModels.Message{}
// @Failure 500 {object} viewModels.Message{}
// @Router /v1/restricted/task/{id} [get]
func (t TaskController) GetTaskListByID(c echo.Context) (err error) {
	auth := models.ConvertUser(c.Get("auth"))

	taskList, err := t.TaskService.GetTaskListByUserID(auth.ID, c.Param("id"))
	if err != nil {
		return echo.ErrInternalServerError
	}
	if len(taskList) == 0 {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"message": "data not found",
		}))
	}
	return c.JSON(http.StatusOK, viewModels.SuccessResponse(taskList[0]))
}

// Update godoc
// @Summary update task
// @Description
// @Tags Task
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param reqBody body map[string]interface{} true "reqBody"
// @Param id path int true "id"
// @Success 200 {object} viewModels.Message{}
// @Failure 400 {object} viewModels.Message{}
// @Failure 401 {object} viewModels.Message{}
// @Failure 500 {object} viewModels.Message{}
// @Router /v1/restricted/task/{id} [put]
func (t TaskController) Update(c echo.Context) (err error) {
	auth := models.ConvertUser(c.Get("auth"))
	fmt.Println(auth.ID)

	// Request Bind And Validation
	request := models.TaskListUpdate{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &request); err != nil {
		return err
	}

	if request.Status != "ongoing" && request.Status != "done" {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"status": "status must be ongoing or done",
		}))
	}

	// check data by jwt
	fmt.Println(auth, "auth")
	taskList, err := t.TaskService.GetTaskListByUserID(auth.ID, c.Param("id"))
	if err != nil {
		return echo.ErrInternalServerError
	}
	fmt.Println(taskList, "taskList")
	if len(taskList) == 0 {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"message": "data not found",
		}))
	}

	paramId := c.Param("id")
	if paramId == "" {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"id": "id must be filled",
		}))
	}

	// Parse the time with the correct time zone (UTC)
	request.DueDate, err = models.ParseTimeWithLocation(request.DueDate, "Asia/Jakarta")
	if err != nil {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"due_date": "due date must be filled",
		}))
	}

	err = t.TaskService.Update(&request, paramId)
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	// Response
	return c.JSON(http.StatusOK, viewModels.SuccessResponse(map[string]string{
		"message": "task updated",
	}))
}

// UpdateStatus godoc
// @Summary update status task
// @Description
// @Tags Task
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param reqBody body map[string]interface{} true "reqBody"
// @Param id path int true "id"
// @Success 200 {object} viewModels.Message{}
// @Failure 400 {object} viewModels.Message{}
// @Failure 401 {object} viewModels.Message{}
// @Failure 500 {object} viewModels.Message{}
// @Router /v1/restricted/task/status/{id} [put]
func (t TaskController) UpdateTaskList(c echo.Context) (err error) {
	auth := models.ConvertUser(c.Get("auth"))

	// Request Bind And Validation
	request := models.TaskList{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &request); err != nil {
		return err
	}

	if request.Status != "done" {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"status": "status must be done",
		}))
	}

	// check data by jwt
	taskList, err := t.TaskService.GetTaskListByUserID(auth.ID, c.Param("id"))
	if err != nil {
		return echo.ErrInternalServerError
	}
	if len(taskList) == 0 {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"message": "data not found",
		}))
	}

	paramId := c.Param("id")
	fmt.Println(taskList, "f")
	if paramId == "" {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"id": "id must be filled",
		}))
	}
	request.CreatedBy = auth.ID
	err = t.TaskService.UpdateStatusListByID(&request, paramId)
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	// Response
	return c.JSON(http.StatusOK, viewModels.SuccessResponse(map[string]string{
		"message": "task status updated",
	}))
}

// Delete godoc
// @Summary delete task
// @Description
// @Tags Task
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "id"
// @Success 200 {object} viewModels.Message{}
// @Failure 500 {object} viewModels.Message{}
// @Router /v1/restricted/task/{id} [delete]
func (t TaskController) DeleteTaskList(c echo.Context) (err error) {
	auth := models.ConvertUser(c.Get("auth"))

	// check data by jwt
	taskList, err := t.TaskService.GetTaskListByUserID(auth.ID, c.Param("id"))
	if err != nil {
		return echo.ErrInternalServerError
	}
	if len(taskList) == 0 {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"message": "data not found",
		}))
	}

	paramId := c.Param("id")
	if paramId == "" {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"id": "id must be filled",
		}))
	}

	err = t.TaskService.DeleteTaskListByID(paramId)
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	// Response
	return c.JSON(http.StatusOK, viewModels.SuccessResponse(map[string]string{
		"message": "task deleted",
	}))
}

// Index godoc
// @Summary create Subtask
// @Description
// @Tags SubTask
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param reqBody body models.SubTaskInterfc true "reqBody"
// @Success 200 {object} viewModels.Message{}
// @Failure 400 {object} viewModels.Message{}
// @Failure 401 {object} viewModels.Message{}
// @Failure 500 {object} viewModels.Message{}
// @Router /v1/restricted/subtask [post]
func (t TaskController) CreateSubTask(c echo.Context) (err error) {
	auth := models.ConvertUser(c.Get("auth"))

	// Request Bind And Validation
	request := models.SubTask{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &request); err != nil {
		return err
	}

	if request.Status != "ongoing" && request.Status != "done" {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"status": "status must be ongoing or done",
		}))
	}

	// Parse the time with the correct time zone (UTC)
	request.DueDate, err = models.ParseTimeWithLocation(request.DueDate, "Asia/Jakarta")
	if err != nil {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"due_date": "due date must be filled",
		}))
	}

	// check data by jwt
	taskList, err := t.TaskService.GetTaskListByUserID(auth.ID, strconv.Itoa(int(request.TaskListID)))
	if err != nil {
		return echo.ErrInternalServerError
	}
	if len(taskList) == 0 {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"message": "task id not found",
		}))
	}

	err = t.TaskService.CreateSubTask(&request)
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	// Response
	return c.JSON(http.StatusOK, viewModels.SuccessResponse(map[string]string{
		"message": "sub task created",
	}))
}

// Update godoc
// @Summary update subtask
// @Description
// @Tags SubTask
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param reqBody body models.SubTaskInterfc true "reqBody"
// @Success 200 {object} viewModels.Message{}
// @Failure 400 {object} viewModels.Message{}
// @Failure 401 {object} viewModels.Message{}
// @Failure 500 {object} viewModels.Message{}
// @Router /v1/restricted/subtask [put]
func (t TaskController) UpdateSubTask(c echo.Context) (err error) {
	auth := models.ConvertUser(c.Get("auth"))

	// Request Bind And Validation
	request := models.SubTaskInterfc{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &request); err != nil {
		return err
	}

	if request.Status != "ongoing" && request.Status != "done" {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"status": "status must be ongoing or done",
		}))
	}

	// Parse the time with the correct time zone (UTC)
	request.DueDate, err = models.ParseTimeWithLocation(request.DueDate, "Asia/Jakarta")
	if err != nil {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"due_date": "due date must be filled",
		}))
	}

	// check data by jwt
	taskList, err := t.TaskService.GetTaskListByUserID(auth.ID, strconv.Itoa(int(request.TaskListID)))
	if err != nil {
		return echo.ErrInternalServerError
	}
	if len(taskList) == 0 {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"message": "task id not found",
		}))
	}

	paramId := c.Param("id")
	if paramId == "" {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"id": "id must be filled",
		}))
	}

	err = t.TaskService.UpdateSubTaskByID(&request, paramId)
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	// Response
	return c.JSON(http.StatusOK, viewModels.SuccessResponse(map[string]string{
		"message": "sub task updated",
	}))
}

// Delete godoc
// @Summary delete task
// @Description
// @Tags Task
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "id"
// @Success 200 {object} viewModels.Message{}
// @Failure 500 {object} viewModels.Message{}
// @Router /v1/restricted/subtask/{id} [delete]
func (t TaskController) DeleteSubTask(c echo.Context) (err error) {
	// auth := models.ConvertUser(c.Get("auth"))

	// check data by jwt
	subtask, err := t.TaskService.GetSubTaskByID(c.Param("id"))
	if err != nil {
		if err.Error() == "record not found" {
			return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
				"message": "data not found",
			}))
		}
		return echo.ErrInternalServerError
	}
	fmt.Println(subtask, "f")

	paramId := c.Param("id")
	if paramId == "" {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"id": "id must be filled",
		}))
	}

	err = t.TaskService.DeleteSubTaskByID(paramId)
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	// Response
	return c.JSON(http.StatusOK, viewModels.SuccessResponse(map[string]string{
		"message": "sub task deleted",
	}))
}
