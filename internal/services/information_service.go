package services

import (
	"ecampus/database"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

type InformationService struct {
	db     *database.DB
	logger *zap.Logger
}

func NewInformationService(db *database.DB, logger *zap.Logger) *InformationService {
	return &InformationService{db: db, logger: logger}
}

// CreateInformation inserts a new information record into the database
func (s *InformationService) CreateInformation(info *database.Information) error {
	if _, err := s.db.Model(info).Insert(); err != nil { // Changed to Model and Insert for go-pg
		s.logger.Error("Failed to create information", zap.Error(err))
		return fmt.Errorf("service: failed to create information: %w", err)
	}
	s.logger.Info("Information created successfully")
	return nil
}

// GetAllInformation retrieves all information records from the database
func (s *InformationService) GetAllInformation() ([]*database.Information, error) {
	var infoList []*database.Information
	if err := s.db.Model(&infoList).Select(); err != nil { // Changed to Model and Select for go-pg
		s.logger.Error("Failed to retrieve information", zap.Error(err))
		return nil, fmt.Errorf("service: failed to retrieve information: %w", err)
	}
	return infoList, nil
}

// GetInformation retrieves an information record by its ID
func (s *InformationService) GetInformation(id snowflake.ID) (*database.Information, error) {
	var info database.Information
	if err := s.db.Model(&info).Where("id = ?", id).Select(); err != nil { // Changed to Model and Select for go-pg
		s.logger.Error("Failed to get information", zap.Error(err))
		return nil, fmt.Errorf("service: failed to get information: %w", err)
	}
	return &info, nil
}

// UpdateInformation updates an existing information record in the database
func (s *InformationService) UpdateInformation(info *database.Information) error {
	if _, err := s.db.Model(info).WherePK().Update(); err != nil { // Changed to Model and Update for go-pg
		s.logger.Error("Failed to update information", zap.Error(err))
		return fmt.Errorf("service: failed to update information: %w", err)
	}
	s.logger.Info("database.Information updated successfully")
	return nil
}

// DeleteInformation deletes an information record by its ID
func (s *InformationService) DeleteInformation(id uint64) error {
	info := &database.Information{ID: id}                          // Create an information instance to delete
	if _, err := s.db.Model(info).WherePK().Delete(); err != nil { // Changed to Model and Delete for go-pg
		s.logger.Error("Failed to delete information", zap.Error(err))
		return fmt.Errorf("service: failed to delete information: %w", err)
	}
	s.logger.Info("database.Information deleted successfully")
	return nil
}

// GetInformationByCategory retrieves all information records for a given category
func (s *InformationService) GetInformationByCategory(category string) ([]*database.Information, error) {
	var infoList []*database.Information
	if err := s.db.Model(&infoList).Where("category = ?", category).Select(); err != nil { // Changed to Model and Select for go-pg
		s.logger.Error("Failed to get information by category", zap.Error(err))
		return nil, fmt.Errorf("service: failed to get information by category: %w", err)
	}
	return infoList, nil
}
