package httptransport

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/h-raju-arch/movie_app_backend/internal/model"
	"github.com/h-raju-arch/movie_app_backend/internal/service"
)

type Movie_handler struct {
	svc service.Movie_Service
}

func New_Movie_Handler(svc service.Movie_Service) *Movie_handler {
	return &Movie_handler{svc: svc}
}

func (h Movie_handler) GetMovies(c *gin.Context) {

	ctx := c.Request.Context()
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie id needed"})
		return
	}

	_, err := uuid.FromString(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	lang := c.DefaultQuery("lang", "en")
	appendtoresponse := c.Query("append_to_response")
	var appends []string
	if appendtoresponse != "" {
		for _, s := range strings.Split(appendtoresponse, ",") {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			appends = append(appends, s)
		}
	}

	res, err := h.svc.GetMovieById(ctx, id, lang, appends)

	if err != nil {
		fmt.Println("GetMovies error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// /------------------------------------------------///
func (h *Movie_handler) SearchMovieHandler(c *gin.Context) {

	ctx := c.Request.Context()
	q := strings.TrimSpace(c.Query("query"))

	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter reqired"})
		return
	}

	includeAdultStr := c.DefaultQuery("include_adult", "false")
	includeAdult := includeAdultStr == "true" || includeAdultStr == "1"

	language := c.DefaultQuery("language", "en-US")
	yearStr := c.DefaultQuery("primary_release_year", c.Query("year"))

	var primaryYear sql.NullInt64
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			primaryYear.Int64 = int64(y)
			primaryYear.Valid = true
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}
	}

	region := c.Query("region")
	var regionNull sql.NullString
	if region != "" {
		regionNull.String = region
		regionNull.Valid = true
	}

	page := 1

	if p := c.DefaultQuery("page", "1"); p != "" {
		if pi, err := strconv.Atoi(p); err == nil && pi >= 1 {
			page = pi
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page (must be >= 1)"})
			return
		}
	}

	pageSize := 20
	if ps := c.DefaultQuery("page_size", "20"); ps != "" {
		if psi, err := strconv.Atoi(ps); err == nil && psi > 0 && psi <= 100 {
			pageSize = psi
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page_size (1-100)"})
			return
		}
	}

	resp, err := h.svc.SearchMovie(ctx, q, language, includeAdult, primaryYear, regionNull, page, pageSize)
	if err != nil {
		fmt.Println("Error Search Movie handler:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

/// ----------------Dicover Handler-----------------------

func (h *Movie_handler) DiscoverMovieHandler(c *gin.Context) {

	ctx := c.Request.Context()
	var params model.DiscoverMoviesParams

	include_adult := c.DefaultQuery("include_adult", "false") == "true"
	params.IncludeAdult = include_adult

	language := c.DefaultQuery("language", "en")
	params.Language = language

	sort_by := c.DefaultQuery("sort_by", "popularity.DESC")
	params.SortBy = sort_by

	with_genres := c.Query("with_genres")
	if with_genres != "" {
		if strings.Contains(with_genres, ",") {
			params.WithGenres = strings.Split(with_genres, ",")
			params.WithGenresAND = true
		} else {
			params.WithGenres = strings.Split(with_genres, "|")
			params.WithGenresAND = false
		}

		for i := range params.WithGenres {
			params.WithGenres[i] = strings.TrimSpace(params.WithGenres[i])
		}
	}

	if v := c.Query("releaseGTE"); v != "" {
		params.ReleaseDateGTE = &v
	}
	if v := c.Query("releaseLTE"); v != "" {
		params.ReleaseDateLTE = &v
	}

	if v := c.Query("VoteAvgGTE"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			params.VoteAvgGTE = &f
		}
	}
	if v := c.Query("VoteAvgLTE"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			params.VoteAvgLTE = &f
		}
	}

	params.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	params.PageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if params.PageSize > 100 {
		params.PageSize = 100
	}

	res, err := h.svc.Discover(ctx, params)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": err,
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
