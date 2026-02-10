package main

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredential = errors.New("invalid credentials")

type UserRepository interface {
	CreateUser(name, email, plainPassword, avatar string) (int, error)
	GetUserByEmailWithProfile(email string) (*User, error)
	GetUsers() ([]*User, error)
	GetUserByEmail(email string) (*User, error)
	Authenticate(email, password string) (int, error)
}

type SQLUserRepository struct {
	db *sql.DB
}

// NewSQLUserRepository creates a new instance of SQLUserRepository
func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) CreateUser(name, email, plainPassword, avatar string) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	userStmt, err := tx.Prepare("INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer userStmt.Close()

	hp, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	res, err := userStmt.Exec(name, email, hp)
	if err != nil {
		return 0, err
	}
	userID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	profileStmt, err := tx.Prepare("INSERT INTO profiles (user_id, avatar) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer profileStmt.Close()

	_, err = profileStmt.Exec(userID, avatar)
	if err != nil {
		return 0, err
	}
	return int(userID), nil
}

func (r *SQLUserRepository) GetUserByEmailWithProfile(email string) (*User, error) {
	query := `
	SELECT u.id, u.name, u.email, u.hashed_password, u.created_at, p.user_id, p.avatar, p.created_at
	FROM users u
	LEFT JOIN profiles p ON u.id = p.user_id
	WHERE u.email = ?`

	row := r.db.QueryRowContext(context.Background(), query, email)
	var user User
	var profile Profile
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.CreatedAt,
		&profile.UserID, &profile.Avatar, &profile.CreatedAt)
	if err != nil {
		return nil, err
	}
	user.Profile = profile
	return &user, nil
}

func (r *SQLUserRepository) GetUserByEmail(email string) (*User, error) {
	stmt := `SELECT u.id, u.name, u.email, u.hashed_password, u.created_at, p.avatar FROM users u INNER JOIN profiles p ON u.id = p.user_id WHERE u.email = ?`
	row := r.db.QueryRow(stmt, email)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.Profile.Avatar)
	if err != nil {
		return nil, err
	}
	user.Profile.UserID = user.ID
	return &user, nil
}

func (r *SQLUserRepository) Authenticate(email, password string) (int, error) {
	user, err := r.GetUserByEmail(email)
	if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredential
		}
		return 0, err
	}
	return user.ID, nil
}

func (r *SQLUserRepository) GetUsers() ([]*User, error) {
	query := `
	SELECT u.id, u.name, u.email, u.hashed_password, u.created_at, p.user_id, p.avatar, p.created_at
	FROM users u
	LEFT JOIN profiles p ON u.id = p.user_id`
	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User
	for rows.Next() {
		var user User
		var profile Profile
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.CreatedAt,
			&profile.UserID, &profile.Avatar, &profile.CreatedAt)
		if err != nil {
			return nil, err
		}
		user.Profile = profile
		users = append(users, &user)
	}
	return users, nil
}
