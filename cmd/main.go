package main

import (
	"github.com/h-raju-arch/movie_app_backend/internal/db"
	movierepo "github.com/h-raju-arch/movie_app_backend/internal/repo/movie_repo"
	"github.com/h-raju-arch/movie_app_backend/internal/service"
	httptransport "github.com/h-raju-arch/movie_app_backend/internal/transport/http"
)

func main() {
	database := db.Open()
	defer database.Close()
	repo := movierepo.New_Movie_Repo(database)
	svc := service.New_Movie_Service(*repo)
	router := httptransport.NewRouter(svc)

	router.Run(":3000")
}
