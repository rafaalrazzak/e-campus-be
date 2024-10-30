package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/rafaalrazzak/e-campus-be/internal/domain/models"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
)

type UserService struct {
	db *database.ECampusDB
}

func NewUserService(db *database.ECampusDB) *UserService {
	return &UserService{db: db}
}

type UserFilters struct {
	Limit   int
	Offset  int
	Filters map[string]interface{}
}

type UserResponse struct {
	Users       []models.User `json:"users"`
	TotalCount  int64         `json:"total_count"`
	TotalPages  int64         `json:"total_pages"`
	CurrentPage int           `json:"current_page"`
}

func (s *UserService) GetUsers(params UserFilters) (*UserResponse, error) {
	users, err := s.fetchUsers(params)
	if err != nil {
		return nil, err
	}

	count, err := s.CountUsers(params.Filters)
	if err != nil {
		return nil, err
	}

	totalPages := (count + int64(params.Limit) - 1) / int64(params.Limit)
	currentPage := (params.Offset / params.Limit) + 1

	return &UserResponse{
		Users:       users,
		TotalCount:  count,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
	}, nil
}

func (s *UserService) fetchUsers(params UserFilters) ([]models.User, error) {
	query := s.db.QB.From("users").
		Select("users.*").
		Limit(uint(params.Limit)).
		Offset(uint(params.Offset))

	if len(params.Filters) > 0 {
		query = query.Where(goqu.Ex(params.Filters))
	}

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err := s.db.Conn.Select(&users, sqlQuery); err != nil {
		return nil, err
	}

	return users, nil
}

type UserDetails struct {
	models.BaseUser
	DepartmentName  string  `db:"department_name"`
	StudyPlanStatus *string `db:"study_plan_status"`
	StudyPlanGrade  *string `db:"study_plan_grade"`
}

func (s *UserService) GetUserByID(userID int64) (*UserDetails, error) {
	var user UserDetails
	query := s.db.QB.From("users").
		Select(
			goqu.I("users.id"),
			goqu.I("users.nim_nip"),
			goqu.I("users.name"),
			goqu.I("users.email"),
			goqu.I("users.role"),
			goqu.I("users.department_code"),
			goqu.I("users.entry_year"),
			goqu.I("users.status"),
			goqu.I("users.address"),
			goqu.I("users.photo_url"),
			goqu.I("users.created_at"),
			goqu.I("users.updated_at"),
			goqu.I("departments.name").As("department_name"),
			goqu.I("study_plans.status").As("study_plan_status"),
			goqu.I("study_plan_details.grade").As("study_plan_grade"),
		).
		LeftJoin(
			goqu.T("departments"),
			goqu.On(goqu.Ex{"users.department_code": goqu.I("departments.code")}),
		).
		LeftJoin(
			goqu.T("study_plans"),
			goqu.On(goqu.Ex{"users.id": goqu.I("study_plans.student_id")}),
		).
		LeftJoin(
			goqu.T("study_plan_details"),
			goqu.On(goqu.Ex{"study_plans.id": goqu.I("study_plan_details.study_plan_id")}),
		).
		Where(goqu.Ex{"users.id": userID})

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	if err := s.db.Conn.Get(&user, sqlQuery); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	exists, err := s.CheckUserExists("email", user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user with this email already exists")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := s.db.QB.Insert("users").Rows(user)
	sqlQuery, args, err := query.ToSQL()
	if err != nil {
		return err
	}

	_, err = s.db.Conn.Exec(sqlQuery, args...)
	return err
}

func (s *UserService) UpdateUser(userID string, updates map[string]interface{}) error {
	if email, ok := updates["email"].(string); ok {
		exists, err := s.CheckUserExists("email", email)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("email already in use")
		}
	}

	updates["updated_at"] = time.Now()

	query := s.db.QB.Update("users").
		Set(updates).
		Where(goqu.Ex{"id": userID})

	sqlQuery, args, err := query.ToSQL()
	if err != nil {
		return err
	}

	result, err := s.db.Conn.Exec(sqlQuery, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *UserService) DeleteUser(userID string) error {
	query := s.db.QB.Delete("users").Where(goqu.Ex{"id": userID})
	sqlQuery, args, err := query.ToSQL()
	if err != nil {
		return err
	}

	result, err := s.db.Conn.Exec(sqlQuery, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *UserService) CountUsers(filters map[string]interface{}) (int64, error) {
	query := s.db.QB.From("users").Select(goqu.COUNT("*"))
	if len(filters) > 0 {
		query = query.Where(goqu.Ex(filters))
	}

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		return 0, err
	}

	var total int64
	if err := s.db.Conn.Get(&total, sqlQuery); err != nil {
		return 0, err
	}

	return total, nil
}

func (s *UserService) CheckUserExists(field, value string) (bool, error) {
	query, _, err := s.db.QB.From("users").
		Where(goqu.Ex{field: value}).
		ToSQL()
	if err != nil {
		return false, err
	}

	var existingUser models.User
	err = s.db.Conn.Get(&existingUser, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
