package movierepo

import (
	"context"
	"fmt"

	"github.com/h-raju-arch/movie_app_backend/internal/model"
)

func (r *Movie_repo) GetMovieBasebyId(ctx context.Context, id, lang string) (model.MovieResponse, error) {
	var res model.MovieResponse

	query := `SELECT 
	           m.id, COALESCE(mt.title,m.title) AS title, COALESCE(mt.overview,m.overview) AS overview,
			   to_char(m.release_date, 'YYYY-MM-DD') AS release_data,
			   ms.vote_average, ms.vote_count,
			   m.poster_path, m.backdrop_path,m.budget,m.revenue,m.homepage
			   FROM movies m
			   LEFT JOIN movie_translations mt ON mt.movie_id = m.id AND mt.language = $2
			   LEFT JOIN movie_stats ms ON ms.movie_id = m.id
			   WHERE m.id = $1;`

	err := r.db.QueryRowContext(ctx, query, id, lang).Scan(&res.ID,
		&res.Title,
		&res.Overview,
		&res.ReleaseDate,
		&res.VoteAverage,
		&res.VoteCount,
		&res.PosterPath,
		&res.BackdropPath,
		&res.Budget,
		&res.Revenue,
		&res.Homepage)
	if err != nil {
		return model.MovieResponse{}, fmt.Errorf("Query movie base: %w", err)
	}
	return res, nil
}
