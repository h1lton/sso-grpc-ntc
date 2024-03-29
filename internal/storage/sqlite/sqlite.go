package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/h1lton/sso-grpc-ntc/internal/domain/models"
	"github.com/h1lton/sso-grpc-ntc/internal/storage"
	"github.com/h1lton/sso-grpc-ntc/pkg/operr"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New создает новый экземпляр хранилища SQLite.
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	// указываем путь до файла бд
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, operr.Error(op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(
	c context.Context,
	email string,
	passHash []byte,
) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare(
		"INSERT INTO users(email, pass_hash) VALUES (?, ?)",
	)
	if err != nil {
		return 0, operr.Error(op, err)
	}

	res, err := stmt.ExecContext(c, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) &&
			sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {

			return 0, operr.Error(op, storage.ErrUserExists)
		}

		return 0, operr.Error(op, err)
	}

	// Получаем ID созданной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, operr.Error(op, err)
	}

	return id, nil
}

func (s *Storage) User(c context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare(
		"SELECT id, email, pass_hash FROM users WHERE email = ?",
	)
	if err != nil {
		return models.User{}, operr.Error(op, err)
	}

	row := stmt.QueryRowContext(c, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, operr.Error(op, storage.ErrUserNotFound)
		}

		return models.User{}, operr.Error(op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(c context.Context, userID int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare(
		"SELECT is_admin FROM users WHERE id = ?",
	)
	if err != nil {
		return false, operr.Error(op, err)
	}

	row := stmt.QueryRowContext(c, userID)

	var is bool
	err = row.Scan(&is)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, operr.Error(op, storage.ErrUserNotFound)
		}

		return false, operr.Error(op, err)
	}

	return is, nil
}

func (s *Storage) App(c context.Context, appID int32) (models.App, error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare(
		"SELECT id, name, secret FROM apps WHERE id = ?",
	)
	if err != nil {
		return models.App{}, operr.Error(op, err)
	}

	row := stmt.QueryRowContext(c, appID)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, operr.Error(op, storage.ErrAppNotFound)
		}

		return models.App{}, operr.Error(op, err)
	}

	return app, nil
}
