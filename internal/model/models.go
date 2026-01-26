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
	Name        string
	Known_for   string
	Credit_type string
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
