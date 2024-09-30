package services

import (
	"ecampus/database"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

// ResourceService handles resource-related operations
type ResourceService struct {
	db     *database.DB
	logger *zap.Logger
}

func NewResourceService(db *database.DB, logger *zap.Logger) *ResourceService {
	return &ResourceService{db: db, logger: logger}
}

// CreateResource inserts a new resource into the database
func (s *ResourceService) CreateResource(resource *database.Resource) error {
	if _, err := s.db.Model(resource).Insert(); err != nil { // Changed to Model and Insert for go-pg
		s.logger.Error("Failed to create resource", zap.Error(err))
		return fmt.Errorf("service: failed to create resource: %w", err)
	}
	s.logger.Info("database.Resource created successfully")
	return nil
}

// GetResourcesBySubjectID retrieves all resources for a given subject
func (s *ResourceService) GetResourcesBySubjectID(subjectID snowflake.ID) ([]*database.Resource, error) {
	var resources []*database.Resource
	if err := s.db.Model(&resources).Where("subject_id = ?", subjectID).Select(); err != nil { // Changed to Model and Select for go-pg
		s.logger.Error("Failed to get resources", zap.Error(err))
		return nil, fmt.Errorf("service: failed to get resources: %w", err)
	}
	return resources, nil
}
