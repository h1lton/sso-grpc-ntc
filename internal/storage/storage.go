package storage

import (
	"errors"
)

var (
	ErrUserExists   = errors.New("пользователь уже существует")
	ErrUserNotFound = errors.New("пользователь не найден")
	ErrAppNotFound  = errors.New("приложение не найдено")
)
