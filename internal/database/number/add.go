package number

import (
	"github.com/delapaska/EgorExplain/internal/models/number"
)

func (r *repository) Add(number number.NumberRequest) error {

	query := `
	INSERT INTO numbers (number)
	VALUES ($1)
`
	_, err := r.db.Exec(query, number.Number)
	if err != nil {
		return err
	}
	return nil
}
