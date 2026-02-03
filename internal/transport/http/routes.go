package httptransport

import (
	"github.com/gin-gonic/gin"
	"github.com/h-raju-arch/movie_app_backend/internal/service"
)

func NewRouter(movie_svc service.Movie_Service) *gin.Engine {
	router := gin.Default()
	h := New_Movie_Handler(movie_svc)

	api := router.Group("/api")
	{
		api.GET("/movie/", h.GetMovies)
		api.GET("/movies/search", h.SearchMovieHandler)
		api.GET("/movies/discover", h.DiscoverMovieHandler)
	}
	return router
}
