package server

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func (d *dbWiz) auth(username, password string) (int, error) {
	var storedPassword string
	var userID int

	err := d.dbWizard.QueryRow("SELECT user_id, password_hash FROM users WHERE username = $1", username).Scan(&userID, &storedPassword)
	if err != nil {
		log.Printf("No users with this username")
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		log.Printf("Invalid password: %v\n", err)
		return 0, err
	}

	return userID, nil
}

func (d *dbWiz) reg(username, email string, hashedPassword []byte) error {
	// Вставка пользователя в базу данных
	_, err := d.dbWizard.Exec("INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
		username, email, hashedPassword)
	if err != nil {
		log.Printf("Error inserting user into the database: %v\n", err)
		return err
	}

	return nil
}
