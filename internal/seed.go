package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/h-raju-arch/movie_app_backend/internal/db"
	_ "github.com/lib/pq"
)

func newUUID() uuid.UUID {
	u, err := uuid.NewV7()
	if err != nil {
		log.Fatal("failed to generate uuidv7", err)
	}
	return u
}

func mustExec(ctx context.Context, execer execer, query string, args ...any) {
	if _, err := execer.ExecContext(ctx, query, args...); err != nil {
		log.Fatalf("exec failed: %v\nquery: %s\nargs: %#v", err, query, args)
	}
}

// execer allows using *sql.DB or *sql.Tx
type execer interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func main() {
	database := db.Open()
	defer database.Close()

	ctx := context.Background()
	log.Println("seeding database..")

	// Use a transaction so partial seeds don't persist on failure
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("begin tx: %v", err)
	}
	// if anything fails we will rollback at the end
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	now := time.Now()

	// --------------------
	// GENRES (12)
	// --------------------
	genreNames := []string{
		"Action", "Adventure", "Animation", "Comedy",
		"Crime", "Documentary", "Drama", "Family",
		"Fantasy", "Historical", "Horror", "Science Fiction",
	}
	genreIDs := make([]uuid.UUID, 0, len(genreNames))
	for _, g := range genreNames {
		id := newUUID()
		genreIDs = append(genreIDs, id)
		mustExec(ctx, tx, `INSERT INTO genres (id, name) VALUES ($1, $2)`, id, g)
	}

	// --------------------
	// COMPANIES (12)
	// --------------------
	companyData := []struct {
		Name    string
		Country string
	}{
		{"Warner Bros.", "US"},
		{"Universal Pictures", "US"},
		{"Paramount Pictures", "US"},
		{"20th Century Studios", "US"},
		{"Columbia Pictures", "US"},
		{"Netflix", "US"},
		{"Amazon Studios", "US"},
		{"Studio Ghibli", "JP"},
		{"Toho", "JP"},
		{"BBC Films", "GB"},
		{"A24", "US"},
		{"Lionsgate", "US"},
	}
	companyIDs := make([]uuid.UUID, 0, len(companyData))
	for _, c := range companyData {
		id := newUUID()
		companyIDs = append(companyIDs, id)
		mustExec(ctx, tx, `INSERT INTO companies (id, name, origin_country) VALUES ($1, $2, $3)`, id, c.Name, c.Country)
	}

	// --------------------
	// MOVIES (15)
	// --------------------
	type movieSeed struct {
		Title    string
		Lang     string
		Overview string
		Release  string // YYYY-MM-DD
		Runtime  int
		Adult    bool
	}
	movies := []movieSeed{
		{"Edge of Tomorrow", "en", "A soldier relives the same day over and over again.", "2014-06-06", 113, false},
		{"Interstellar", "en", "A team travels through a wormhole in space.", "2014-11-07", 169, false},
		{"Inception", "en", "A thief who steals corporate secrets through dream-sharing.", "2010-07-16", 148, false},
		{"The Dark Knight", "en", "Batman faces the Joker.", "2008-07-18", 152, false},
		{"Spirited Away", "ja", "A young girl enters a world of spirits.", "2001-07-20", 125, false},
		{"Parasite", "ko", "A dark satire about class and family.", "2019-05-30", 132, false},
		{"Mad Max: Fury Road", "en", "In a post-apocalyptic wasteland, Max teams with Furiosa.", "2015-05-15", 120, false},
		{"The Matrix", "en", "A hacker discovers the true nature of reality.", "1999-03-31", 136, false},
		{"The Shawshank Redemption", "en", "Two imprisoned men bond over years.", "1994-09-22", 142, false},
		{"The Godfather", "en", "The aging patriarch transfers control of his empire.", "1972-03-24", 175, false},
		{"Get Out", "en", "A young African-American visits his white girlfriend's family.", "2017-02-24", 104, false},
		{"The Conjuring", "en", "Paranormal investigators help a family.", "2013-07-19", 112, false},
		{"Your Name", "ja", "Two teens share a profound connection.", "2016-08-26", 107, false},
		{"Blade Runner 2049", "en", "A young blade runner discovers a long-buried secret.", "2017-10-06", 164, false},
		{"The Irishman", "en", "A mob hitman's possible involvement in historic events.", "2019-11-27", 209, false},
	}
	movieIDs := make([]uuid.UUID, 0, len(movies))
	for _, m := range movies {
		id := newUUID()
		movieIDs = append(movieIDs, id)
		mustExec(ctx, tx, `
			INSERT INTO movies (
				id, title, original_language, overview, release_date, runtime, adult, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
		`, id, m.Title, m.Lang, m.Overview, m.Release, m.Runtime, m.Adult, now)
	}

	// --------------------
	// MOVIE STATS (15)
	// --------------------
	// Use some made-up popularity / vote_average / vote_count values
	stats := []struct {
		Popularity  float64
		VoteAverage float64
		VoteCount   int
	}{
		{85.4, 8.1, 12450},
		{92.7, 8.6, 19870},
		{87.3, 8.8, 21000},
		{95.0, 9.0, 22000},
		{88.0, 8.6, 17000},
		{89.5, 8.6, 15000},
		{82.0, 8.1, 9800},
		{90.1, 8.7, 18000},
		{91.0, 9.3, 23000},
		{94.5, 9.2, 25000},
		{80.2, 7.7, 11000},
		{78.0, 7.5, 10500},
		{86.6, 8.4, 14000},
		{84.9, 8.0, 12500},
		{75.4, 7.8, 7200},
	}
	for i, s := range stats {
		mustExec(ctx, tx, `INSERT INTO movie_stats (movie_id, popularity, vote_average, vote_count) VALUES ($1, $2, $3, $4)`,
			movieIDs[i], s.Popularity, s.VoteAverage, s.VoteCount)
	}

	// --------------------
	// MOVIE ↔ GENRES mappings
	// --------------------
	// Map movies to a couple of genres each
	movieToGenres := [][]int{
		{0, 11}, // Edge of Tomorrow: Action, Sci-Fi
		{8, 11}, // Interstellar: Fantasy? Sci-Fi (using indexes)
		{0, 11}, // Inception
		{0, 4},  // The Dark Knight -> Action, Crime
		{2, 7},  // Spirited Away -> Animation, Family
		{6, 4},  // Parasite -> Drama, Crime
		{0, 11}, // Mad Max -> Action, Sci-Fi
		{11, 6}, // The Matrix -> Sci-Fi, Drama
		{6, 1},  // Shawshank -> Drama, Adventure
		{6, 4},  // The Godfather -> Drama, Crime
		{4, 3},  // Get Out -> Crime, Comedy (loosely)
		{10, 6}, // The Conjuring -> Horror, Drama
		{2, 7},  // Your Name -> Animation, Family
		{11, 6}, // Blade Runner 2049 -> Sci-Fi, Drama
		{6, 4},  // The Irishman -> Drama, Crime
	}
	for mi, gids := range movieToGenres {
		for _, gidx := range gids {
			mustExec(ctx, tx, `INSERT INTO movie_genres (movie_id, genre_id) VALUES ($1, $2)`, movieIDs[mi], genreIDs[gidx])
		}
	}

	// --------------------
	// MOVIE ↔ COMPANIES mappings
	// --------------------
	// Assign one or two companies per movie
	movieToCompanies := [][]int{
		{0},  // Warner
		{1},  // Universal
		{2},  // Paramount
		{0},  // Warner
		{7},  // Studio Ghibli
		{10}, // A24 (Parasite was actually CJ/Barunson etc., but for seed it's okay)
		{0},  // Warner
		{3},  // 20th Century
		{4},  // Columbia
		{2},  // Paramount
		{10}, // A24
		{11}, // Lionsgate (Conjuring was New Line, fine)
		{7},  // Studio Ghibli / Toho
		{1},  // Universal
		{5},  // Netflix
	}
	for mi, cidxs := range movieToCompanies {
		for _, cidx := range cidxs {
			mustExec(ctx, tx, `INSERT INTO movie_companies (movie_id, company_id) VALUES ($1, $2)`, movieIDs[mi], companyIDs[cidx])
		}
	}

	// --------------------
	// PEOPLE (20)
	// --------------------
	people := []struct {
		Name     string
		KnownFor string
	}{
		{"Tom Cruise", "Acting"},
		{"Matthew McConaughey", "Acting"},
		{"Christopher Nolan", "Directing"},
		{"Leonardo DiCaprio", "Acting"},
		{"Brad Pitt", "Acting"},
		{"Quentin Tarantino", "Directing"},
		{"Hayao Miyazaki", "Directing"},
		{"Bong Joon-ho", "Directing"},
		{"George Miller", "Directing"},
		{"Keanu Reeves", "Acting"},
		{"Morgan Freeman", "Acting"},
		{"Al Pacino", "Acting"},
		{"Jordan Peele", "Directing"},
		{"James Wan", "Directing"},
		{"Makoto Shinkai", "Directing"},
		{"Ryan Gosling", "Acting"},
		{"Harrison Ford", "Acting"},
		{"Robert De Niro", "Acting"},
		{"Martin Scorsese", "Directing"},
		{"Adam Driver", "Acting"},
	}
	personIDs := make([]uuid.UUID, 0, len(people))
	for _, p := range people {
		id := newUUID()
		personIDs = append(personIDs, id)
		mustExec(ctx, tx, `INSERT INTO people (id, name, known_for) VALUES ($1, $2, $3)`, id, p.Name, p.KnownFor)
	}

	// --------------------
	// CREDITS (cast & crew) - create some mappings
	// --------------------
	// We'll create multiple credits and map to movies by index
	type creditSeed struct {
		MovieIndex int
		PersonIdx  int
		CType      string
		Character  *string
		CastOrder  *int
	}
	char1 := "William Cage"
	char2 := "Cooper"
	char3 := "Dom Cobb"
	char4 := "Bruce Wayne"
	char5 := "Chihiro Ogino"
	char6 := "Kim Ki-taek"
	char7 := "Imperator Furiosa"
	char8 := "Neo"
	char9 := "Andy Dufresne"
	char10 := "Michael Corleone"

	ord1 := 1
	ord2 := 1
	ord3 := 1
	ord4 := 1
	ord5 := 1
	ord6 := 1

	credits := []creditSeed{
		{0, 0, "cast", &char1, &ord1},   // movie 0 -> Tom Cruise
		{1, 1, "cast", &char2, &ord2},   // movie 1 -> Matthew
		{2, 3, "cast", &char3, &ord3},   // Inception -> Leo (for seed)
		{3, 4, "cast", &char4, &ord4},   // Dark Knight -> Brad Pitt (example)
		{4, 6, "crew", &char5, nil},     // Spirited Away -> Hayao Miyazaki (director)
		{5, 7, "crew", &char6, nil},     // Parasite -> Bong Joon-ho (director)
		{6, 8, "crew", &char7, nil},     // Mad Max -> George Miller
		{7, 9, "cast", &char8, &ord5},   // The Matrix -> Keanu Reeves
		{8, 10, "cast", &char9, &ord6},  // Shawshank -> Morgan Freeman
		{9, 11, "cast", &char10, &ord1}, // Godfather -> Al Pacino
		{10, 12, "crew", nil, nil},      // Get Out -> Jordan Peele
		{11, 13, "crew", nil, nil},      // The Conjuring -> James Wan
		{12, 14, "crew", nil, nil},      // Your Name -> Makoto Shinkai
		{13, 15, "cast", nil, nil},      // Blade Runner -> Ryan Gosling (example)
		{14, 18, "crew", nil, nil},      // The Irishman -> Martin Scorsese
	}

	for _, cs := range credits {
		credID := newUUID()
		var charParam interface{}
		var orderParam interface{}
		if cs.Character != nil {
			charParam = *cs.Character
		} else {
			charParam = nil
		}
		if cs.CastOrder != nil {
			orderParam = *cs.CastOrder
		} else {
			orderParam = nil
		}
		mustExec(ctx, tx, `
			INSERT INTO credits (id, movie_id, person_id, credit_type, character_name, cast_order)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, credID, movieIDs[cs.MovieIndex], personIDs[cs.PersonIdx], cs.CType, charParam, orderParam)
	}

	// --------------------
	// IMAGES (at least 15)
	// --------------------
	type imageSeed struct {
		MovieIndex int
		Path       string
		Type       string
		Width      int
		Height     int
		Lang       *string
	}
	langEn := "en"
	langJa := "ja"
	imageSeeds := []imageSeed{
		{0, "/posters/edge_of_tomorrow.jpg", "poster", 500, 750, &langEn},
		{0, "/backdrops/edge_of_tomorrow.jpg", "backdrop", 1920, 1080, nil},
		{1, "/posters/interstellar.jpg", "poster", 500, 750, &langEn},
		{2, "/posters/inception.jpg", "poster", 500, 750, &langEn},
		{3, "/posters/dark_knight.jpg", "poster", 500, 750, &langEn},
		{4, "/posters/spirited_away.jpg", "poster", 500, 750, &langJa},
		{5, "/posters/parasite.jpg", "poster", 500, 750, nil},
		{6, "/posters/mad_max_fury_road.jpg", "poster", 500, 750, nil},
		{7, "/posters/matrix.jpg", "poster", 500, 750, &langEn},
		{8, "/posters/shawshank.jpg", "poster", 500, 750, &langEn},
		{9, "/posters/godfather.jpg", "poster", 500, 750, &langEn},
		{10, "/posters/get_out.jpg", "poster", 500, 750, &langEn},
		{11, "/posters/conjuring.jpg", "poster", 500, 750, &langEn},
		{12, "/posters/your_name.jpg", "poster", 500, 750, &langJa},
		{13, "/posters/blade_runner_2049.jpg", "poster", 500, 750, &langEn},
		{14, "/posters/the_irishman.jpg", "poster", 500, 750, &langEn},
		// add some backdrops too:
		{2, "/backdrops/inception.jpg", "backdrop", 1920, 1080, &langEn},
		{5, "/backdrops/parasite.jpg", "backdrop", 1920, 1080, nil},
	}
	for _, img := range imageSeeds {
		id := newUUID()
		var langParam interface{}
		if img.Lang != nil {
			langParam = *img.Lang
		} else {
			langParam = nil
		}
		mustExec(ctx, tx, `
			INSERT INTO images (id, movie_id, file_path, type, width, height, language)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, id, movieIDs[img.MovieIndex], img.Path, img.Type, img.Width, img.Height, langParam)
	}

	// commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("tx commit: %v", err)
	}

	log.Println("✅ Database seeded successfully")
}
