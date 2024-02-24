package operr

import "fmt"

// Error оборачивает ошибку операции в формат "op: err"
//
// Пример: "Auth.Login: недействительные учетные данные"
func Error(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
