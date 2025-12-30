package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sy-miller/sso-grpc/internal/domain/models"
	"github.com/sy-miller/sso-grpc/internal/storage"
)

const pkgFn = "storage.sqlite"

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = pkgFn + ".New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (string, error) {
	const fn = pkgFn + ".SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(id, email, pass_hash) VALUES(?,?,?)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	userId := uuid.NewString()
	_, err = stmt.ExecContext(ctx, userId, email, passHash)
	if err != nil {

		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return "", fmt.Errorf("%s: %w", fn, storage.ErrUserExists)
		}

		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return userId, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const fn = pkgFn + ".GetUserByEmail"
	var user models.User

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return user, fmt.Errorf("%s: %w", fn, err)
	}


	row := stmt.QueryRowContext(ctx, email)

	if err := row.Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", fn, storage.ErrUserNotFound)
		}
		return user, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId string) (bool, error) {
	const fn = pkgFn + ".IsAdmin"
	var isAdmin bool

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return isAdmin, fmt.Errorf("%s: %w", fn, err)
	}


	row := stmt.QueryRowContext(ctx, userId)

	if err := row.Scan(&isAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return isAdmin, fmt.Errorf("%s: %w", fn, storage.ErrUserNotFound)
		}
		return isAdmin, fmt.Errorf("%s: %w", fn, err)
	}

	return isAdmin, nil
}
