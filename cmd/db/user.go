package cmd

import "github.com/jmoiron/sqlx"

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
	query := `INSERT INTO users_profiles (phone, password_hash, first_name, last_name, email)
              VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, phone, passwordHash, firstName, lastName, email)
	return err
}
