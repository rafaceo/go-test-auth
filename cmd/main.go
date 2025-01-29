package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

var db *sqlx.DB

func InitDB(dsn string) (*sqlx.DB, error) {
	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateUser(phone, passwordHash, firstName, lastName, email string) error {
	query := `INSERT INTO users (phone, password_hash, first_name, last_name, email)
              VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, phone, passwordHash, firstName, lastName, email)
	return err
}
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

	// Проверка пароля
	if !validatePassword(user.Password) {
		http.Error(w, "Пароль не соответствует требованиям", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка при хэшировании пароля", http.StatusInternalServerError)
		return
	}

	err = CreateUser(user.Phone, string(hashedPassword), user.FirstName, user.LastName, user.Email)
	if err != nil {
		http.Error(w, "Ошибка при добавлении пользователя в базу данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно зарегистрирован"})
}

func validatePassword(password string) bool {
	// Проверка на минимум 8 символов
	if len(password) < 8 {
		return false
	}

	// Проверка на наличие хотя бы одной заглавной буквы
	hasUpper := false
	// Проверка на наличие хотя бы одной цифры
	hasDigit := false
	// Проверка на наличие хотя бы одного специального символа
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

	// Все условия должны быть выполнены
	return hasUpper && hasDigit && hasSpecial
}

func main() {

	dsn := "postgres://rafael:StrongAndSmart1@localhost:5432/testing?sslmode=disable"
	_, err := InitDB(dsn)
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных:", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v4/users", RegisterUser).Methods("POST")

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
