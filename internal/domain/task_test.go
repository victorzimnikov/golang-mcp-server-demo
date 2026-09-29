package domain_test

import (
	"errors"
	"testing"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

func TestValidateTaskStatus(t *testing.T) {
	tests := []struct {
		name   string
		status domain.TaskStatus
		want   bool
	}{
		{
			name:   "todo status",
			status: domain.TaskStatusTodo,
			want:   true,
		},
		{
			name:   "in_progress status",
			status: domain.TaskStatusInProgress,
			want:   true,
		},
		{
			name:   "blocked status",
			status: domain.TaskStatusBlocked,
			want:   true,
		},
		{
			name:   "done status",
			status: domain.TaskStatusDone,
			want:   true,
		},
		{
			name:   "cancelled status",
			status: domain.TaskStatusCancelled,
			want:   true,
		},
		{
			name:   "unknown status",
			status: domain.TaskStatus("unknown"),
			want:   false,
		},
		{
			name:   "empty status",
			status: domain.TaskStatus(""),
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := domain.ValidateTaskStatus(test.status)

			if got != test.want {
				t.Errorf("value '%s' got '%v', want '%v'", test.status, got, test.want)
			}
		})
	}
}

func TestCanTransitionTaskStatus(t *testing.T) {
	tests := []struct {
		name       string
		fromStatus domain.TaskStatus
		toStatus   domain.TaskStatus
		want       bool
	}{
		{
			name:       "todo → in_progress",
			fromStatus: domain.TaskStatusTodo,
			toStatus:   domain.TaskStatusInProgress,
			want:       true,
		},
		{
			name:       "todo → blocked",
			fromStatus: domain.TaskStatusTodo,
			toStatus:   domain.TaskStatusBlocked,
			want:       true,
		},
		{
			name:       "todo → cancelled",
			fromStatus: domain.TaskStatusTodo,
			toStatus:   domain.TaskStatusCancelled,
			want:       true,
		},
		{
			name:       "todo → done",
			fromStatus: domain.TaskStatusTodo,
			toStatus:   domain.TaskStatusDone,
			want:       true,
		},
		{
			name:       "in_progress → done",
			fromStatus: domain.TaskStatusInProgress,
			toStatus:   domain.TaskStatusDone,
			want:       true,
		},
		{
			name:       "in_progress → blocked",
			fromStatus: domain.TaskStatusInProgress,
			toStatus:   domain.TaskStatusBlocked,
			want:       true,
		},
		{
			name:       "blocked → in_progress",
			fromStatus: domain.TaskStatusBlocked,
			toStatus:   domain.TaskStatusInProgress,
			want:       true,
		},
		{
			name:       "blocked → cancelled",
			fromStatus: domain.TaskStatusBlocked,
			toStatus:   domain.TaskStatusCancelled,
			want:       true,
		},
		{
			name:       "todo → todo",
			fromStatus: domain.TaskStatusTodo,
			toStatus:   domain.TaskStatusTodo,
			want:       false,
		},
		{
			name:       "done → in_progress",
			fromStatus: domain.TaskStatusDone,
			toStatus:   domain.TaskStatusInProgress,
			want:       false,
		},
		{
			name:       "cancelled → in_progress",
			fromStatus: domain.TaskStatusCancelled,
			toStatus:   domain.TaskStatusInProgress,
			want:       false,
		},
		{
			name:       "unknown → in_progress",
			fromStatus: domain.TaskStatus("unknown"),
			toStatus:   domain.TaskStatusInProgress,
			want:       false,
		},
		{
			name:       "todo → unknown",
			fromStatus: domain.TaskStatusTodo,
			toStatus:   domain.TaskStatus("unknown"),
			want:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := domain.CanTransitionTaskStatus(test.fromStatus, test.toStatus)

			if got != test.want {
				t.Errorf("value '%s' → '%v' got '%v', want '%v'", test.fromStatus, test.toStatus, got, test.want)
			}
		})
	}
}

func TestNewTask(t *testing.T) {
	projectIDMock := int64(1)
	titleMock := "    Test title  "
	descriptionMock := "Description"
	priorityMock := domain.TaskPriorityHigh
	sourceMock := domain.SourceCodex

	t.Run("success create task", func(t *testing.T) {
		task, err := domain.NewTask(projectIDMock, titleMock, descriptionMock, priorityMock, sourceMock)

		if err != nil {
			t.Fatalf("task error got %v, want nil", err)
		}

		if task == nil {
			t.Fatal("task got nil, want Task{}")
		}

		if task.Title != "Test title" {
			t.Errorf("task title got '%v', want 'Test title'", task.Title)
		}

		if task.Description != descriptionMock {
			t.Errorf("task description got '%s', want '%s'", task.Description, descriptionMock)
		}

		if task.Status != domain.TaskStatusTodo {
			t.Errorf("task status got '%s', want '%s'", task.Status, domain.TaskStatusTodo)
		}

		if task.Version != 1 {
			t.Errorf("task version got '%d', want '%d'", task.Version, 1)
		}

		if task.ProjectID != projectIDMock {
			t.Error("task projectID is required")
		}

		if task.Priority != priorityMock {
			t.Error("task priority is required")
		}

		if task.Source != sourceMock {
			t.Error("task source is required")
		}
	})

	t.Run("empty task title", func(t *testing.T) {
		task, err := domain.NewTask(projectIDMock, "", descriptionMock, priorityMock, sourceMock)

		if err == nil {
			t.Error("task error got 'nil', want 'error'")
		}

		if task != nil {
			t.Error("task got 'Task{}', want 'nil'")
		}
	})

	t.Run("task title with only spaces", func(t *testing.T) {
		task, err := domain.NewTask(projectIDMock, "        ", descriptionMock, priorityMock, sourceMock)

		if err == nil {
			t.Error("task error got 'nil', want 'error'")
		}

		if task != nil {
			t.Error("task got 'Task{}', want 'nil'")
		}
	})

	t.Run("task projectID equal 0", func(t *testing.T) {
		task, err := domain.NewTask(int64(0), titleMock, descriptionMock, priorityMock, sourceMock)

		if err == nil {
			t.Error("task error got 'nil', want 'error'")
		}

		if task != nil {
			t.Error("task got 'Task{}', want 'nil'")
		}
	})

	t.Run("task negative projectID", func(t *testing.T) {
		task, err := domain.NewTask(int64(-10), titleMock, descriptionMock, priorityMock, sourceMock)

		if err == nil {
			t.Error("task error got 'nil', want 'error'")
		}

		if task != nil {
			t.Error("task got 'Task{}', want 'nil'")
		}
	})

	t.Run("task priority is unknown", func(t *testing.T) {
		task, err := domain.NewTask(projectIDMock, titleMock, descriptionMock, domain.TaskPriority("unknown"), sourceMock)

		if err == nil {
			t.Error("task error got 'nil', want 'error'")
		}

		if task != nil {
			t.Error("task got 'Task{}', want 'nil'")
		}
	})

	t.Run("task source is unknown", func(t *testing.T) {
		task, err := domain.NewTask(projectIDMock, titleMock, descriptionMock, priorityMock, domain.Source("unknown"))

		if err == nil {
			t.Error("task error got 'nil', want 'error'")
		}

		if task != nil {
			t.Error("task got 'Task{}', want 'nil'")
		}
	})
}

func TestTaskChangeStatus(t *testing.T) {
	projectIDMock := int64(1)
	titleMock := "    Test title  "
	descriptionMock := "Description"
	priorityMock := domain.TaskPriorityHigh
	sourceMock := domain.SourceCodex

	t.Run("task status changed todo → in_progress", func(t *testing.T) {
		task, err := domain.NewTask(projectIDMock, titleMock, descriptionMock, priorityMock, sourceMock)
		if err != nil {
			t.Fatalf("new task: %v", err)
		}

		err = task.ChangeStatus(domain.TaskStatusInProgress)
		if err != nil {
			t.Errorf("task change status: %v", err)
		}

		if task.Status != domain.TaskStatusInProgress {
			t.Errorf("task change status got '%s', want '%s'", task.Status, domain.TaskStatusInProgress)
		}
	})

	t.Run("task status changed todo → done", func(t *testing.T) {
		task, err := domain.NewTask(projectIDMock, titleMock, descriptionMock, priorityMock, sourceMock)
		if err != nil {
			t.Fatalf("new task: %v", err)
		}

		err = task.ChangeStatus(domain.TaskStatusDone)
		if err != nil {
			t.Fatalf("task change status: %v", err)
		}

		if task.Status != domain.TaskStatusDone {
			t.Errorf("task status got '%s', want '%s'", task.Status, domain.TaskStatusDone)
		}
	})

	t.Run("task status not changed", func(t *testing.T) {
		task := domain.Task{
			Status: domain.TaskStatusDone,
		}

		err := task.ChangeStatus(domain.TaskStatusTodo)
		if !errors.Is(err, domain.ErrInvalidTaskStatusTransition) {
			t.Errorf("task change status: %v", err)
		}

		if task.Status != domain.TaskStatusDone {
			t.Errorf("task status got '%s', want '%s'", task.Status, domain.TaskStatusDone)
		}
	})
}

func TestValidateTaskPriority(t *testing.T) {
	tests := []struct {
		name     string
		priority domain.TaskPriority
		want     bool
	}{
		{
			name:     "low priority",
			priority: domain.TaskPriorityLow,
			want:     true,
		},
		{
			name:     "medium priority",
			priority: domain.TaskPriorityMedium,
			want:     true,
		},
		{
			name:     "high priority",
			priority: domain.TaskPriorityHigh,
			want:     true,
		},
		{
			name:     "urgent priority",
			priority: domain.TaskPriorityUrgent,
			want:     true,
		},
		{
			name:     "empty priority",
			priority: domain.TaskPriority(""),
			want:     false,
		},
		{
			name:     "unknown priority",
			priority: domain.TaskPriority("unknown"),
			want:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := domain.ValidateTaskPriority(test.priority)

			if got != test.want {
				t.Errorf("value '%s' got '%v', want '%v'", test.priority, got, test.want)
			}
		})
	}
}

func TestValidateSource(t *testing.T) {
	tests := []struct {
		name   string
		source domain.Source
		want   bool
	}{
		{
			name:   "codex source",
			source: domain.SourceCodex,
			want:   true,
		},
		{
			name:   "claude source",
			source: domain.SourceClaude,
			want:   true,
		},
		{
			name:   "human source",
			source: domain.SourceHuman,
			want:   true,
		},
		{
			name:   "empty source",
			source: domain.Source(""),
			want:   false,
		},
		{
			name:   "unknown source",
			source: domain.Source("unknown"),
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := domain.ValidateSource(test.source)

			if got != test.want {
				t.Errorf("value '%s' got '%v', want '%v'", test.source, got, test.want)
			}
		})
	}
}
