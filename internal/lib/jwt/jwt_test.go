package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sy-miller/sso-grpc/internal/domain/models"
)

func TestNewToken(t *testing.T) {
	tests := []struct {
		name      string
		user      models.User
		app       models.App
		duration  time.Duration
		wantError bool
	}{
		{
			name: "valid token creation",
			user: models.User{
				ID:    "id-1",
				Email: "test@example.com",
			},
			app: models.App{
				ID:     1,
				Secret: "test-secret-key",
			},
			duration:  time.Hour,
			wantError: false,
		},
		{
			name: "token with different duration",
			user: models.User{
				ID:    "id-2",
				Email: "user@test.com",
			},
			app: models.App{
				ID:     2,
				Secret: "another-secret",
			},
			duration:  24 * time.Hour,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := NewToken(tt.user, tt.app, tt.duration)

			if (err != nil) != tt.wantError {
				t.Errorf("NewToken() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if err == nil && tokenString == "" {
				t.Error("NewToken() returned empty token string")
				return
			}

			if err == nil {
				token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(tt.app.Secret), nil
				})

				if err != nil || !token.Valid {
					t.Error("NewToken() generated invalid token")
					return
				}

				claims := token.Claims.(*jwt.MapClaims)
				if (*claims)["uid"] != tt.user.ID || (*claims)["email"] != tt.user.Email {
					t.Error("NewToken() claims do not match input")
				}
			}
		})
	}
}
