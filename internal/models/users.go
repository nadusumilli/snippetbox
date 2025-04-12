package models

import (
	"database/sql"
	"errors"
	appErrors "snippetbox/internal/errors"
	"strings"

	bycrptyp "golang.org/x/crypto/bcrypt"

	"github.com/lib/pq" // New import
)

type User struct {
	ID             int
	Name           string
	Username       string
	Email          string
	Password       string
	Created        string
	Updated        string
	HashedPassword string
}

type UserModel struct {
	DB *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

func (m *UserModel) Insert(name, email, password string) (int, error) {
	var id int

	hashed_password, err := bycrptyp.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}
	query := `
		INSERT INTO users (name, email, hashed_password, created, updated)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id`

	_, err = m.DB.Exec(query, name, email, string(hashed_password))
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && strings.Contains(pgErr.Constraint, "users_uc_email") {
				return 0, appErrors.ErrDuplicateEmail
			}
		}
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	query := `
		SELECT id, name, email, hashed_password, created, updated
		FROM users
		WHERE id = $1`
	row := m.DB.QueryRow(query, id)

	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.Created, &u.Updated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &u, nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, name, email, hashed_password, created, updated
		FROM users
		WHERE email = $1`
	row := m.DB.QueryRow(query, email)

	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.Created, &u.Updated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &u, nil
}

func (m *UserModel) Update(id int, name, email string) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, updated = NOW()
		WHERE id = $3`
	_, err := m.DB.Exec(query, name, email, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) UpdatePassword(id int, password string) error {
	query := `
		UPDATE users
		SET hashed_password = $1, updated = NOW()
		WHERE id = $2`
	_, err := m.DB.Exec(query, password, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Delete(id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {

	var id int
	var hashedPassword []byte

	stmt := `Select id, hashed_password FROM users WHERE email = $1`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, appErrors.ErrInvalidCredentials
		}
		return 0, err
	}

	err = bycrptyp.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if err == bycrptyp.ErrMismatchedHashAndPassword {
			return 0, appErrors.ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}
