package movierepo

import "database/sql"

type Movie_repo struct {
	db *sql.DB
}

func New_Movie_Repo(db *sql.DB) *Movie_repo {
	return &Movie_repo{db: db}
}
