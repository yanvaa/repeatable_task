package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Status        taskdomain.Status      `json:"status"`
	StartDateTime *time.Time             `json:"start_datetime,omitempty"`
	EndDateTime   *time.Time             `json:"end_datetime,omitempty"`
	RepeatRule    *taskdomain.RepeatRule `json:"repeat_rule,omitempty"`
}

type taskDTO struct {
	ID            int64                  `json:"id"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Status        taskdomain.Status      `json:"status"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	StartDateTime *time.Time             `json:"start_datetime,omitempty"`
	EndDateTime   *time.Time             `json:"end_datetime,omitempty"`
	RepeatRule    *taskdomain.RepeatRule `json:"repeat_rule,omitempty"`
}

type calendarResponseDTO struct {
	Date  string    `json:"date"`
	Tasks []taskDTO `json:"tasks"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:            task.ID,
		Title:         task.Title,
		Description:   task.Description,
		Status:        task.Status,
		CreatedAt:     task.CreatedAt,
		UpdatedAt:     task.UpdatedAt,
		StartDateTime: task.StartDateTime,
		EndDateTime:   task.EndDateTime,
		RepeatRule:    task.RepeatRule,
	}
}
