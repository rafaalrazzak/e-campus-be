package services

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/internal/constants"
	"github.com/rafaalrazzak/e-campus-be/internal/domain/models"
	"github.com/rafaalrazzak/e-campus-be/internal/utils"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/redis"
)

type AuthService struct {
	db          *database.ECampusDB
	redisClient *redis.ECampusRedisDB
	config      config.Config
}

type RegisterUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	NimNip   string `json:"nim_nip"`
	Role     string `json:"role"`
}

func NewAuthService(db *database.ECampusDB, redisClient *redis.ECampusRedisDB, cfg config.Config) *AuthService {
	return &AuthService{
		db:          db,
		redisClient: redisClient,
		config:      cfg,
	}
}

func (s *AuthService) AuthenticateUser(credentials models.User) (string, error) {
	if err := s.validateLoginInput(credentials); err != nil {
		return "", err
	}

	dbUser, err := s.getUserFromDB(credentials.Email)
	if err != nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	if err := s.verifyPassword(dbUser.Password, credentials.Password); err != nil {
		return "", err
	}

	token, err := s.createSession(dbUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GetSession(token string) (*int64, error) {
	userId, sessionId, err := s.ParseToken(token)
	if err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf(constants.Redis.SessionKey, userId, sessionId)
	_, err = s.redisClient.Client.Get(context.Background(), redisKey).Result()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid session")
	}

	return &userId, nil
}

func (s *AuthService) validateLoginInput(user models.User) error {
	if user.Email == "" || user.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email and Password are required")
	}
	return nil
}

func (s *AuthService) getUserFromDB(email string) (models.User, error) {
	var dbUser models.User
	query := s.db.QB.From("users").Where(goqu.Ex{"email": email})
	sql, _, _ := query.ToSQL()
	err := s.db.Conn.Get(&dbUser, sql)
	return dbUser, err
}

func (s *AuthService) verifyPassword(hashedPassword, inputPassword string) error {
	isValid, err := utils.VerifyData(hashedPassword, inputPassword)
	if err != nil || !isValid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}
	return nil
}

func (s *AuthService) ParseToken(token string) (int64, int64, error) {
	decrypt, err := utils.DecryptSessionToken(token, s.config)
	if err != nil {
		return 0, 0, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	userId, sessionId, err := utils.ParseSessionToken(decrypt)
	if err != nil {
		return 0, 0, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}
	return userId, sessionId, nil
}

func (s *AuthService) createSession(dbUser models.User) (string, error) {
	sessionId := utils.GenerateSessionToken()
	sessionToken := fmt.Sprintf("%d::%d", dbUser.ID, sessionId)
	token, err := utils.GenerateSessionEncryption(sessionToken, s.config)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}

	if err := s.storeSession(dbUser.ID, sessionId); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) storeSession(userID, sessionID int64) error {
	ctx := context.Background()
	redisKey := fmt.Sprintf(constants.Redis.SessionKey, userID, sessionID)

	// Store a blank session placeholder in Redis and set expiration
	if err := s.redisClient.Client.Set(ctx, redisKey, "", constants.App.SessionExpiration).Err(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to store session ID")
	}

	if err := s.redisClient.Client.Expire(ctx, redisKey, constants.App.SessionExpiration).Err(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to set session expiration")
	}

	return nil
}
