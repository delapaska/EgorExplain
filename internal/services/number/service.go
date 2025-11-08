package number

import (
	"github.com/delapaska/EgorExplain/internal/database/number"
	number2 "github.com/delapaska/EgorExplain/internal/models/number"
)

type service struct {
	repo number.Repository
}

type Service interface {
	Add(number number2.NumberRequest) error
}

func New(repo number.Repository) Service {
	return &service{
		repo: repo,
	}
}
