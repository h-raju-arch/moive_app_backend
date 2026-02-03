package movierepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/h-raju-arch/movie_app_backend/internal/model"
	"github.com/lib/pq"
)

func (r Movie_repo) DiscoverMovies(ctx context.Context, p model.DiscoverMoviesParams) ([]model.DiscoverItem, int, error) {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize <= 0 || p.PageSize > 100 {
		p.PageSize = 20
	}
	offset := (p.Page - 1) * p.PageSize

	// --------- SELECT (fixed ARRAY[]::text[] and corrected correlated subquery mg.movie_id = m.id)
	BaseQuery := `
    SELECT 
      m.id,
      COALESCE(mt.title,m.title) AS title,
      COALESCE(mt.overview,m.overview) AS overview,
      to_char(m.release_date, 'YYYY-MM-DD') AS release_date,
      ms.vote_average, ms.vote_count,
      m.poster_path, m.backdrop_path, ms.popularity,
      (SELECT COALESCE(array_agg(g.id::text), ARRAY[]::text[]) 
         FROM movie_genres mg JOIN genres g ON mg.genre_id = g.id 
         WHERE mg.movie_id = m.id
      ) AS genre_ids,
      COUNT(*) OVER() AS total_count
  `

	fromWhere := `FROM movies m 
    LEFT JOIN movie_translations mt ON mt.movie_id = m.id AND mt.language = $1 
    LEFT JOIN movie_stats ms ON ms.movie_id = m.id`

	// --------- where builder (fix "1==1" -> "1=1")
	where := []string{"1=1"}
	args := []any{p.Language}
	argsPos := 2

	addArg := func(val any) string {
		placeholder := fmt.Sprintf("$%d", argsPos)
		args = append(args, val)
		argsPos++
		return placeholder
	}

	// adult filter
	if !p.IncludeAdult {
		where = append(where, "m.adult = false")
	}

	// release date filters
	if p.ReleaseDateGTE != nil {
		where = append(where, "m.release_date >= "+addArg(*p.ReleaseDateGTE))
	}
	if p.ReleaseDateLTE != nil {
		where = append(where, "m.release_date <= "+addArg(*p.ReleaseDateLTE))
	}

	// vote average filters
	if p.VoteAvgGTE != nil {
		where = append(where, "ms.vote_average >= "+addArg(*p.VoteAvgGTE))
	}
	if p.VoteAvgLTE != nil {
		where = append(where, "ms.vote_average <= "+addArg(*p.VoteAvgLTE))
	}

	// genres: pass pq.Array(p.WithGenres) and correctly reference mg/genre columns
	if len(p.WithGenres) > 0 {
		// ensure we actually have concrete values (dedupe if you want)
		genrePos := addArg(pq.Array(p.WithGenres))

		if p.WithGenresAND {
			// Correct HAVING counts distinct matched genre_ids, not movie_id
			where = append(where,
				fmt.Sprintf("m.id IN (SELECT mg2.movie_id FROM movie_genres mg2 WHERE mg2.genre_id = ANY(%s::uuid[]) GROUP BY mg2.movie_id HAVING COUNT(DISTINCT mg2.genre_id) = %d)",
					genrePos, len(p.WithGenres)))
		} else {
			// EXISTS: check mg3.genre_id (not g.id) and use the placeholder
			where = append(where,
				fmt.Sprintf("EXISTS (SELECT 1 FROM movie_genres mg3 WHERE mg3.movie_id = m.id AND mg3.genre_id = ANY(%s::uuid[]))",
					genrePos))
		}
	}

	// ordering (kept your logic, with corrected column names)
	orderBy := "ms.popularity DESC NULLS LAST, m.created_at DESC"
	if p.SortBy != "" {
		parts := strings.SplitN(p.SortBy, ".", 2)
		sortMap := map[string]string{
			"popularity":   "ms.popularity",
			"release_date": "m.release_date",
			"vote_average": "ms.vote_average",
		}
		field := parts[0]
		dir := "desc"
		if len(parts) == 2 && (strings.ToLower(parts[1]) == "asc" || strings.ToLower(parts[1]) == "desc") {
			dir = strings.ToLower(parts[1])
		}
		if col, ok := sortMap[field]; ok {
			orderBy = fmt.Sprintf("%s %s NULLS LAST", col, dir)
		}
	}

	// limit & offset
	limitPlaceholder := addArg(p.PageSize)
	offsetPlaceholder := addArg(offset)

	sqlStr := fmt.Sprintf("%s %s WHERE %s ORDER BY %s LIMIT %s OFFSET %s",
		BaseQuery, fromWhere, strings.Join(where, " AND "), orderBy, limitPlaceholder, offsetPlaceholder)

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("Error Querying Discover movies: %w", err)
	}
	defer rows.Close()

	var items []model.DiscoverItem
	totalCount := 0

	for rows.Next() {
		var (
			id          uuid.UUID
			title       string
			overview    sql.NullString
			releaseDate sql.NullString
			voteAvg     sql.NullFloat64
			voteCount   sql.NullInt64
			poster      sql.NullString
			backdrop    sql.NullString
			popularity  sql.NullFloat64
			genreIDs    pq.StringArray
			total       sql.NullInt64
		)

		if err := rows.Scan(&id, &title, &overview, &releaseDate, &voteAvg, &voteCount, &poster, &backdrop, &popularity, &genreIDs, &total); err != nil {
			return nil, 0, fmt.Errorf("Error on rows discovermovie: %w", err)
		}

		it := model.DiscoverItem{ID: id, Title: title}
		if overview.Valid {
			it.Overview = &overview.String
		}
		if releaseDate.Valid {
			it.ReleaseDate = &releaseDate.String
		}
		if voteAvg.Valid {
			v := voteAvg.Float64
			it.VoteAverage = &v
		}
		if voteCount.Valid {
			vc := int(voteCount.Int64)
			it.VoteCount = &vc
		}
		if poster.Valid {
			it.PosterPath = &poster.String
		}
		if backdrop.Valid {
			it.BackdropPath = &backdrop.String
		}
		if popularity.Valid {
			pv := popularity.Float64
			it.Popularity = &pv
		}
		it.GenreIDs = []string(genreIDs)

		items = append(items, it)
		if total.Valid {
			totalCount = int(total.Int64)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration: %w", err)
	}

	return items, totalCount, nil
}
