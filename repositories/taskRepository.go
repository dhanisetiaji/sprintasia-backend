package repositories

import (
	"fmt"
	"gotham/infrastructures"
	"gotham/models"
	"gotham/models/scopes"
	"gotham/utils"
	"strconv"
)

type IUTaskListRepository interface {
	Migratable
	// Seedable

	GetTaskListByID(ID uint) (models.TaskList, error)
	GetTaskListByUserID(ID uint, dataid string) (tasklist []models.TaskList, err error)
	GetTaskListPagination(UserId uint, search utils.ISearch, pagination scopes.GormPager, order scopes.GormOrderer) (tasklist []models.TaskList, totalCount int64, err error)
	// Create & Save & Updates & Delete
	Create(task *models.TaskList) (err error)
	UpdateStatusByID(task *models.TaskList, param string) (err error)
	Update(task *models.TaskListUpdate, param string) (err error)
	Delete(paramId string) (err error)
}

type TaskListRepository struct {
	infrastructures.IGormDatabase
}

// // Migrate implements IUTaskListRepository.
// func (*TaskListRepository) Migrate() error {
// 	panic("unimplemented")
// }

// // Seed implements IUTaskListRepository.
// func (*TaskListRepository) Seed() error {
// 	panic("unimplemented")
// }

/**
 * Migrate
 *
 * @return error
 */

func (repository *TaskListRepository) Migrate() (err error) {
	return repository.DB().AutoMigrate(models.TaskList{})
}
func (repository *TaskListRepository) GetTaskListByID(ID uint) (tasklist models.TaskList, err error) {
	err = repository.DB().First(&tasklist, ID).Error
	return
}

func (repository *TaskListRepository) GetTaskListPagination(UserId uint, search utils.ISearch, pagination scopes.GormPager, order scopes.GormOrderer) (tasklist []models.TaskList, totalCount int64, err error) {
	query := repository.DB().Model(&models.TaskList{}).Where("created_by = ?", UserId)
	fmt.Println(search.GetSearch(), "search")
	if search.GetSearch() != "" {
		query = query.Where("name LIKE ?", "%"+search.GetSearch()+"%")
	}
	err = query.Count(&totalCount).Scopes(order.ToOrder(models.TaskList{}.TableName(), "id", "id", "created_at", "updated_at")).Scopes(pagination.ToPaginate()).Find(&tasklist).Error
	return
}

func (repository *TaskListRepository) Create(task *models.TaskList) (err error) {
	err = repository.DB().Create(&task).Error
	return
}

func (repository *TaskListRepository) UpdateStatusByID(task *models.TaskList, param string) (err error) {
	err = repository.DB().Model(&task).Where("id = ? AND created_by = ?", param, task.CreatedBy).Updates(map[string]interface{}{"status": task.Status}).Error
	return
}

func (repository *TaskListRepository) GetTaskListByUserID(ID uint, dataId string) (tasklist []models.TaskList, err error) {
	fmt.Println(ID, dataId)
	// Convert dataId to uint
	parsedID, err := strconv.ParseUint(dataId, 10, 64)
	if err != nil {
		return nil, err
	}
	err = repository.DB().Model(models.TaskList{}).Where("created_by = ? AND id = ?", ID, parsedID).Find(&tasklist).Error
	return
}

func (repository *TaskListRepository) Update(task *models.TaskListUpdate, param string) (err error) {
	err = repository.DB().Model(models.TaskList{}).Where("id = ?", param).Updates(map[string]interface{}{"name": task.Name, "description": task.Description, "due_date": task.DueDate}).Error
	return
}

func (repository *TaskListRepository) Delete(paramId string) (err error) {
	err = repository.DB().Where("id = ?", paramId).Delete(models.TaskList{}).Error
	return
}
