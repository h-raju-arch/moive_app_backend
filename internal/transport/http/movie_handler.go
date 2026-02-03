package httptransport

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter required"})
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
		fmt.Println("GetMovies error: %w", err)
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page_size (1-100)"})
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
		fmt.Println("Error Search Movie handler: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, resp)
}
