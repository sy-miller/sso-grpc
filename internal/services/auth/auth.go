package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/sy-miller/sso-grpc/internal/domain/models"
	"github.com/sy-miller/sso-grpc/internal/lib/jwt"
	"github.com/sy-miller/sso-grpc/internal/storage"
)

const pkgFn = "svcAuth"

type AuthSvc struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTtl     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (string, error)
}

type UserProvider interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId string) (bool, error)
}

type AppProvider interface {
	GetAppById(ctx context.Context, appId int) (models.App, error)
}

// New returns a new instance of Auth service
func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTtl time.Duration) *AuthSvc {
	return &AuthSvc{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTtl:     tokenTtl,
	}
}

func (a *AuthSvc) Login(ctx context.Context, email string, password string, appId int) (string, error) {
	const fn = pkgFn + ".Login"

	log := a.log.With(
		slog.String("fn", fn),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", fn, ErrInvalidCredentials)
		}

		log.Error("failed to get user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", fn, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("invalid password", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", fn, ErrInvalidCredentials)
	}

	app, err := a.appProvider.GetAppById(ctx, appId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Error("app not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", fn, ErrInvalidAppId)
		}

		log.Error("failed to get app", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTtl)
	if err != nil {
		a.log.Error("failed to generate token", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return token, nil
}

func (a *AuthSvc) RegisterNewUser(ctx context.Context, email string, password string) (string, error) {
	const fn = pkgFn + ".RegisterNewUser"

	log := a.log.With(
		slog.String("fn", fn),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("user already exists", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", fn, ErrUserExists)
		}

		log.Error("failed to save user", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	log.Info("user registered", slog.String("userId", id))

	return id, nil
}

func (a *AuthSvc) IsAdmin(ctx context.Context, userId string) (bool, error) {
	fn := pkgFn + ".IsAdmin"

	log := a.log.With(
		slog.String("fn", fn),
		slog.String("userId", userId),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("checked if user is admin", slog.Bool("isAdmin", isAdmin))

	return isAdmin, nil
}
