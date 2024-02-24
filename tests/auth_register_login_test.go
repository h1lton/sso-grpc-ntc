package tests

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ssov1 "sso-grpc-ntc/pkg/api"
	"sso-grpc-ntc/tests/suite"
	"testing"
	"time"
)

const (
	appID         = 1
	appSecret     = "test-secret"
	pssDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	c, st := suite.New(t)

	email := gofakeit.Email()
	password := randomPassword()

	respReg, err := st.AuthClient.Register(c, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(c, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(
		token, func(token *jwt.Token) (interface{}, error) {
			return []byte(appSecret), nil
		},
	)
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1
	assert.InDelta(
		t,
		loginTime.Add(st.Cfg.TokenTTL).Unix(),
		claims["exp"].(float64),
		deltaSeconds,
	)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	c, st := suite.New(t)

	email := gofakeit.Email()
	password := gofakeit.Email()

	respReg, err := st.AuthClient.Register(c, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(c, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "Пользователь уже существует")
}

func TestRegister_FailCases(t *testing.T) {
	c, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Регистрация с пустым паролем",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "пароль не указан",
		},
		{
			name:        "Регистрация с пустым email",
			email:       "",
			password:    randomPassword(),
			expectedErr: "email не указан",
		},
		{
			name:        "Регистрация со всеми пустыми полями",
			email:       "",
			password:    "",
			expectedErr: "email не указан",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(c, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

		})
	}
}

func randomPassword() string {
	return gofakeit.Password(
		true,
		true,
		true,
		true,
		false,
		pssDefaultLen,
	)
}
