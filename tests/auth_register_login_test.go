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

	respReg, err := st.AuthService.Register(c, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthService.Login(c, &ssov1.LoginRequest{
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
