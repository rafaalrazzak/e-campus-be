package services

import (
	"ecampus/database"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

type SubjectService struct {
	db     *database.DB
	logger *zap.Logger
}

func NewSubjectService(db *database.DB, logger *zap.Logger) *SubjectService {
	return &SubjectService{db: db, logger: logger}
}

// CreateSubject inserts a new subject into the database
func (s *SubjectService) CreateSubject(subject *database.Subject) error {
	if _, err := s.db.Model(subject).Insert(); err != nil { // Changed to Model and Insert for go-pg
		s.logger.Error("Failed to create subject", zap.Error(err))
		return fmt.Errorf("service: failed to create subject: %w", err)
	}
	s.logger.Info("database.Subject created successfully")
	return nil
}

// GetSubjectsByUserID retrieves all subjects for a given user
func (s *SubjectService) GetSubjectsByUserID(userID snowflake.ID) ([]*database.Subject, error) {
	var subjects []*database.Subject
	if err := s.db.Model(&subjects).Where("user_id = ?", userID).Select(); err != nil { // Changed to Model and Select for go-pg
		s.logger.Error("Failed to get subjects", zap.Error(err))
		return nil, fmt.Errorf("service: failed to get subjects: %w", err)
	}
	return subjects, nil
}

// UpdateSubject handles subject updates
func (s *SubjectService) UpdateSubject(subject *database.Subject) error {
	if _, err := s.db.Model(subject).WherePK().Update(); err != nil { // Changed to Model and Update for go-pg
		s.logger.Error("Failed to update subject", zap.Error(err))
		return fmt.Errorf("service: failed to update subject: %w", err)
	}
	s.logger.Info("database.Subject updated successfully")
	return nil
}

// DeleteSubject handles subject deletion
func (s *SubjectService) DeleteSubject(id uint64) error {
	subject := &database.Subject{ID: id}                              // Create a subject instance to delete
	if _, err := s.db.Model(subject).WherePK().Delete(); err != nil { // Changed to Model and Delete for go-pg
		s.logger.Error("Failed to delete subject", zap.Error(err))
		return fmt.Errorf("service: failed to delete subject: %w", err)
	}
	s.logger.Info("database.Subject deleted successfully")
	return nil
}
