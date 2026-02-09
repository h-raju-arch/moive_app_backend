package model

import "github.com/gofrs/uuid/v5"

type MovieResponse struct {
	ID                  uuid.UUID          `json:"id"`
	Title               string             `json:"title"`
	Overview            *string            `json:"overview,omitempty"`
	ReleaseDate         *string            `json:"release_date,omitempty"`
	VoteAverage         *float64           `json:"vote_average,omitempty"`
	VoteCount           *int               `json:"vote_count,omitempty"`
	PosterPath          *string            `json:"poster_path,omitempty"`
	BackdropPath        *string            `json:"backdrop_path,omitempty"`
	Budget              *int64             `json:"budget,omitempty"`
	Revenue             *int64             `json:"revenue,omitempty"`
	Genres              []string           `json:"genres,omitempty"`
	ProductionCompanies []string           `json:"production_companies,omitempty"`
	SpokenLanguages     []map[string]any   `json:"spoken_languages,omitempty"`
	Homepage            *string            `json:"homepage,omitempty"`
	Credits             []Credits_Response `json:"credits,omitempty"`
	Videos              []map[string]any   `json:"videos,omitempty"`
	Images              []map[string]any   `json:"images,omitempty"`
}

type Credits_Response struct {
	Name        string `json:"name"`
	Known_for   string `json:"known_for"`
	Credit_type string `json:"credit_type"`
}

type MovieSearchItem struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Overview    *string   `json:"overview,omitempty"`
	ReleaseDate *string   `json:"release_date,omitempty"`
	VoteAverage *float64  `json:"vote_average,omitempty"`
	Popularity  *float64  `json:"popularity,omitempty"`
}

type SearchResponse struct {
	Page         int               `json:"page"`
	TotalResults int               `json:"total_results"`
	TotalPages   int               `json:"total_pages"`
	Results      []MovieSearchItem `json:"results"`
}

type DiscoverMoviesParams struct {
	WithGenres     []string // UUID strings
	WithGenresAND  bool     // true if comma (AND), false if OR (pipe) — only used if WithGenres not empty
	IncludeAdult   bool
	Language       string
	ReleaseDateGTE *string // YYYY-MM-DD
	ReleaseDateLTE *string
	VoteAvgGTE     *float64
	VoteAvgLTE     *float64
	SortBy         string // e.g., "popularity.desc"
	Page           int
	PageSize       int
}

type DiscoverItem struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Overview     *string   `json:"overview,omitempty"`
	ReleaseDate  *string   `json:"release_date,omitempty"`
	VoteAverage  *float64  `json:"vote_average,omitempty"`
	VoteCount    *int      `json:"vote_count,omitempty"`
	PosterPath   *string   `json:"poster_path,omitempty"`
	BackdropPath *string   `json:"backdrop_path,omitempty"`
	Popularity   *float64  `json:"popularity,omitempty"`
	GenreIDs     []string  `json:"genre_ids,omitempty"`
}

type DiscoverMoviesResponse struct {
	Page         int            `json:"page"`
	PageSize     int            `json:"page_size"`
	TotalResults int            `json:"total_results"`
	TotalPages   int            `json:"total_pages"`
	Results      []DiscoverItem `json:"results"`
}
