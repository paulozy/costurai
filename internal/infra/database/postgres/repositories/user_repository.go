package repositories

import (
	"database/sql"

	"github.com/paulozy/costurai/internal/entity"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(user *entity.User) error {
	stmt, err := r.DB.Prepare(`
		INSERT INTO users (id, email, password, name, location, created_at, updated_at)
		VALUES ($1, $2, $3, $4, ST_MakePoint($5, $6)::geography, $7, $8)
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		user.ID,
		user.Email,
		user.Password,
		user.Name,
		user.Location.Latitude,
		user.Location.Longitude,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User

	row := r.DB.QueryRow(`
		SELECT 
			id, 
			email, 
			password, 
			name, 
			ST_X(location::geometry), 
			ST_Y(location::geometry), 
			created_at, 
			updated_at
		FROM users
		WHERE email = $1
	`, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Location.Latitude,
		&user.Location.Longitude,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Exists(email string) (bool, error) {
	row := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
		)
	`, email)

	var exists bool
	err := row.Scan(&exists)

	return exists, err
}

func (r *UserRepository) FindByID(id string) (*entity.User, error) {
	var user entity.User

	row := r.DB.QueryRow(`
		SELECT *
		FROM users
		WHERE id = $1
	`, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Location.Latitude,
		&user.Location.Longitude,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return &user, err
}

func (r *UserRepository) Update(user *entity.User) error {
	stmt, err := r.DB.Prepare(`
		UPDATE users
		SET name = $2, location = ST_MakePoint($3, $4)::geography, updated_at = $5
		WHERE id = $1
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		user.ID,
		user.Name,
		user.Location.Latitude,
		user.Location.Longitude,
		user.UpdatedAt,
	)

	return err
}
