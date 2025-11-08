package number

import (
	"github.com/delapaska/EgorExplain/internal/models"
)

func (r *repository) Add(number models.NumberRequest) error {

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
