package database

import "avito_tech_backend/internal/config"

func (storage *Storage) CreateUser(user config.User) error {
	_, err := storage.DB.Exec(`
		INSERT INTO users (user_id, email, password, role, created_at) 
		VALUES ($1, $2, $3, $4, $5)`,
		user.ID, user.Email, user.Password, user.Role, user.CreatedAt)
	return err
}

func (storage *Storage) GetUserByEmail(email string) (config.User, error) {
	var user config.User
	err := storage.DB.QueryRow(`
		SELECT user_id, password, email, role, created_at 
		FROM users 
		WHERE email = $1`, email).
		Scan(&user.ID, &user.Password, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return config.User{}, err
	}
	return user, err
}
