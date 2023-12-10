package models

import (
	"time"
)

type TaskList struct {
	ID          uint      `gorm:"primaryKey;auto_increment" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Status      string    `gorm:"size:10;not null" json:"status"`
	DueDate     time.Time `json:"due_date"`
	Description string    `gorm:"size:255;not null" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uint      `gorm:"not null" json:"created_by"`
	TaskListID  uint      `gorm:"default:0" json:"task_list_id"`
	SubTask     []SubTask `gorm:"foreignKey:TaskListID;references:ID" json:"sub_task"`
	User        User      `gorm:"foreignKey:CreatedBy;references:ID" json:"user"`
}

type TaskListCreate struct {
	Name        string    `json:"name" validate:"required,min=4,max=50"`
	Status      string    `json:"status" validate:"required"`
	DueDate     time.Time `json:"due_date" validate:"required"`
	Description string    `json:"description"`
}

type TaskListUpdate struct {
	Name        string    `json:"name" validate:"required,min=4,max=50"`
	Description string    `json:"description" validate:"required"`
	Status      string    `json:"status" validate:"required"`
	DueDate     time.Time `json:"due_date" validate:"required"`
}

type SubTask struct {
	ID          uint      `gorm:"primaryKey;auto_increment" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Status      string    `gorm:"size:10;not null" json:"status"`
	DueDate     time.Time `json:"due_date"`
	Description string    `gorm:"size:255;not null" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	TaskListID  uint      `gorm:"not null" json:"task_list_id"`
}

type SubTaskInterfc struct {
	Name        string    `json:"name" validate:"required,min=4,max=50"`
	Status      string    `json:"status" validate:"required"`
	DueDate     time.Time `json:"due_date" validate:"required"`
	Description string    `json:"description"`
	TaskListID  uint      `json:"task_list_id" validate:"required"`
}

/**
 * TableName
 *
 * @return string
 */
func (TaskList) TableName() string {
	return "task_list"
}
func (SubTask) TableName() string {
	return "sub_task"
}

// create ParseTimeWithLocation(request.DueDate, "Asia/Jakarta")
func ParseTimeWithLocation(date time.Time, location string) (time.Time, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}

	return time.ParseInLocation("2006-01-02 15:04:05", date.Format("2006-01-02 15:04:05"), loc)
}
