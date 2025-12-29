package auth

import (
	"context"
	"log/slog"

	"github.com/sy-miller/sso-grpc/internal/domain/models"
)

type AuthSvc struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
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
func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider) *AuthSvc {
	return &AuthSvc{
		log: log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
	}
}

func (a *AuthSvc) Login(ctx context.Context, email string, password string, appId int) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (a *AuthSvc) RegisterNewUser(ctx context.Context, email string, password string) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (a *AuthSvc) IsAdmin(ctx context.Context, userId string) (bool, error) {
	panic("not implemented") // TODO: Implement
}
