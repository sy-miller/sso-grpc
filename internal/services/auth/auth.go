package auth

import "context"

type AuthSvc struct{}

func New() *AuthSvc {
	return &AuthSvc{}
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
