package services

import (
	"gotham/models"
	"gotham/models/scopes"
	"gotham/repositories"
	"gotham/utils"
)

type ITaskService interface {
	GetTaskListByUserID(ID uint, dataid string) (tasklist []models.TaskList, err error)
	GetTaskListPagination(UserId uint, search utils.ISearch, pagination utils.IPagination, order utils.IOrder) (tasklist []models.TaskList, totalCount int64, err error)
	CreateTaskList(taskList *models.TaskList) (err error)
	UpdateStatusListByID(taskList *models.TaskList, param string) (err error)
	Update(task *models.TaskListUpdate, param string) (err error)
	DeleteTaskListByID(paramId string) (err error)
}

type TaskService struct {
	TaskListRepository repositories.IUTaskListRepository
}

func (service *TaskService) GetTaskListPagination(UserId uint, search utils.ISearch, pagination utils.IPagination, order utils.IOrder) (tasklist []models.TaskList, totalCount int64, err error) {
	return service.TaskListRepository.GetTaskListPagination(UserId, search, &scopes.GormPagination{Pagination: pagination.Get()}, &scopes.GormOrder{Order: order.Get()})
}

func (service *TaskService) CreateTaskList(taskList *models.TaskList) (err error) {
	return service.TaskListRepository.Create(taskList)
}

func (service *TaskService) UpdateStatusListByID(taskList *models.TaskList, param string) (err error) {
	return service.TaskListRepository.UpdateStatusByID(taskList, param)
}

func (service *TaskService) GetTaskListByUserID(ID uint, dataid string) (tasklist []models.TaskList, err error) {
	return service.TaskListRepository.GetTaskListByUserID(ID, dataid)
}

func (service *TaskService) Update(task *models.TaskListUpdate, param string) (err error) {
	return service.TaskListRepository.Update(task, param)
}

func (service *TaskService) DeleteTaskListByID(paramId string) (err error) {
	return service.TaskListRepository.Delete(paramId)
}
