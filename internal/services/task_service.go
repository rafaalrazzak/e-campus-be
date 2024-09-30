package services

import (
	"ecampus/database"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

type TaskService struct {
	db     *database.DB
	logger *zap.Logger
}

func NewTaskService(db *database.DB, logger *zap.Logger) *TaskService {
	return &TaskService{db: db, logger: logger}
}

// CreateTask inserts a new task into the database
func (s *TaskService) CreateTask(task *database.Task) error {
	if _, err := s.db.Model(task).Insert(); err != nil { // Changed to Model and Insert for go-pg
		s.logger.Error("Failed to create task", zap.Error(err))
		return fmt.Errorf("service: failed to create task: %w", err)
	}
	s.logger.Info("database.Task created successfully")
	return nil
}

// GetTasksBySubjectID retrieves all tasks for a given subject
func (s *TaskService) GetTasksBySubjectID(subjectID snowflake.ID) ([]*database.Task, error) {
	var tasks []*database.Task
	if err := s.db.Model(&tasks).Where("subject_id = ?", subjectID).Select(); err != nil { // Changed to Model and Select for go-pg
		s.logger.Error("Failed to get tasks", zap.Error(err))
		return nil, fmt.Errorf("service: failed to get tasks: %w", err)
	}
	return tasks, nil
}

// CompleteTask sets the task's completion status to true
func (s *TaskService) CompleteTask(taskID uint64) error {
	task := &database.Task{ID: taskID, Completed: true}            // Create a task instance to update
	if _, err := s.db.Model(task).WherePK().Update(); err != nil { // Changed to Model and Update for go-pg
		s.logger.Error("Failed to complete task", zap.Error(err))
		return fmt.Errorf("service: failed to complete task: %w", err)
	}
	s.logger.Info("database.Task completed successfully")
	return nil
}

// UpdateTask updates an existing task in the database
func (s *TaskService) UpdateTask(task *database.Task) error {
	if _, err := s.db.Model(task).WherePK().Update(); err != nil { // Changed to Model and Update for go-pg
		s.logger.Error("Failed to update task", zap.Error(err))
		return fmt.Errorf("service: failed to update task: %w", err)
	}
	s.logger.Info("database.Task updated successfully")
	return nil
}

// DeleteTask handles task deletion
func (s *TaskService) DeleteTask(id uint64) error {
	task := &database.Task{ID: id}                                 // Create a task instance to delete
	if _, err := s.db.Model(task).WherePK().Delete(); err != nil { // Changed to Model and Delete for go-pg
		s.logger.Error("Failed to delete task", zap.Error(err))
		return fmt.Errorf("service: failed to delete task: %w", err)
	}
	s.logger.Info("database.Task deleted successfully")
	return nil
}
