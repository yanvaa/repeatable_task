package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:         normalized.Title,
		Description:   normalized.Description,
		Status:        normalized.Status,
		StartDateTime: normalized.StartDateTime,
		EndDateTime:   normalized.EndDateTime,
		RepeatRule:    normalized.RepeatRule,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:            id,
		Title:         normalized.Title,
		Description:   normalized.Description,
		Status:        normalized.Status,
		StartDateTime: normalized.StartDateTime,
		EndDateTime:   normalized.EndDateTime,
		RepeatRule:    normalized.RepeatRule,
		UpdatedAt:     s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func (s *Service) GetForDateRange(ctx context.Context, from, to time.Time) (map[time.Time][]*taskdomain.Task, error) {
	if from.After(to) {
		return nil, fmt.Errorf("%w: from date must be before to date", ErrInvalidInput)
	}

	tasks, err := s.repo.GetByDateRange(ctx, from, to)
	if err != nil {
		return nil, err
	}

	result := make(map[time.Time][]*taskdomain.Task)
	for current := from; !current.After(to); current = current.AddDate(0, 0, 1) {
		for _, task := range tasks {
			taskCopy := task
			if s.appearsOnDate(&taskCopy, current) {
				result[current] = append(result[current], &taskCopy)
			}
		}
	}

	return result, nil
}

func (s *Service) appearsOnDate(task *taskdomain.Task, datetime time.Time) bool {

	if task.StartDateTime != nil {
		if datetime.Before(*task.StartDateTime) {
			return false
		}
	}

	if task.EndDateTime != nil {
		if datetime.After(*task.EndDateTime) {
			return false
		}
	}

	if task.RepeatRule == nil {
		return true
	}

	switch task.RepeatRule.Type {
	case taskdomain.RepeatDaily:
		if task.StartDateTime == nil {
			return false
		}
		daysDiff := int(datetime.Sub(*task.StartDateTime).Hours() / 24)
		return daysDiff >= 0 && daysDiff%task.RepeatRule.Interval == 0

	case taskdomain.RepeatMonthly:
		targetDay := task.RepeatRule.DayOfMonth
		lastDay := time.Date(datetime.Year(), datetime.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		if targetDay > lastDay {
			targetDay = lastDay
		}
		return datetime.Day() == targetDay

	case taskdomain.RepeatSpecific:
		dateStr := datetime.Format("2006-01-02")
		if task.RepeatRule.SpecificDate == dateStr {
			return true
		}
		return false

	case taskdomain.RepeatParity:
		if task.RepeatRule.Parity == "even" {
			return datetime.Day()%2 == 0
		}
		return datetime.Day()%2 == 1

	default:
		return false
	}
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.StartDateTime != nil && input.EndDateTime != nil && input.StartDateTime.After(*input.EndDateTime) {
		return CreateInput{}, fmt.Errorf("%w: start date must be before end date", ErrInvalidInput)
	}

	if input.RepeatRule != nil {
		if !input.RepeatRule.Valid() {
			return CreateInput{}, fmt.Errorf("%w: %v", ErrInvalidInput, input.RepeatRule.Type)
		}

		if input.StartDateTime == nil || input.EndDateTime == nil {
			return CreateInput{}, fmt.Errorf("%w: start_datetime and end_datetime are required for recurring tasks", ErrInvalidInput)
		}
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.StartDateTime != nil && input.EndDateTime != nil && input.StartDateTime.After(*input.EndDateTime) {
		return UpdateInput{}, fmt.Errorf("%w: start date must be before end date", ErrInvalidInput)
	}

	if input.RepeatRule != nil {
		if !input.RepeatRule.Valid() {
			return UpdateInput{}, fmt.Errorf("%w: %v", ErrInvalidInput, input.RepeatRule.Type)
		}
	}

	return input, nil
}
