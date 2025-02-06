package main

import (
	cmd "github.com/rafaceo/go-test-auth/cmd/db"
	"log"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"

	"github.com/rafaceo/go-test-auth/cmd/handler"
	"github.com/rafaceo/go-test-auth/cmd/repository"
	"github.com/rafaceo/go-test-auth/cmd/service"
	"github.com/rafaceo/go-test-auth/cmd/transport/https"

	_ "github.com/lib/pq"
)

func main() {

	dsn := "postgres://admin:admin@localhost:5432/testing?sslmode=disable"
	db, err := cmd.InitDB(dsn)
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных:", err)
	}
	defer db.Close()

	// Создание репозитория
	authRepo := repository.NewAuthRepository(db)

	logger := kitlog.NewLogfmtLogger(log.Writer())
	authService := service.NewAuthService(authRepo)

	r := mux.NewRouter()

	r.HandleFunc("/api/v4/users", handler.RegisterUser).Methods("POST")

	r.HandleFunc("/api/v4/auth/refresh", service.RefreshTokenHandler).Methods("POST")

	r.HandleFunc("/api/v4/auth/logout", service.HandlerLogout).Methods("POST")

	authHandlers := https.GetAuthHandlers(authService, logger)

	for _, h := range authHandlers {
		r.HandleFunc(h.Path, func(w http.ResponseWriter, r *http.Request) {
			h.Handler.ServeHTTP(w, r)
		}).Methods(h.Methods...)
	}

	// Запускаем сервер
	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
