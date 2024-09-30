package services

import (
	"ecampus/database"
	"ecampus/internal"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"strconv"
)

type UserService struct {
	db     *database.DB
	logger *zap.Logger
}

func NewUserService(db *database.DB, logger *zap.Logger) *UserService {
	return &UserService{db: db, logger: logger}
}

func (s *UserService) CreateUser(user *database.User) error {
	// Hash the password
	hashedPassword, err := internal.HashPassword(user.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	if _, err := s.db.Model(user).Insert(); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	s.logger.Info("database.User created successfully", zap.String("userID", strconv.FormatUint(user.ID, 10)))
	return nil
}

func (s *UserService) GetUserByID(id snowflake.ID) (*database.User, error) {
	var user database.User
	if err := s.db.Model(&user).Where("id = ?", id.Int64()).Select(); err != nil {
		s.logger.Error("Failed to get user", zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (s *UserService) GetAllUsers() ([]*database.User, error) {
	var users []*database.User
	if err := s.db.Model(&users).Select(); err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	s.logger.Info("Retrieved all users", zap.Int("count", len(users)))
	return users, nil
}

func (s *UserService) UpdateUser(user *database.User) error {
	if _, err := s.db.Model(user).WherePK().Update(); err != nil {
		s.logger.Error("Failed to update user", zap.Error(err))
		return fmt.Errorf("failed to update user: %w", err)
	}
	s.logger.Info("database.User updated successfully", zap.String("userID", strconv.FormatUint(user.ID, 10)))
	return nil
}

func (s *UserService) DeleteUser(id snowflake.ID) error {
	if _, err := s.db.Model(&database.User{}).Where("id = ?", id.Int64()).Delete(); err != nil {
		s.logger.Error("Failed to delete user", zap.Error(err))
		return fmt.Errorf("failed to delete user: %w", err)
	}
	s.logger.Info("database.User deleted successfully", zap.String("userID", id.String()))
	return nil
}
