package cmd

import "github.com/jmoiron/sqlx"

var Db *sqlx.DB

func InitDB(dsn string) (*sqlx.DB, error) {
	var err error
	Db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return Db, nil
}

func CreateUser(phone, passwordHash, firstName, lastName, email string) error {
	query := `INSERT INTO users_profiles (phone, password_hash, first_name, last_name, email)
              VALUES ($1, $2, $3, $4, $5)`
	_, err := Db.Exec(query, phone, passwordHash, firstName, lastName, email)
	return err
}
