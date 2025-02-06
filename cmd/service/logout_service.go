package service

import (
	"encoding/json"
	"fmt"
	cmd "github.com/rafaceo/go-test-auth/cmd/db"
	"log"
	"net/http"
)

// LogoutRequest описывает структуру входных данных
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// HandlerLogout обрабатывает выход из системы
func HandlerLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Убедимся, что заголовок Content-Type = application/json
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Неверный Content-Type", http.StatusBadRequest)
		return
	}

	var req LogoutRequest

	// Декодируем JSON-запрос
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.RefreshToken == "" {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	err = revokeToken(req.RefreshToken)
	if err != nil {
		log.Println("Ошибка при удалении refreshToken:", err)
		http.Error(w, "Токен не найден", http.StatusForbidden)
		return
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Выход выполнен успешно"}`))
}

func revokeToken(token string) error {
	// Удаляем токен из базы данных
	log.Printf("Попытка удалить токен: '%s'", token)
	result, err := cmd.Db.Exec("UPDATE users_profiles SET refresh_token = NULL WHERE refresh_token = $1", token)

	if err != nil {
		log.Printf("Ошибка при удалении токена: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Println("Запись с таким токеном не найдена")
		return fmt.Errorf("токен не найден") // <-- возвращаем ошибку
	}

	log.Println("Токен успешно удалён")
	return nil
}
