package handler

import (
	"encoding/json"
	cmd "github.com/rafaceo/go-test-auth/cmd/db"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Phone     string `json:"phone"`
		Password  string `json:"password"`
		FirstName string `json:"firstName,omitempty"`
		LastName  string `json:"lastName,omitempty"`
		Email     string `json:"email,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	phoneRegex := `^\+[1-9]\d{1,14}$`
	rePhone := regexp.MustCompile(phoneRegex)
	if !rePhone.MatchString(user.Phone) {
		http.Error(w, "Неверный формат номера телефона", http.StatusBadRequest)
		return
	}

	if !validatePassword(user.Password) {
		http.Error(w, "Пароль не соответствует требованиям", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка при хэшировании пароля", http.StatusInternalServerError)
		return
	}

	err = cmd.CreateUser(user.Phone, string(hashedPassword), user.FirstName, user.LastName, user.Email)
	if err != nil {
		http.Error(w, "Ошибка при добавлении пользователя в базу данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно зарегистрирован"})
}

func validatePassword(password string) bool {

	if len(password) < 8 {
		return false
	}

	hasUpper := false

	hasDigit := false

	hasSpecial := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
		if strings.ContainsAny(string(char), "!@#$%^&*") {
			hasSpecial = true
		}
	}

	return hasUpper && hasDigit && hasSpecial
}
