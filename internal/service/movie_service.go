package service

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/h-raju-arch/movie_app_backend/internal/model"
	movierepo "github.com/h-raju-arch/movie_app_backend/internal/repo/movie_repo"
)

type Movie_Service interface {
	GetMovieById(ctx context.Context, id, lang string, appendtoresponse []string) (model.MovieResponse, error)
	SearchMovie(ctx context.Context, searchQuery string, language string, includeAdult bool, primaryYear sql.NullInt64, region sql.NullString, page, pageSize int) (model.SearchResponse, error)
}

type movie_service struct {
	repo movierepo.Movie_repo
}

func New_Movie_Service(r movierepo.Movie_repo) *movie_service {
	return &movie_service{repo: r}
}

func (r movie_service) GetMovieById(ctx context.Context, id, lang string, appendtoresponse []string) (model.MovieResponse, error) {
	movie, err := r.repo.GetMovieBasebyId(ctx, id, lang)

	if err != nil {
		return model.MovieResponse{}, fmt.Errorf("service: Get base movie: %w", err)
	}

	var genres []string
	var companies []string
	var credits []model.Credits_Response

	for _, append := range appendtoresponse {
		switch append {
		case "genre":
			genres, err = r.repo.FetchGenres(ctx, id)
			if err != nil {
				return model.MovieResponse{}, fmt.Errorf("fetch Genres: %w", err)
			}
		case "companies":
			companies, err = r.repo.FetchCompanies(ctx, id)
			if err != nil {
				return model.MovieResponse{}, fmt.Errorf("companies Genres: %w", err)
			}

		case "credits":
			credits, err = r.repo.FetchCredits(ctx, id)
			if err != nil {
				return model.MovieResponse{}, fmt.Errorf("fetch credits: %w", err)
			}
		}
	}

	res := model.MovieResponse{
		ID:           movie.ID,
		Title:        movie.Title,
		Overview:     movie.Overview,
		ReleaseDate:  movie.ReleaseDate,
		VoteAverage:  movie.VoteAverage,
		VoteCount:    movie.VoteCount,
		PosterPath:   movie.PosterPath,
		BackdropPath: movie.BackdropPath,
		Budget:       movie.Budget,
		Revenue:      movie.Revenue,
		Homepage:     movie.Homepage,
	}

	if contains(appendtoresponse, "genre") {
		res.Genres = genres
	}
	if contains(appendtoresponse, "companies") {
		res.ProductionCompanies = companies
	}
	if contains(appendtoresponse, "credits") {
		res.Credits = credits
	}
	return res, nil
}

func contains(list []string, s string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

func (r movie_service) SearchMovie(ctx context.Context, searchQuery string, language string, includeAdult bool, primaryYear sql.NullInt64, region sql.NullString, page, pageSize int) (model.SearchResponse, error) {

	total, items, err := r.repo.SearchMovie(ctx, searchQuery, includeAdult, language, primaryYear, region, page, pageSize)

	if err != nil {
		return model.SearchResponse{}, fmt.Errorf("service: SearchMovie : %w", err)
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}

	return model.SearchResponse{
		Page:         page,
		TotalResults: total,
		TotalPages:   totalPages,
		Results:      items,
	}, nil
}
