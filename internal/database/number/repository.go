package number

import (
	"database/sql"
	"github.com/delapaska/EgorExplain/internal/models/number"
)

type repository struct {
	db *sql.DB
}

type Repository interface {
	Add(record number.NumberRequest) error
}

func New(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}
