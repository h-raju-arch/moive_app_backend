package movierepo

import (
	"context"
	"database/sql"

	"github.com/h-raju-arch/movie_app_backend/internal/model"
)

// MovieRepository defines the interface for movie data access operations.
// This interface allows for easy mocking in unit tests.
type MovieRepository interface {
	GetMovieBasebyId(ctx context.Context, id, lang string) (model.MovieResponse, error)
	FetchGenres(ctx context.Context, id string) ([]string, error)
	FetchCompanies(ctx context.Context, id string) ([]string, error)
	FetchCredits(ctx context.Context, id string) ([]model.Credits_Response, error)
	SearchMovie(ctx context.Context, query string, includeAdult bool, lang string, year sql.NullInt64, region sql.NullString, page, pageSize int) (int, []model.MovieSearchItem, error)
	DiscoverMovies(ctx context.Context, params model.DiscoverMoviesParams) ([]model.DiscoverItem, int, error)
}
