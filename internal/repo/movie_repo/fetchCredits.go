package movierepo

import (
	"context"
	"fmt"

	"github.com/h-raju-arch/movie_app_backend/internal/model"
)

func (r Movie_repo) FetchCredits(ctx context.Context, id string) ([]model.Credits_Response, error) {
	query := `SELECT p.name,p.known_for,c.credit_type
	          FROM people p JOIN credits c on p.id = c.person_id
			  WHERE c.movie_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return []model.Credits_Response{}, fmt.Errorf("Query FetchCredits : %w", err)
	}
	defer rows.Close()

	var resp []model.Credits_Response
	for rows.Next() {
		var temp model.Credits_Response
		err := rows.Scan(&temp.Name, &temp.Known_for, &temp.Credit_type)
		if err != nil {
			return []model.Credits_Response{}, fmt.Errorf("Error Credit rows scan: %w", err)
		}
		resp = append(resp, temp)
	}

	if err := rows.Err(); err != nil {
		return []model.Credits_Response{}, fmt.Errorf("Error Fetch Credits row: %w", err)
	}
	return resp, nil
}
