package movierepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/h-raju-arch/movie_app_backend/internal/model"
)

func (r Movie_repo) SearchMovie(ctx context.Context, queryStr string, adult bool, lang string, primaryYear sql.NullInt64, region sql.NullString, page, pageSize int) (int, []model.MovieSearchItem, error) {
	fmt.Println("inside repo ", primaryYear)

	if page < 1 {
		page = 1
	}

	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var year interface{}

	if primaryYear.Valid {
		year = primaryYear.Int64
	} else {
		year = nil
	}

	fmt.Println("interface ", year)

	var regionParam interface{}
	if region.Valid {
		regionParam = region.String
	} else {
		regionParam = nil
	}

	query := `
SELECT
  m.id,
  COALESCE(mt.title, m.title) AS title,
  COALESCE(mt.overview, m.overview) AS overview,
  to_char(m.release_date, 'YYYY-MM-DD') AS release_date,
  ms.vote_average,
  ms.popularity,
  COUNT(*) OVER() AS total_count
FROM movies m
LEFT JOIN movie_translations mt ON mt.movie_id = m.id AND mt.language = $2
LEFT JOIN movie_stats ms ON ms.movie_id = m.id
WHERE (
    COALESCE(mt.title, m.title) ILIKE '%' || $1 || '%'
    OR COALESCE(mt.overview, m.overview) ILIKE '%' || $1 || '%'
)
AND ($3 OR m.adult = false)
AND ($4::int IS NULL OR EXTRACT(YEAR FROM m.release_date)::int = $4::int)
AND (
  $5::text IS NULL OR EXISTS (
    SELECT 1
    FROM movie_companies mc
    JOIN companies c ON c.id = mc.company_id
    WHERE mc.movie_id = m.id
      AND c.origin_country = $5::text
  )
)
ORDER BY ms.popularity DESC NULLS LAST, m.created_at DESC
LIMIT $6 OFFSET $7;
`

	rows, err := r.db.QueryContext(ctx, query,
		queryStr,    // $1
		lang,        // $2
		adult,       // $3
		year,        // $4 (sql.NullString)
		regionParam, // $5 (sql.NullString)
		pageSize,    // $6
		offset,      // $7
	)
	if err != nil {
		return 0, []model.MovieSearchItem{}, fmt.Errorf("Error Query search Movie: %w", err)
	}
	defer rows.Close()

	var result []model.MovieSearchItem
	totalCount := 0

	for rows.Next() {
		var (
			id          uuid.UUID
			title       string
			overview    sql.NullString
			releaseDate sql.NullString
			voteAverage sql.NullFloat64
			popularity  sql.NullFloat64
			total       sql.NullInt64
		)

		if err := rows.Scan(&id, &title, &overview, &releaseDate, &voteAverage, &popularity, &total); err != nil {

			return 0, []model.MovieSearchItem{}, fmt.Errorf("Error Search movie row scan: %w", err)
		}

		item := model.MovieSearchItem{
			ID:    id,
			Title: title,
		}

		if overview.Valid {
			item.Overview = &overview.String
		}
		if releaseDate.Valid {
			item.ReleaseDate = &releaseDate.String
		}
		if voteAverage.Valid {
			v := voteAverage.Float64
			item.VoteAverage = &v
		}
		if popularity.Valid {
			p := popularity.Float64
			item.Popularity = &p
		}

		result = append(result, item)

		if total.Valid {
			totalCount = int(total.Int64)
		}
	}

	if err := rows.Err(); err != nil {

		return 0, []model.MovieSearchItem{}, fmt.Errorf("Error Search movie rows: %w", err)
	}

	return totalCount, result, nil
}
