package repositories

import (
	"database/sql"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(database *sql.DB) *UserRepository {
	return &UserRepository{db: database}
}

func (r *UserRepository) GetAllUsers() ([]*db.User, error) {
	rows, err := r.db.Query("SELECT id, name, email, created_at FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*db.User
	for rows.Next() {
		var user db.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (*db.User, error) {
	var user db.User
	if err := r.db.QueryRow("SELECT id, name, email, created_at FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user *db.User) (*db.User, error) {
	var id int
	if err := r.db.QueryRow(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at",
		user.Name, user.Email,
	).Scan(&id, &user.CreatedAt); err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

func (r *UserRepository) UpdateUser(id int, user *db.User) (*db.User, error) {
	result, err := r.db.Exec(
		"UPDATE users SET name = $1, email = $2 WHERE id = $3",
		user.Name, user.Email, id,
	)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, db.ErrRecordNotFound
	}
	user.ID = id
	return user, nil
}

func (r *UserRepository) DeleteUser(id int) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return db.ErrRecordNotFound
	}
	return nil
}
