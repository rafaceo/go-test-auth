package main

import (
	"github.com/joho/godotenv"
	cmd "github.com/rafaceo/go-test-auth/cmd/db"
	"github.com/rafaceo/go-test-auth/config"
	"github.com/rafaceo/go-test-auth/utils"
	"log"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	authRepoPkg "github.com/rafaceo/go-test-auth/cmd/repository/postgres"
	authServicePkg "github.com/rafaceo/go-test-auth/cmd/service"
)

func main() {

	dsn := "postgres://admin:admin@localhost:5432/testing?sslmode=disable"
	db, err := cmd.InitDB(dsn)
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных:", err)
	}
	defer db.Close()

	if err := config.LoadConfig(); err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	log.Println("Пытаемся загрузить .env...")
	erro := godotenv.Load(".env")
	if erro != nil {
		log.Println("Ошибка загрузки .env:", err)
	}

	logger := kitlog.NewLogfmtLogger(log.Writer())

	jwtSecret := config.AllConfigs.Env.JwtSecret
	log.Println("jwtSecret:", jwtSecret)
	if config.AllConfigs.Env.JwtSecret == "" {
		log.Fatal("Ошибка: jwtSecret пуст или не загружен")
	}

	authRepo := authRepoPkg.NewAuthRepository(db)
	authService := authServicePkg.NewAuthService(authRepo, jwtSecret)

	router := utils.CreateHTTPRouting(authService, logger, db)

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
