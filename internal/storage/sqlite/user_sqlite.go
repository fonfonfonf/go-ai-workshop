package sqlite

import (
	"database/sql"
	"log"
	"time"

	"github.com/nattakan-n/ai-training-backend/internal/models"
	"github.com/nattakan-n/ai-training-backend/internal/repositories"
)

var _ repositories.UserRepository = (*UserSQLiteRepository)(nil)

type UserSQLiteRepository struct {
	db *sql.DB
}

func NewUserSQLiteRepository(db *sql.DB) (*UserSQLiteRepository, error) {
	repo := &UserSQLiteRepository{db: db}
	if err := repo.ensureSchema(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *UserSQLiteRepository) ensureSchema() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			name TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`)
	return err
}

func (r *UserSQLiteRepository) GetByEmail(email string) (models.User, bool, error) {
	var u models.User
	log.Printf("email %v", email)
	row := r.db.QueryRow("SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE email = ?", email)
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("err ErrNoRows %v", err)
			return models.User{}, false, nil
		}
		log.Printf("err raw %v", err)
		return models.User{}, false, err
	}
	return u, true, nil
}

func (r *UserSQLiteRepository) GetByID(id string) (models.User, bool, error) {
	var u models.User
	row := r.db.QueryRow("SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE id = ?", id)
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, false, nil
		}
		return models.User{}, false, err
	}
	return u, true, nil
}

func (r *UserSQLiteRepository) Create(user models.User) (models.User, error) {
	log.Printf("seed email created: %s", user.Email)
	_, err := r.db.Exec("INSERT INTO users (id, email, password_hash, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", user.ID, user.Email, user.PasswordHash, user.Name, time.Now(), time.Now())
	return user, err
}

func (r *UserSQLiteRepository) Update(user models.User) (models.User, error) {
	_, err := r.db.Exec("UPDATE users SET email = ?, password_hash = ?, name = ?, updated_at = ? WHERE id = ?", user.Email, user.PasswordHash, user.Name, time.Now(), user.ID)
	return user, err
}


