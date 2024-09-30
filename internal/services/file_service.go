package services

import (
	"ecampus/database"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

type FileService struct {
	db     *database.DB
	logger *zap.Logger
}

func NewFileService(db *database.DB, logger *zap.Logger) *FileService {
	return &FileService{db: db, logger: logger}
}

// CreateFile inserts a new file into the database
func (s *FileService) CreateFile(file *database.File) error {
	if _, err := s.db.Model(file).Insert(); err != nil { // Insert new file
		s.logger.Error("Failed to create file", zap.Error(err))
		return fmt.Errorf("service: failed to create file: %w", err)
	}
	s.logger.Info("database.File created successfully")
	return nil
}

// GetAllFiles retrieves all files from the database
func (s *FileService) GetAllFiles() ([]*database.File, error) {
	var files []*database.File
	if err := s.db.Model(&files).Select(); err != nil { // Retrieve all files
		s.logger.Error("Failed to retrieve files", zap.Error(err))
		return nil, fmt.Errorf("service: failed to retrieve files: %w", err)
	}
	return files, nil
}

// GetFile retrieves a file by its ID
func (s *FileService) GetFile(id snowflake.ID) (*database.File, error) {
	var file database.File
	if err := s.db.Model(&file).Where("id = ?", id).Select(); err != nil { // Retrieve file by ID
		s.logger.Error("Failed to get file", zap.Error(err))
		return nil, fmt.Errorf("service: failed to get file: %w", err)
	}
	return &file, nil
}

// UpdateFile updates an existing file in the database
func (s *FileService) UpdateFile(file *database.File) error {
	if _, err := s.db.Model(file).WherePK().Update(); err != nil { // Update file by primary key
		s.logger.Error("Failed to update file", zap.Error(err))
		return fmt.Errorf("service: failed to update file: %w", err)
	}
	s.logger.Info("database.File updated successfully")
	return nil
}

// DeleteFile deletes a file by its ID
func (s *FileService) DeleteFile(id uint64) error {
	file := &database.File{ID: id}                                 // Create a file instance to delete
	if _, err := s.db.Model(file).WherePK().Delete(); err != nil { // Delete file by primary key
		s.logger.Error("Failed to delete file", zap.Error(err))
		return fmt.Errorf("service: failed to delete file: %w", err)
	}
	s.logger.Info("database.File deleted successfully")
	return nil
}

// GetFilesByResourceID retrieves all files for a given resource
func (s *FileService) GetFilesByResourceID(resourceID snowflake.ID) ([]*database.File, error) {
	var files []*database.File
	if err := s.db.Model(&files).Where("resource_id = ?", resourceID).Select(); err != nil { // Get files by resource ID
		s.logger.Error("Failed to get files by resource ID", zap.Error(err))
		return nil, fmt.Errorf("service: failed to get files by resource ID: %w", err)
	}
	return files, nil
}
