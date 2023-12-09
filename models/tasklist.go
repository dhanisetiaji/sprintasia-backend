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
	DueDate     time.Time `json:"due_date" validate:"required"`
}

/**
 * TableName
 *
 * @return string
 */
func (TaskList) TableName() string {
	return "task_list"
}

// create ParseTimeWithLocation(request.DueDate, "Asia/Jakarta")
func ParseTimeWithLocation(date time.Time, location string) (time.Time, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}

	return time.ParseInLocation("2006-01-02 15:04:05", date.Format("2006-01-02 15:04:05"), loc)
}
