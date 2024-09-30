package database

import (
	"ecampus/config"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

// DB represents the database connection
type DB struct {
	*pg.DB
	logger *zap.Logger
}

// New initializes and returns a new database connection
func New(logger *zap.Logger) (*DB, error) {
	// Retrieve database URL from config
	dbConfig, err := config.New() // Handle the error
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	opts, err := pg.ParseURL(dbConfig.DatabaseURL) // Correctly use the database URL
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	db := pg.Connect(opts)

	// Check the database connection
	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	logger.Info("Database connected successfully", zap.String("database", dbConfig.DatabaseURL)) // Log the database URL for debugging
	return &DB{DB: db, logger: logger}, nil
}

// CreateSchema creates the necessary tables if they don't exist
func (db *DB) CreateSchema() error {
	// Define models
	models := []interface{}{
		(*User)(nil),
		(*Subject)(nil),
		(*Task)(nil),
		(*Resource)(nil),
		(*File)(nil),
		(*Information)(nil),
	}

	// Create the tables
	for _, model := range models {
		if err := db.Model(model).CreateTable(&orm.CreateTableOptions{IfNotExists: true}); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	db.logger.Info("Schema created successfully")
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if err := db.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	db.logger.Info("Database connection closed")
	return nil
}

// User model
type User struct {
	ID       uint64 `pg:",pk"`
	Name     string
	Email    string `pg:",unique,notnull"`
	Password string `pg:",notnull"`
	Group    int
	Major    string
	Year     int
}

// Subject model
type Subject struct {
	ID     uint64 `pg:",pk"`
	Name   string
	UserID uint64
}

// Task model
type Task struct {
	ID        uint64 `pg:",pk"`
	Title     string
	Content   string
	Completed bool `pg:",default:false"`
	SubjectID uint64
}

// Resource model
type Resource struct {
	ID        uint64 `pg:",pk"`
	Title     string
	Content   string
	Type      string
	SubjectID uint64
	AuthorID  uint64
}

// File model
type File struct {
	ID         uint64 `pg:",pk"`
	FileName   string
	FilePath   string
	ResourceID uint64
	TaskID     uint64
}

// Information model
type Information struct {
	ID       uint64 `pg:",pk"`
	Title    string
	Content  string
	AuthorID uint64 `pg:",notnull"`
}
