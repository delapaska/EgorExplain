package number

import (
	"database/sql"
	"github.com/delapaska/EgorExplain/internal/models"
)

type repository struct {
	db *sql.DB
}

type Repository interface {
	Add(record models.NumberRequest) error
}

func New(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}
