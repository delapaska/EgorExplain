package number

import "github.com/delapaska/EgorExplain/internal/models"

func (s *service) Add(number models.NumberRequest) error {
	return s.repo.Add(number)
}
