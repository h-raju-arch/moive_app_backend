package movierepo

import (
	"context"
	"fmt"
)

func (r Movie_repo) FetchGenres(ctx context.Context, id string) ([]string, error) {

	var res []string
	query := `SELECT g.name FROM genres g JOIN movie_genres mg ON g.id = mg.genre_id  WHERE mg.movie_id = $1`

	rows, err := r.db.QueryContext(ctx, query, id)

	if err != nil {
		return nil, fmt.Errorf("Error Query FetchGenres: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("Erro Fetch Genres row scan: %w", err)
		}
		res = append(res, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error Fetch Gneres row: %w", err)
	}
	return res, nil
}
