package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/h-raju-arch/movie_app_backend/internal/model"
)

// MockMovieRepo is a manual mock implementation of MovieRepository
type MockMovieRepo struct {
	GetMovieBasebyIdFunc func(ctx context.Context, id, lang string) (model.MovieResponse, error)
	FetchGenresFunc      func(ctx context.Context, id string) ([]string, error)
	FetchCompaniesFunc   func(ctx context.Context, id string) ([]string, error)
	FetchCreditsFunc     func(ctx context.Context, id string) ([]model.Credits_Response, error)
	SearchMovieFunc      func(ctx context.Context, query string, includeAdult bool, lang string, year sql.NullInt64, region sql.NullString, page, pageSize int) (int, []model.MovieSearchItem, error)
	DiscoverMoviesFunc   func(ctx context.Context, params model.DiscoverMoviesParams) ([]model.DiscoverItem, int, error)
}

func (m *MockMovieRepo) GetMovieBasebyId(ctx context.Context, id, lang string) (model.MovieResponse, error) {
	if m.GetMovieBasebyIdFunc != nil {
		return m.GetMovieBasebyIdFunc(ctx, id, lang)
	}
	return model.MovieResponse{}, nil
}

func (m *MockMovieRepo) FetchGenres(ctx context.Context, id string) ([]string, error) {
	if m.FetchGenresFunc != nil {
		return m.FetchGenresFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockMovieRepo) FetchCompanies(ctx context.Context, id string) ([]string, error) {
	if m.FetchCompaniesFunc != nil {
		return m.FetchCompaniesFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockMovieRepo) FetchCredits(ctx context.Context, id string) ([]model.Credits_Response, error) {
	if m.FetchCreditsFunc != nil {
		return m.FetchCreditsFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockMovieRepo) SearchMovie(ctx context.Context, query string, includeAdult bool, lang string, year sql.NullInt64, region sql.NullString, page, pageSize int) (int, []model.MovieSearchItem, error) {
	if m.SearchMovieFunc != nil {
		return m.SearchMovieFunc(ctx, query, includeAdult, lang, year, region, page, pageSize)
	}
	return 0, nil, nil
}

func (m *MockMovieRepo) DiscoverMovies(ctx context.Context, params model.DiscoverMoviesParams) ([]model.DiscoverItem, int, error) {
	if m.DiscoverMoviesFunc != nil {
		return m.DiscoverMoviesFunc(ctx, params)
	}
	return nil, 0, nil
}

// Test GetMovieById
func TestGetMovieById_Success(t *testing.T) {
	movieID := uuid.Must(uuid.NewV4())
	expectedMovie := model.MovieResponse{
		ID:    movieID,
		Title: "Test Movie",
	}

	mockRepo := &MockMovieRepo{
		GetMovieBasebyIdFunc: func(ctx context.Context, id, lang string) (model.MovieResponse, error) {
			if id != movieID.String() {
				t.Errorf("expected id %s, got %s", movieID.String(), id)
			}
			if lang != "en" {
				t.Errorf("expected lang 'en', got %s", lang)
			}
			return expectedMovie, nil
		},
	}

	svc := New_Movie_Service(mockRepo)
	result, err := svc.GetMovieById(context.Background(), movieID.String(), "en", []string{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Title != expectedMovie.Title {
		t.Errorf("expected title %s, got %s", expectedMovie.Title, result.Title)
	}
	if result.ID != expectedMovie.ID {
		t.Errorf("expected id %s, got %s", expectedMovie.ID, result.ID)
	}
}

func TestGetMovieById_WithGenres(t *testing.T) {
	movieID := uuid.Must(uuid.NewV4())
	expectedMovie := model.MovieResponse{
		ID:    movieID,
		Title: "Test Movie",
	}
	expectedGenres := []string{"Action", "Comedy"}

	mockRepo := &MockMovieRepo{
		GetMovieBasebyIdFunc: func(ctx context.Context, id, lang string) (model.MovieResponse, error) {
			return expectedMovie, nil
		},
		FetchGenresFunc: func(ctx context.Context, id string) ([]string, error) {
			return expectedGenres, nil
		},
	}

	svc := New_Movie_Service(mockRepo)
	result, err := svc.GetMovieById(context.Background(), movieID.String(), "en", []string{"genres"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Genres) != len(expectedGenres) {
		t.Errorf("expected %d genres, got %d", len(expectedGenres), len(result.Genres))
	}
	for i, g := range expectedGenres {
		if result.Genres[i] != g {
			t.Errorf("expected genre %s at index %d, got %s", g, i, result.Genres[i])
		}
	}
}

func TestGetMovieById_WithCompanies(t *testing.T) {
	movieID := uuid.Must(uuid.NewV4())
	expectedMovie := model.MovieResponse{
		ID:    movieID,
		Title: "Test Movie",
	}
	expectedCompanies := []string{"Warner Bros", "Universal"}

	mockRepo := &MockMovieRepo{
		GetMovieBasebyIdFunc: func(ctx context.Context, id, lang string) (model.MovieResponse, error) {
			return expectedMovie, nil
		},
		FetchCompaniesFunc: func(ctx context.Context, id string) ([]string, error) {
			return expectedCompanies, nil
		},
	}

	svc := New_Movie_Service(mockRepo)
	result, err := svc.GetMovieById(context.Background(), movieID.String(), "en", []string{"companies"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.ProductionCompanies) != len(expectedCompanies) {
		t.Errorf("expected %d companies, got %d", len(expectedCompanies), len(result.ProductionCompanies))
	}
}

func TestGetMovieById_WithCredits(t *testing.T) {
	movieID := uuid.Must(uuid.NewV4())
	expectedMovie := model.MovieResponse{
		ID:    movieID,
		Title: "Test Movie",
	}
	expectedCredits := []model.Credits_Response{
		{Name: "Actor 1", Known_for: "Acting", Credit_type: "cast"},
		{Name: "Director 1", Known_for: "Directing", Credit_type: "crew"},
	}

	mockRepo := &MockMovieRepo{
		GetMovieBasebyIdFunc: func(ctx context.Context, id, lang string) (model.MovieResponse, error) {
			return expectedMovie, nil
		},
		FetchCreditsFunc: func(ctx context.Context, id string) ([]model.Credits_Response, error) {
			return expectedCredits, nil
		},
	}

	svc := New_Movie_Service(mockRepo)
	result, err := svc.GetMovieById(context.Background(), movieID.String(), "en", []string{"credits"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Credits) != len(expectedCredits) {
		t.Errorf("expected %d credits, got %d", len(expectedCredits), len(result.Credits))
	}
}

func TestGetMovieById_RepoError(t *testing.T) {
	mockRepo := &MockMovieRepo{
		GetMovieBasebyIdFunc: func(ctx context.Context, id, lang string) (model.MovieResponse, error) {
			return model.MovieResponse{}, errors.New("movie not found")
		},
	}

	svc := New_Movie_Service(mockRepo)
	_, err := svc.GetMovieById(context.Background(), "invalid-id", "en", []string{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "service: Get base movie: movie not found" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGetMovieById_GenresError(t *testing.T) {
	movieID := uuid.Must(uuid.NewV4())
	expectedMovie := model.MovieResponse{
		ID:    movieID,
		Title: "Test Movie",
	}

	mockRepo := &MockMovieRepo{
		GetMovieBasebyIdFunc: func(ctx context.Context, id, lang string) (model.MovieResponse, error) {
			return expectedMovie, nil
		},
		FetchGenresFunc: func(ctx context.Context, id string) ([]string, error) {
			return nil, errors.New("failed to fetch genres")
		},
	}

	svc := New_Movie_Service(mockRepo)
	_, err := svc.GetMovieById(context.Background(), movieID.String(), "en", []string{"genres"})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// Test SearchMovie
func TestSearchMovie_Success(t *testing.T) {
	items := []model.MovieSearchItem{
		{Title: "Movie 1"},
		{Title: "Movie 2"},
	}

	mockRepo := &MockMovieRepo{
		SearchMovieFunc: func(ctx context.Context, query string, includeAdult bool, lang string, year sql.NullInt64, region sql.NullString, page, pageSize int) (int, []model.MovieSearchItem, error) {
			if query != "test" {
				t.Errorf("expected query 'test', got %s", query)
			}
			return 2, items, nil
		},
	}

	svc := New_Movie_Service(mockRepo)
	result, err := svc.SearchMovie(context.Background(), "test", "en", false, sql.NullInt64{}, sql.NullString{}, 1, 20)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.TotalResults != 2 {
		t.Errorf("expected 2 total results, got %d", result.TotalResults)
	}
	if len(result.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(result.Results))
	}
}

func TestSearchMovie_Pagination(t *testing.T) {
	items := []model.MovieSearchItem{
		{Title: "Movie 1"},
	}

	mockRepo := &MockMovieRepo{
		SearchMovieFunc: func(ctx context.Context, query string, includeAdult bool, lang string, year sql.NullInt64, region sql.NullString, page, pageSize int) (int, []model.MovieSearchItem, error) {
			return 45, items, nil
		},
	}

	svc := New_Movie_Service(mockRepo)
	result, err := svc.SearchMovie(context.Background(), "test", "en", false, sql.NullInt64{}, sql.NullString{}, 1, 20)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.TotalResults != 45 {
		t.Errorf("expected 45 total results, got %d", result.TotalResults)
	}
	// 45 / 20 = 2.25, ceil = 3 pages
	if result.TotalPages != 3 {
		t.Errorf("expected 3 total pages, got %d", result.TotalPages)
	}
	if result.Page != 1 {
		t.Errorf("expected page 1, got %d", result.Page)
	}
}

func TestSearchMovie_EmptyResult(t *testing.T) {
	mockRepo := &MockMovieRepo{
		SearchMovieFunc: func(ctx context.Context, query string, includeAdult bool, lang string, year sql.NullInt64, region sql.NullString, page, pageSize int) (int, []model.MovieSearchItem, error) {
			return 0, []model.MovieSearchItem{}, nil
		},
	}

	svc := New_Movie_Service(mockRepo)
	result, err := svc.SearchMovie(context.Background(), "nonexistent", "en", false, sql.NullInt64{}, sql.NullString{}, 1, 20)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.TotalResults != 0 {
		t.Errorf("expected 0 total results, got %d", result.TotalResults)
	}
	if result.TotalPages != 0 {
		t.Errorf("expected 0 total pages, got %d", result.TotalPages)
	}
}

func TestSearchMovie_Error(t *testing.T) {
	mockRepo := &MockMovieRepo{
		SearchMovieFunc: func(ctx context.Context, query string, includeAdult bool, lang string, year sql.NullInt64, region sql.NullString, page, pageSize int) (int, []model.MovieSearchItem, error) {
			return 0, nil, errors.New("database error")
		},
	}

	svc := New_Movie_Service(mockRepo)
	_, err := svc.SearchMovie(context.Background(), "test", "en", false, sql.NullInt64{}, sql.NullString{}, 1, 20)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// Test Discover
func TestDiscover_Success(t *testing.T) {
	items := []model.DiscoverItem{
		{Title: "Discover Movie 1"},
		{Title: "Discover Movie 2"},
	}

	mockRepo := &MockMovieRepo{
		DiscoverMoviesFunc: func(ctx context.Context, params model.DiscoverMoviesParams) ([]model.DiscoverItem, int, error) {
			return items, 2, nil
		},
	}

	svc := New_Movie_Service(mockRepo)
	params := model.DiscoverMoviesParams{
		Language: "en",
		Page:     1,
		PageSize: 20,
	}
	result, err := svc.Discover(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(result.Results))
	}
	if result.Page != 1 {
		t.Errorf("expected page 1, got %d", result.Page)
	}
}

func TestDiscover_WithGenres(t *testing.T) {
	items := []model.DiscoverItem{
		{Title: "Action Movie"},
	}

	mockRepo := &MockMovieRepo{
		DiscoverMoviesFunc: func(ctx context.Context, params model.DiscoverMoviesParams) ([]model.DiscoverItem, int, error) {
			if len(params.WithGenres) == 0 {
				t.Error("expected genres to be passed")
			}
			return items, 1, nil
		},
	}

	svc := New_Movie_Service(mockRepo)
	params := model.DiscoverMoviesParams{
		Language:   "en",
		WithGenres: []string{"action-genre-id"},
		Page:       1,
		PageSize:   20,
	}
	result, err := svc.Discover(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(result.Results))
	}
}

func TestDiscover_Error(t *testing.T) {
	mockRepo := &MockMovieRepo{
		DiscoverMoviesFunc: func(ctx context.Context, params model.DiscoverMoviesParams) ([]model.DiscoverItem, int, error) {
			return nil, 0, errors.New("database error")
		},
	}

	svc := New_Movie_Service(mockRepo)
	params := model.DiscoverMoviesParams{
		Language: "en",
		Page:     1,
		PageSize: 20,
	}
	_, err := svc.Discover(context.Background(), params)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// Test helper function contains
func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		list     []string
		s        string
		expected bool
	}{
		{"found", []string{"a", "b", "c"}, "b", true},
		{"not found", []string{"a", "b", "c"}, "d", false},
		{"empty list", []string{}, "a", false},
		{"empty string found", []string{"a", "", "c"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.list, tt.s)
			if result != tt.expected {
				t.Errorf("contains(%v, %s) = %v, expected %v", tt.list, tt.s, result, tt.expected)
			}
		})
	}
}
