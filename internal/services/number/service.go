package number

import (
	"github.com/delapaska/EgorExplain/internal/database/number"
	"github.com/delapaska/EgorExplain/internal/models"
)

type service struct {
	repo number.Repository
}

type Service interface {
	Add(number models.NumberRequest) error
}

func New(repo number.Repository) Service {
	return &service{
		repo: repo,
	}
}
