package httptransport

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/h-raju-arch/movie_app_backend/internal/model"
)

// MockMovieService is a manual mock implementation of Movie_Service
type MockMovieService struct {
	GetMovieByIdFunc func(ctx context.Context, id, lang string, appendtoresponse []string) (model.MovieResponse, error)
	SearchMovieFunc  func(ctx context.Context, searchQuery string, language string, includeAdult bool, primaryYear sql.NullInt64, region sql.NullString, page, pageSize int) (model.SearchResponse, error)
	DiscoverFunc     func(ctx context.Context, params model.DiscoverMoviesParams) (model.DiscoverMoviesResponse, error)
}

func (m *MockMovieService) GetMovieById(ctx context.Context, id, lang string, appendtoresponse []string) (model.MovieResponse, error) {
	if m.GetMovieByIdFunc != nil {
		return m.GetMovieByIdFunc(ctx, id, lang, appendtoresponse)
	}
	return model.MovieResponse{}, nil
}

func (m *MockMovieService) SearchMovie(ctx context.Context, searchQuery string, language string, includeAdult bool, primaryYear sql.NullInt64, region sql.NullString, page, pageSize int) (model.SearchResponse, error) {
	if m.SearchMovieFunc != nil {
		return m.SearchMovieFunc(ctx, searchQuery, language, includeAdult, primaryYear, region, page, pageSize)
	}
	return model.SearchResponse{}, nil
}

func (m *MockMovieService) Discover(ctx context.Context, params model.DiscoverMoviesParams) (model.DiscoverMoviesResponse, error) {
	if m.DiscoverFunc != nil {
		return m.DiscoverFunc(ctx, params)
	}
	return model.DiscoverMoviesResponse{}, nil
}

func setupTestRouter(handler *Movie_handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/movies", handler.GetMovies)
	r.GET("/search", handler.SearchMovieHandler)
	r.GET("/discover", handler.DiscoverMovieHandler)
	return r
}

// GetMovies handler tests
func TestGetMovies_MissingID(t *testing.T) {
	mockSvc := &MockMovieService{}
	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/movies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["error"] != "Movie id needed" {
		t.Errorf("expected error 'Movie id needed', got %s", response["error"])
	}
}

func TestGetMovies_InvalidID(t *testing.T) {
	mockSvc := &MockMovieService{}
	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/movies?id=not-a-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetMovies_Success(t *testing.T) {
	movieID := uuid.Must(uuid.NewV4())
	expectedMovie := model.MovieResponse{
		ID:    movieID,
		Title: "Inception",
	}

	mockSvc := &MockMovieService{
		GetMovieByIdFunc: func(ctx context.Context, id, lang string, appendtoresponse []string) (model.MovieResponse, error) {
			if id != movieID.String() {
				t.Errorf("expected id %s, got %s", movieID.String(), id)
			}
			return expectedMovie, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/movies?id="+movieID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.MovieResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Title != "Inception" {
		t.Errorf("expected title 'Inception', got %s", response.Title)
	}
}

func TestGetMovies_WithLanguage(t *testing.T) {
	movieID := uuid.Must(uuid.NewV4())
	expectedMovie := model.MovieResponse{
		ID:    movieID,
		Title: "インセプション", // Japanese title
	}

	mockSvc := &MockMovieService{
		GetMovieByIdFunc: func(ctx context.Context, id, lang string, appendtoresponse []string) (model.MovieResponse, error) {
			if lang != "ja" {
				t.Errorf("expected lang 'ja', got %s", lang)
			}
			return expectedMovie, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/movies?id="+movieID.String()+"&lang=ja", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetMovies_WithAppendToResponse(t *testing.T) {
	movieID := uuid.Must(uuid.NewV4())
	expectedMovie := model.MovieResponse{
		ID:     movieID,
		Title:  "Inception",
		Genres: []string{"Sci-Fi", "Action"},
	}

	mockSvc := &MockMovieService{
		GetMovieByIdFunc: func(ctx context.Context, id, lang string, appendtoresponse []string) (model.MovieResponse, error) {
			if len(appendtoresponse) != 2 {
				t.Errorf("expected 2 append items, got %d", len(appendtoresponse))
			}
			return expectedMovie, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/movies?id="+movieID.String()+"&append_to_response=genres,credits", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetMovies_ServiceError(t *testing.T) {
	movieID := uuid.Must(uuid.NewV4())

	mockSvc := &MockMovieService{
		GetMovieByIdFunc: func(ctx context.Context, id, lang string, appendtoresponse []string) (model.MovieResponse, error) {
			return model.MovieResponse{}, errors.New("movie not found")
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/movies?id="+movieID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// SearchMovieHandler tests
func TestSearchMovie_MissingQuery(t *testing.T) {
	mockSvc := &MockMovieService{}
	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["error"] != "query parameter reqired" {
		t.Errorf("expected error 'query parameter reqired', got %s", response["error"])
	}
}

func TestSearchMovie_EmptyQuery(t *testing.T) {
	mockSvc := &MockMovieService{}
	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/search?query=", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSearchMovie_Success(t *testing.T) {
	movieID := uuid.Must(uuid.NewV4())
	expectedResponse := model.SearchResponse{
		Page:         1,
		TotalResults: 10,
		TotalPages:   1,
		Results: []model.MovieSearchItem{
			{ID: movieID, Title: "Inception"},
		},
	}

	mockSvc := &MockMovieService{
		SearchMovieFunc: func(ctx context.Context, searchQuery string, language string, includeAdult bool, primaryYear sql.NullInt64, region sql.NullString, page, pageSize int) (model.SearchResponse, error) {
			if searchQuery != "inception" {
				t.Errorf("expected query 'inception', got %s", searchQuery)
			}
			return expectedResponse, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/search?query=inception", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.SearchResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.TotalResults != 10 {
		t.Errorf("expected 10 total results, got %d", response.TotalResults)
	}
}

func TestSearchMovie_WithPagination(t *testing.T) {
	mockSvc := &MockMovieService{
		SearchMovieFunc: func(ctx context.Context, searchQuery string, language string, includeAdult bool, primaryYear sql.NullInt64, region sql.NullString, page, pageSize int) (model.SearchResponse, error) {
			if page != 2 {
				t.Errorf("expected page 2, got %d", page)
			}
			if pageSize != 50 {
				t.Errorf("expected pageSize 50, got %d", pageSize)
			}
			return model.SearchResponse{Page: 2}, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/search?query=test&page=2&page_size=50", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSearchMovie_InvalidYear(t *testing.T) {
	mockSvc := &MockMovieService{}
	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/search?query=test&year=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSearchMovie_InvalidPageSize(t *testing.T) {
	mockSvc := &MockMovieService{}
	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/search?query=test&page_size=200", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSearchMovie_WithAdultContent(t *testing.T) {
	mockSvc := &MockMovieService{
		SearchMovieFunc: func(ctx context.Context, searchQuery string, language string, includeAdult bool, primaryYear sql.NullInt64, region sql.NullString, page, pageSize int) (model.SearchResponse, error) {
			if !includeAdult {
				t.Error("expected includeAdult to be true")
			}
			return model.SearchResponse{}, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/search?query=test&include_adult=true", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// DiscoverMovieHandler tests
func TestDiscoverMovie_Success(t *testing.T) {
	movieID := uuid.Must(uuid.NewV4())
	expectedResponse := model.DiscoverMoviesResponse{
		Page:     1,
		PageSize: 20,
		Results: []model.DiscoverItem{
			{ID: movieID, Title: "Popular Movie"},
		},
	}

	mockSvc := &MockMovieService{
		DiscoverFunc: func(ctx context.Context, params model.DiscoverMoviesParams) (model.DiscoverMoviesResponse, error) {
			return expectedResponse, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/discover", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.DiscoverMoviesResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(response.Results))
	}
}

func TestDiscoverMovie_WithGenres(t *testing.T) {
	mockSvc := &MockMovieService{
		DiscoverFunc: func(ctx context.Context, params model.DiscoverMoviesParams) (model.DiscoverMoviesResponse, error) {
			if len(params.WithGenres) == 0 {
				t.Error("expected genres to be passed")
			}
			return model.DiscoverMoviesResponse{}, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/discover?with_genres=action|comedy", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDiscoverMovie_WithGenresAND(t *testing.T) {
	mockSvc := &MockMovieService{
		DiscoverFunc: func(ctx context.Context, params model.DiscoverMoviesParams) (model.DiscoverMoviesResponse, error) {
			if !params.WithGenresAND {
				t.Error("expected WithGenresAND to be true when comma-separated")
			}
			return model.DiscoverMoviesResponse{}, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/discover?with_genres=action,comedy", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDiscoverMovie_WithIncludeAdult(t *testing.T) {
	mockSvc := &MockMovieService{
		DiscoverFunc: func(ctx context.Context, params model.DiscoverMoviesParams) (model.DiscoverMoviesResponse, error) {
			if !params.IncludeAdult {
				t.Error("expected IncludeAdult to be true")
			}
			return model.DiscoverMoviesResponse{}, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/discover?include_adult=true", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDiscoverMovie_WithFilters(t *testing.T) {
	mockSvc := &MockMovieService{
		DiscoverFunc: func(ctx context.Context, params model.DiscoverMoviesParams) (model.DiscoverMoviesResponse, error) {
			if params.ReleaseDateGTE == nil || *params.ReleaseDateGTE != "2020-01-01" {
				t.Error("expected release date filter")
			}
			if params.VoteAvgGTE == nil || *params.VoteAvgGTE != 7.5 {
				t.Error("expected vote average filter")
			}
			return model.DiscoverMoviesResponse{}, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/discover?releaseGTE=2020-01-01&VoteAvgGTE=7.5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDiscoverMovie_PageSizeCapped(t *testing.T) {
	mockSvc := &MockMovieService{
		DiscoverFunc: func(ctx context.Context, params model.DiscoverMoviesParams) (model.DiscoverMoviesResponse, error) {
			if params.PageSize != 100 {
				t.Errorf("expected pageSize to be capped at 100, got %d", params.PageSize)
			}
			return model.DiscoverMoviesResponse{}, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/discover?page_size=500", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDiscoverMovie_ServiceError(t *testing.T) {
	mockSvc := &MockMovieService{
		DiscoverFunc: func(ctx context.Context, params model.DiscoverMoviesParams) (model.DiscoverMoviesResponse, error) {
			return model.DiscoverMoviesResponse{}, errors.New("database error")
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/discover", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("expected status %d, got %d", http.StatusBadGateway, w.Code)
	}
}

func TestDiscoverMovie_SortBy(t *testing.T) {
	mockSvc := &MockMovieService{
		DiscoverFunc: func(ctx context.Context, params model.DiscoverMoviesParams) (model.DiscoverMoviesResponse, error) {
			if params.SortBy != "vote_average.DESC" {
				t.Errorf("expected sort_by 'vote_average.DESC', got %s", params.SortBy)
			}
			return model.DiscoverMoviesResponse{}, nil
		},
	}

	handler := New_Movie_Handler(mockSvc)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/discover?sort_by=vote_average.DESC", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
