package auth

import (
	"context"
	"errors"
	"github.com/h1lton/sso-grpc-ntc/internal/domain/models"
	"github.com/h1lton/sso-grpc-ntc/internal/jwt"
	"github.com/h1lton/sso-grpc-ntc/internal/storage"
	"github.com/h1lton/sso-grpc-ntc/pkg/logger/sl"
	"github.com/h1lton/sso-grpc-ntc/pkg/operr"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		c context.Context,
		email string,
		passHash []byte,
	) (userID int64, err error)
}

type UserProvider interface {
	User(c context.Context, email string) (models.User, error)
	IsAdmin(c context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(c context.Context, appID int32) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("недействительные учетные данные")
	ErrUserExists         = errors.New("пользователь уже существует")
	ErrUserNotFound       = errors.New("пользователь не найден")
	ErrInvalidAppID       = errors.New("неверный id приложения")
)

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: userProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(
	c context.Context,
	email string,
	password string,
	appID int32,
) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("попытка войти в систему пользователя")

	user, err := a.usrProvider.User(c, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("пользователь не найден", sl.Err(err))

			return "", operr.Error(op, ErrInvalidCredentials)
		}

		log.Error("не удалось получить пользователя", sl.Err(err))

		return "", operr.Error(op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		log.Info("неверный пароль", sl.Err(err))

		return "", operr.Error(op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(c, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("приложение не найдено", sl.Err(err))

			return "", operr.Error(op, ErrInvalidAppID)
		}

		log.Error("не удалось получить приложение", sl.Err(err))

		return "", operr.Error(op, err)
	}

	log.Info("пользователь успешно вошел в систему")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("не удалось сгенерировать JWT-токен", sl.Err(err))

		return "", operr.Error(op, err)
	}

	return token, nil
}

// Register регистрирует пользователя и возвращает его ID.
func (a *Auth) Register(
	c context.Context,
	email string,
	password string,
) (int64, error) {
	const op = "auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("регистрация пользователя")

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Error("не удалось сгенерировать хэш пароля", sl.Err(err))

		return 0, operr.Error(op, err)
	}

	id, err := a.usrSaver.SaveUser(c, email, passwordHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("пользователь уже существует", sl.Err(err))

			return 0, operr.Error(op, ErrUserExists)
		}
		log.Error("не удалось сохранить пользователя", sl.Err(err))

		return 0, operr.Error(op, err)
	}

	log.Info("пользователь зарегистрирован")

	return id, err
}

func (a *Auth) IsAdmin(
	c context.Context,
	userID int64,
) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)

	log.Info("проверка")

	is, err := a.usrProvider.IsAdmin(c, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("пользователь не найден", sl.Err(err))

			return false, operr.Error(op, ErrUserNotFound)
		}

		log.Error("не удалось проверить пользователя", sl.Err(err))

		return false, operr.Error(op, err)
	}

	log.Info(
		"проверенно",
		slog.Bool("is_admin", is),
	)

	return is, nil
}
