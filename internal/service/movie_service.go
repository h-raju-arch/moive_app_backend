package service

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sync"

	"github.com/h-raju-arch/movie_app_backend/internal/model"
	movierepo "github.com/h-raju-arch/movie_app_backend/internal/repo/movie_repo"
)

type Movie_Service interface {
	GetMovieById(ctx context.Context, id, lang string, appendtoresponse []string) (model.MovieResponse, error)
	SearchMovie(ctx context.Context, searchQuery string, language string, includeAdult bool, primaryYear sql.NullInt64, region sql.NullString, page, pageSize int) (model.SearchResponse, error)
	Discover(ctx context.Context, params model.DiscoverMoviesParams) (model.DiscoverMoviesResponse, error)
}

type movie_service struct {
	repo movierepo.Movie_repo
}

func New_Movie_Service(r movierepo.Movie_repo) *movie_service {
	return &movie_service{repo: r}
}

//getMoviebyId

func (r movie_service) GetMovieById(ctx context.Context, id, lang string, appendtoresponse []string) (model.MovieResponse, error) {
	movie, err := r.repo.GetMovieBasebyId(ctx, id, lang)

	if err != nil {
		return model.MovieResponse{}, fmt.Errorf("service: Get base movie: %w", err)
	}

	var genres []string
	var companies []string
	var credits []model.Credits_Response
	var wg sync.WaitGroup

	type result struct {
		typ       string
		genres    []string
		companies []string
		credits   []model.Credits_Response
		err       error
	}

	resultCh := make(chan result, len(appendtoresponse))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, append := range appendtoresponse {
		typ := append
		wg.Add(1)

		go func(typ string) {
			defer wg.Done()
			switch typ {
			case "genres":
				genres, err = r.repo.FetchGenres(ctx, id)
				resultCh <- result{genres: genres, typ: typ, err: err}
				if err != nil {
					cancel()
				}

			case "companies":
				companies, err = r.repo.FetchCompanies(ctx, id)
				resultCh <- result{companies: companies, typ: typ, err: err}

				if err != nil {
					cancel()
				}

			case "credits":
				credits, err = r.repo.FetchCredits(ctx, id)
				resultCh <- result{credits: credits, typ: typ, err: err}

				if err != nil {
					cancel()
				}
			default:
				resultCh <- result{typ: typ}
			}
		}(typ)

	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for itr := range resultCh {
		if itr.err != nil {
			return movie, itr.err
		}
		if itr.typ == "genres" {
			genres = append(genres, itr.genres...)
		}
		if itr.typ == "companies" {
			companies = append(companies, itr.companies...)
		}
		if itr.typ == "credits" {
			credits = append(credits, itr.credits...)
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

	if contains(appendtoresponse, "genres") {
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

//Search

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

// discover
func (r movie_service) Discover(ctx context.Context, params model.DiscoverMoviesParams) (model.DiscoverMoviesResponse, error) {

	resp, totalCount, err := r.repo.DiscoverMovies(ctx, params)

	if err != nil {
		fmt.Println("service", err)
		return model.DiscoverMoviesResponse{}, fmt.Errorf("service: DisoverMovie: %w", err)
	}

	result := model.DiscoverMoviesResponse{
		Results:  resp,
		Page:     params.Page,
		PageSize: params.PageSize,
	}

	if totalCount > 0 {
		result.TotalPages = int(math.Ceil(float64(result.TotalResults)) / math.Ceil(float64(result.PageSize)))
	}
	return result, nil
}
