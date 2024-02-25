package auth

import (
	"context"
	"errors"
	"github.com/h1lton/sso-grpc-ntc/internal/services/auth"
	ssov1 "github.com/h1lton/sso-grpc-ntc/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const emptyValue = 0

type Auth interface {
	Login(
		c context.Context,
		email string,
		password string,
		appID int32,
	) (token string, err error)
	Register(
		c context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(c context.Context, userID int64) (bool, error)
}

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(server *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(server, &ServerAPI{auth: auth})
}

// Обработчики...

func (s *ServerAPI) Login(
	c context.Context,
	r *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(r); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(c, r.GetEmail(), r.GetPassword(), r.GetAppId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(
				codes.InvalidArgument,
				"Неправильный email или пароль",
			)
		}
		if errors.Is(err, auth.ErrInvalidAppID) {
			return nil, status.Error(
				codes.InvalidArgument,
				"неверный id приложения",
			)
		}

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *ServerAPI) Register(
	c context.Context,
	r *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(r); err != nil {
		return nil, err
	}

	userID, err := s.auth.Register(c, r.GetEmail(), r.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(
				codes.AlreadyExists,
				"Пользователь уже существует",
			)
		}

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.RegisterResponse{UserId: userID}, nil
}

func (s *ServerAPI) IsAdmin(
	c context.Context,
	r *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(r); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(c, r.GetUserId())
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"Пользователь не найден",
			)
		}

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

// Валидаторы...

func validateLogin(r *ssov1.LoginRequest) error {
	if r.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email не указан")
	}

	if r.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "пароль не указан")
	}

	if r.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app id не указан")
	}

	return nil
}

func validateRegister(r *ssov1.RegisterRequest) error {
	if r.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email не указан")
	}

	if r.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "пароль не указан")
	}

	return nil
}

func validateIsAdmin(r *ssov1.IsAdminRequest) error {
	if r.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user id не указан")
	}

	return nil
}
