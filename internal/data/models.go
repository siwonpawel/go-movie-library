package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies interface {
		Insert(movie *Movie) error
		GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error)
		Get(id int64) (*Movie, error)
		Update(movie *Movie) error
		Delete(id int64) error
	}
	Permissions interface {
		GetAllForUser(int64) (Persmissions, error)
	}
	Tokens interface {
		New(userID int64, ttl time.Duration, scope string) (*Token, error)
		DeleteAllForUser(scope string, userId int64) error
	}
	Users interface {
		Insert(user *User) error
		GetForToken(tokenScope, tokenPlaintext string) (*User, error)
		GetByEmail(email string) (*User, error)
		Update(user *User) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      MovieModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
