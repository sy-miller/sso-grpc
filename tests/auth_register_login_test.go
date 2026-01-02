package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ssov1 "github.com/sy-miller/protos/gen/go/sso"
	"github.com/sy-miller/sso-grpc/tests/suite"
)

const (
	emptyAppId = 0
	appID      = 1
	appSecret  = "test-secret"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()
	// Register inital user
	rResp, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	require.NotEmpty(t, rResp.GetUserId())

	lResp, err := st.AuthClient.Login(
		ctx,
		&ssov1.LoginRequest{
			Email:    email,
			Password: pass,
			AppId:    appID,
		},
	)
	require.NoError(t, err)

	loginTime := time.Now()

	token := lResp.GetToken()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok, "expected jwt.MapClaims type")

	assert.Equal(t, rResp.GetUserId(), claims["uid"].(string))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL.ToTimeDuration()).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_DuplicateRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()
	// Register inital user
	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)

	// Attempt duplicate registration
	_, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user already exists")
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name     string
		email    string
		password string
		errMsg   string
	}{
		{
			name:     "empty email",
			email:    "",
			password: randomFakePassword(),
			errMsg:   "email is required",
		},
		{
			name:     "empty password",
			email:    gofakeit.Email(),
			password: "",
			errMsg:   "password is required",
		},
		{
			name:     "both empty",
			email:    "",
			password: "",
			errMsg:   "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()
	// Register inital user
	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)

	tests := []struct {
		name     string
		email    string
		password string
		appId    int32
		errMsg   string
	}{
		{
			name:     "empty email",
			email:    "",
			password: pass,
			appId:    appID,
			errMsg:   "email is required",
		},
		{
			name:     "empty password",
			email:    email,
			password: "",
			appId:    appID,
			errMsg:   "password is required",
		},
		{
			name:     "empty email and password",
			email:    "",
			password: "",
			appId:    appID,
			errMsg:   "email is required",
		},
		{
			name:     "invalid app ID",
			email:    email,
			password: pass,
			appId:    emptyAppId,
			errMsg:   "app_id is required",
		},
		{
			name:     "wrong password",
			email:    email,
			password: "wrongpassword",
			appId:    appID,
			errMsg:   "invalid credentials",
		},
		{
			name:     "unregistered email",
			email:    "nonexistent@example.com",
			password: pass,
			appId:    appID,
			errMsg:   "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appId,
			})
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}
