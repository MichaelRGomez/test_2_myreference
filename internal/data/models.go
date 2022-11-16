// Filename: test2/internal/data/models.go
package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// A wrapper for our data models
type Models struct {
	//References ReferenceModel
	Users UserModel
}

// create a new model
func NewModels(db *sql.DB) Models {
	return Models{
		//References: ReferenceModel{DB: db},
		Users: UserModel{DB: db},
	}
}
