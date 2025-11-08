package number

import (
	"github.com/delapaska/EgorExplain/internal/models/number"
)

func (s *service) Add(number number.NumberRequest) error {
	return s.repo.Add(number)
}
