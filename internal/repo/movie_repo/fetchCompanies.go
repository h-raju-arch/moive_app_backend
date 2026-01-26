package movierepo

import (
	"context"
	"fmt"
)

func (r Movie_repo) FetchCompanies(ctx context.Context, id string) ([]string, error) {

	query := `SELECT c.name from companies c join movie_companies mc on
	         c.id = mc.company_id WHERE mc.movie_id = $1`

	rows, err := r.db.QueryContext(ctx, query, id)

	if err != nil {
		return nil, fmt.Errorf("Error Query Fetch Companies: %w", err)
	}

	defer rows.Close()
	var resp []string

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("Error Fetch Comapanies scan: %w", err)
		}
		resp = append(resp, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error fetch companies rows: %w", err)
	}
	return resp, nil
}
