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
	log.Println("Seeding database...")

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

	// ====================
	// LANGUAGES (12)
	// ====================
	languages := []struct {
		ISO  string
		Name string
	}{
		{"en", "English"},
		{"ja", "Japanese"},
		{"ko", "Korean"},
		{"es", "Spanish"},
		{"fr", "French"},
		{"de", "German"},
		{"it", "Italian"},
		{"zh", "Chinese"},
		{"hi", "Hindi"},
		{"pt", "Portuguese"},
		{"ru", "Russian"},
		{"ar", "Arabic"},
	}
	for _, l := range languages {
		mustExec(ctx, tx, `INSERT INTO languages (iso_639_1, name) VALUES ($1, $2) ON CONFLICT DO NOTHING`, l.ISO, l.Name)
	}
	log.Println("  - Languages seeded")

	// ====================
	// GENRES (15)
	// ====================
	genreNames := []string{
		"Action", "Adventure", "Animation", "Comedy", "Crime",
		"Documentary", "Drama", "Family", "Fantasy", "Horror",
		"Mystery", "Romance", "Science Fiction", "Thriller", "War",
	}
	genreIDs := make(map[string]uuid.UUID)
	for _, g := range genreNames {
		id := newUUID()
		genreIDs[g] = id
		mustExec(ctx, tx, `INSERT INTO genres (id, name) VALUES ($1, $2)`, id, g)
	}
	log.Println("  - Genres seeded")

	// ====================
	// COMPANIES (18)
	// ====================
	companyData := []struct {
		Name     string
		Country  string
		Homepage string
	}{
		{"Warner Bros. Pictures", "US", "https://www.warnerbros.com"},
		{"Universal Pictures", "US", "https://www.universalpictures.com"},
		{"Paramount Pictures", "US", "https://www.paramount.com"},
		{"20th Century Studios", "US", "https://www.20thcenturystudios.com"},
		{"Columbia Pictures", "US", "https://www.sonypictures.com"},
		{"Walt Disney Pictures", "US", "https://www.disney.com"},
		{"Netflix", "US", "https://www.netflix.com"},
		{"Amazon Studios", "US", "https://studios.amazon.com"},
		{"A24", "US", "https://a24films.com"},
		{"Lionsgate", "US", "https://www.lionsgate.com"},
		{"Studio Ghibli", "JP", "https://www.ghibli.jp"},
		{"Toho", "JP", "https://www.toho.co.jp"},
		{"CJ Entertainment", "KR", "https://www.cjenm.com"},
		{"BBC Films", "GB", "https://www.bbc.co.uk/bbcfilm"},
		{"Canal+", "FR", "https://www.canalplus.com"},
		{"Gaumont", "FR", "https://www.gaumont.com"},
		{"Legendary Entertainment", "US", "https://www.legendary.com"},
		{"Blumhouse Productions", "US", "https://www.blumhouse.com"},
	}
	companyIDs := make(map[string]uuid.UUID)
	for _, c := range companyData {
		id := newUUID()
		companyIDs[c.Name] = id
		mustExec(ctx, tx, `INSERT INTO companies (id, name, origin_country, homepage) VALUES ($1, $2, $3, $4)`,
			id, c.Name, c.Country, c.Homepage)
	}
	log.Println("  - Companies seeded")

	// ====================
	// PEOPLE (35) - Actors, Directors, Writers
	// ====================
	peopleData := []struct {
		Name        string
		KnownFor    string
		ProfilePath string
	}{
		// Directors
		{"Christopher Nolan", "Directing", "/cGOPbv9wA5gEejkUN892JrveARt.jpg"},
		{"Steven Spielberg", "Directing", "/tZxcg19YQ3e8fJ0pOs7hjlnmmr6.jpg"},
		{"Martin Scorsese", "Directing", "/9U9Y5GQuWX3EZy39B8nkk4NY01S.jpg"},
		{"Quentin Tarantino", "Directing", "/1gjcpAa99FAOWGnrUvHEXXsRs7o.jpg"},
		{"Denis Villeneuve", "Directing", "/zdDx9Xs93UIrJFWYApYR28J8M6b.jpg"},
		{"Bong Joon-ho", "Directing", "/tKLJBqbdH6HFj2QxLA5o8Zk7IVs.jpg"},
		{"Hayao Miyazaki", "Directing", "/mG3cfxtA5jqDc7fpKgyzZMKoXDh.jpg"},
		{"Makoto Shinkai", "Directing", "/yTjVpqkmGLMKUJjxrro1cq5bYgK.jpg"},
		{"Jordan Peele", "Directing", "/kFUKn5g3ebPSOAMUH2wJ1jHg7BQ.jpg"},
		{"Greta Gerwig", "Directing", "/9xYRFjfYDlIVPHlxCQp4xCFKviA.jpg"},
		// Actors
		{"Leonardo DiCaprio", "Acting", "/wo2hJpn04vbtmh0B9utCFdsQhxM.jpg"},
		{"Tom Hanks", "Acting", "/xndWFsBlClOJFRdhSt4NBwiPq2o.jpg"},
		{"Robert De Niro", "Acting", "/cT8htcckIuyI1Lqwt1CvD02ynTh.jpg"},
		{"Meryl Streep", "Acting", "/emAAzyK1sz6bIG0Ni4gLejhSYhh.jpg"},
		{"Scarlett Johansson", "Acting", "/6NsMbJXRlDZuDzatN2akFdGuTvx.jpg"},
		{"Brad Pitt", "Acting", "/cckcYc2v0yh1tc9QjRelptcOBko.jpg"},
		{"Tom Cruise", "Acting", "/eOh4ubpOm2Igdg0QH2ghj0mFtC.jpg"},
		{"Keanu Reeves", "Acting", "/4D0PpNI0kmP58hgrwGC3wCjxhnm.jpg"},
		{"Morgan Freeman", "Acting", "/oIciQWr8VwKoR8TmAw1owaiZFyb.jpg"},
		{"Al Pacino", "Acting", "/ks7Ba8x9EJr2L1uQ2K45iWzrT5E.jpg"},
		{"Ryan Gosling", "Acting", "/lyUyVARQKhGxaxy0FbPJCQRpiaW.jpg"},
		{"Emma Stone", "Acting", "/eWjkPYeXkgkJLAqLPdj4pLqLUXo.jpg"},
		{"Margot Robbie", "Acting", "/euDPyqLnuwaWMHgjkgniUvrgQBi.jpg"},
		{"Timothee Chalamet", "Acting", "/BE2sdjpgsa2rNTFa66f7upkaOP.jpg"},
		{"Florence Pugh", "Acting", "/fhEsn35uSwqUQPwDpK3Oc2xsULk.jpg"},
		{"Song Kang-ho", "Acting", "/dVyYvqYzxs8pxCq0odZBDDUUKAb.jpg"},
		{"Park So-dam", "Acting", "/2Hd7M7kQF7nEpW2akDtYOFCqxo.jpg"},
		// Writers
		{"Aaron Sorkin", "Writing", "/yle4Y79SDNqUnuHFsYwHqT6a6mH.jpg"},
		{"Charlie Kaufman", "Writing", "/4U3OqTxfHPIRnDSvhjqLz5P8gqF.jpg"},
		{"Diablo Cody", "Writing", "/4thGZbKNRuBWgVPdPL8swGG3m2F.jpg"},
		// Cinematographers
		{"Roger Deakins", "Camera", "/rB3Epv97LqvToYLyiUBaLKjIbHy.jpg"},
		{"Emmanuel Lubezki", "Camera", "/5mW7RTt8C8iQxS6tLqsD8qYnYhd.jpg"},
		// Composers
		{"Hans Zimmer", "Sound", "/tpQnDeHY15szIXvpnhlprufz4d.jpg"},
		{"John Williams", "Sound", "/KL0BK3V7m5Q1J7Xay9H1FnS4BGz.jpg"},
		{"Ludwig Goransson", "Sound", "/kxPVz9gNvTJkrL0gLPJgKVGqGrs.jpg"},
	}
	personIDs := make(map[string]uuid.UUID)
	for _, p := range peopleData {
		id := newUUID()
		personIDs[p.Name] = id
		mustExec(ctx, tx, `INSERT INTO people (id, name, known_for, profile_path) VALUES ($1, $2, $3, $4)`,
			id, p.Name, p.KnownFor, p.ProfilePath)
	}
	log.Println("  - People seeded")

	// ====================
	// MOVIES (25)
	// ====================
	type movieSeed struct {
		Title            string
		OriginalTitle    string
		Lang             string
		TagLine          string
		Overview         string
		Release          string
		Runtime          int
		Adult            bool
		Budget           int64
		Revenue          int64
		Homepage         string
		PosterPath       string
		BackdropPath     string
		Genres           []string
		Companies        []string
		SpokenLanguages  []string
		Popularity       float64
		VoteAverage      float64
		VoteCount        int
		CastMembers      []string
		Director         string
		Writer           string
		Composer         string
		HasTranslations  bool
	}

	movies := []movieSeed{
		{
			Title: "Inception", OriginalTitle: "Inception", Lang: "en",
			TagLine:  "Your mind is the scene of the crime.",
			Overview: "A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O.",
			Release:  "2010-07-16", Runtime: 148, Adult: false,
			Budget: 160000000, Revenue: 836836967,
			Homepage:     "https://www.warnerbros.com/movies/inception",
			PosterPath:   "/9gk7adHYeDvHkCSEqAvQNLV5Ber.jpg",
			BackdropPath: "/s3TBrRGB1iav7gFOCNx3H31MoES.jpg",
			Genres:       []string{"Action", "Science Fiction", "Adventure"},
			Companies:    []string{"Warner Bros. Pictures", "Legendary Entertainment"},
			SpokenLanguages: []string{"en", "ja", "fr"},
			Popularity: 87.3, VoteAverage: 8.4, VoteCount: 34500,
			CastMembers: []string{"Leonardo DiCaprio", "Tom Hanks"},
			Director: "Christopher Nolan", Writer: "Christopher Nolan", Composer: "Hans Zimmer",
			HasTranslations: true,
		},
		{
			Title: "The Dark Knight", OriginalTitle: "The Dark Knight", Lang: "en",
			TagLine:  "Why so serious?",
			Overview: "Batman raises the stakes in his war on crime. With the help of Lt. Jim Gordon and District Attorney Harvey Dent, Batman sets out to dismantle the remaining criminal organizations that plague the streets.",
			Release:  "2008-07-18", Runtime: 152, Adult: false,
			Budget: 185000000, Revenue: 1004558444,
			Homepage:     "https://www.warnerbros.com/movies/dark-knight",
			PosterPath:   "/qJ2tW6WMUDux911r6m7haRef0WH.jpg",
			BackdropPath: "/nMKdUUepR0i5zn0y1T4CsSB5chy.jpg",
			Genres:       []string{"Action", "Crime", "Drama", "Thriller"},
			Companies:    []string{"Warner Bros. Pictures", "Legendary Entertainment"},
			SpokenLanguages: []string{"en"},
			Popularity: 95.0, VoteAverage: 9.0, VoteCount: 29800,
			CastMembers: []string{"Brad Pitt", "Morgan Freeman"},
			Director: "Christopher Nolan", Writer: "Christopher Nolan", Composer: "Hans Zimmer",
			HasTranslations: true,
		},
		{
			Title: "Interstellar", OriginalTitle: "Interstellar", Lang: "en",
			TagLine:  "Mankind was born on Earth. It was never meant to die here.",
			Overview: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
			Release:  "2014-11-07", Runtime: 169, Adult: false,
			Budget: 165000000, Revenue: 677471339,
			Homepage:     "https://www.interstellarmovie.com",
			PosterPath:   "/gEU2QniE6E77NI6lCU6MxlNBvIx.jpg",
			BackdropPath: "/xJHokMbljvjADYdit5fK5VQsXEG.jpg",
			Genres:       []string{"Adventure", "Drama", "Science Fiction"},
			Companies:    []string{"Warner Bros. Pictures", "Paramount Pictures", "Legendary Entertainment"},
			SpokenLanguages: []string{"en"},
			Popularity: 92.7, VoteAverage: 8.6, VoteCount: 32100,
			CastMembers: []string{"Leonardo DiCaprio", "Tom Hanks"},
			Director: "Christopher Nolan", Writer: "Christopher Nolan", Composer: "Hans Zimmer",
			HasTranslations: true,
		},
		{
			Title: "Parasite", OriginalTitle: "기생충", Lang: "ko",
			TagLine:  "Act like you own the place.",
			Overview: "All unemployed, Ki-taek's family takes peculiar interest in the wealthy and glamorous Parks for their livelihood until they get entangled in an unexpected incident.",
			Release:  "2019-05-30", Runtime: 132, Adult: false,
			Budget: 11400000, Revenue: 258773700,
			Homepage:     "https://www.parasite-movie.com",
			PosterPath:   "/7IiTTgloJzvGI1TAYymCfbfl3vT.jpg",
			BackdropPath: "/TU9NIjwzjoKPwQHoHshkFcQUCG.jpg",
			Genres:       []string{"Comedy", "Thriller", "Drama"},
			Companies:    []string{"CJ Entertainment"},
			SpokenLanguages: []string{"ko", "en"},
			Popularity: 89.5, VoteAverage: 8.5, VoteCount: 16200,
			CastMembers: []string{"Song Kang-ho", "Park So-dam"},
			Director: "Bong Joon-ho", Writer: "Bong Joon-ho", Composer: "Ludwig Goransson",
			HasTranslations: true,
		},
		{
			Title: "Spirited Away", OriginalTitle: "千と千尋の神隠し", Lang: "ja",
			TagLine:  "A young girl enters a world ruled by spirits.",
			Overview: "A young girl, Chihiro, becomes trapped in a strange new world of spirits. When her parents undergo a mysterious transformation, she must call upon the courage she never knew she had.",
			Release:  "2001-07-20", Runtime: 125, Adult: false,
			Budget: 19000000, Revenue: 395580000,
			Homepage:     "https://www.ghibli.jp/works/chihiro",
			PosterPath:   "/39wmItIWsg5sZMyRUHLkWBcuVCM.jpg",
			BackdropPath: "/6oaL4DP75yABrd5EbC4H2zq5ghc.jpg",
			Genres:       []string{"Animation", "Family", "Fantasy"},
			Companies:    []string{"Studio Ghibli", "Toho"},
			SpokenLanguages: []string{"ja"},
			Popularity: 88.0, VoteAverage: 8.5, VoteCount: 14800,
			CastMembers: []string{},
			Director: "Hayao Miyazaki", Writer: "Hayao Miyazaki", Composer: "Hans Zimmer",
			HasTranslations: true,
		},
		{
			Title: "Your Name", OriginalTitle: "君の名は。", Lang: "ja",
			TagLine:  "What's your name?",
			Overview: "High schoolers Mitsuha and Taki are complete strangers living separate lives. But one night, they suddenly switch places. This bizarre occurrence continues to happen randomly, and the two must adjust their lives around each other.",
			Release:  "2016-08-26", Runtime: 107, Adult: false,
			Budget: 25000000, Revenue: 380140450,
			Homepage:     "https://www.kiminona.com",
			PosterPath:   "/q719jXXEzOoYaps6babgKnONONX.jpg",
			BackdropPath: "/dIWwZW7dJJtqC6CgWzYkNVKIUm8.jpg",
			Genres:       []string{"Animation", "Romance", "Drama"},
			Companies:    []string{"Toho"},
			SpokenLanguages: []string{"ja"},
			Popularity: 86.6, VoteAverage: 8.4, VoteCount: 10500,
			CastMembers: []string{},
			Director: "Makoto Shinkai", Writer: "Makoto Shinkai", Composer: "Ludwig Goransson",
			HasTranslations: true,
		},
		{
			Title: "The Matrix", OriginalTitle: "The Matrix", Lang: "en",
			TagLine:  "Welcome to the Real World.",
			Overview: "A computer hacker learns from mysterious rebels about the true nature of his reality and his role in the war against its controllers.",
			Release:  "1999-03-31", Runtime: 136, Adult: false,
			Budget: 63000000, Revenue: 466364845,
			Homepage:     "https://www.warnerbros.com/movies/matrix",
			PosterPath:   "/f89U3ADr1oiB1s9GkdPOEpXUk5H.jpg",
			BackdropPath: "/l4QHerTSbMI7qgvasqxP36pqjN6.jpg",
			Genres:       []string{"Action", "Science Fiction"},
			Companies:    []string{"Warner Bros. Pictures"},
			SpokenLanguages: []string{"en"},
			Popularity: 90.1, VoteAverage: 8.7, VoteCount: 24100,
			CastMembers: []string{"Keanu Reeves"},
			Director: "Steven Spielberg", Writer: "Aaron Sorkin", Composer: "Hans Zimmer",
			HasTranslations: false,
		},
		{
			Title: "The Shawshank Redemption", OriginalTitle: "The Shawshank Redemption", Lang: "en",
			TagLine:  "Fear can hold you prisoner. Hope can set you free.",
			Overview: "Framed in the 1940s for the double murder of his wife and her lover, upstanding banker Andy Dufresne begins a new life at the Shawshank prison, where he puts his accounting skills to work for an pointedly cruel warden.",
			Release:  "1994-09-23", Runtime: 142, Adult: false,
			Budget: 25000000, Revenue: 58300000,
			Homepage:     "",
			PosterPath:   "/q6y0Go1tsGEsmtFryDOJo3dEmqu.jpg",
			BackdropPath: "/kXfqcdQKsToO0OUXHcrrNCHDBzO.jpg",
			Genres:       []string{"Drama", "Crime"},
			Companies:    []string{"Columbia Pictures"},
			SpokenLanguages: []string{"en"},
			Popularity: 91.0, VoteAverage: 9.3, VoteCount: 25600,
			CastMembers: []string{"Morgan Freeman", "Tom Hanks"},
			Director: "Steven Spielberg", Writer: "Aaron Sorkin", Composer: "John Williams",
			HasTranslations: false,
		},
		{
			Title: "The Godfather", OriginalTitle: "The Godfather", Lang: "en",
			TagLine:  "An offer you can't refuse.",
			Overview: "Spanning the years 1945 to 1955, a chronicle of the fictional Italian-American Corleone crime family. When organized crime family patriarch, Vito Corleone barely survives an attempt on his life, his youngest son, Michael steps in to take care of the would-be killers.",
			Release:  "1972-03-14", Runtime: 175, Adult: false,
			Budget: 6000000, Revenue: 286000000,
			Homepage:     "",
			PosterPath:   "/3bhkrj58Vtu7enYsRolD1fZdja1.jpg",
			BackdropPath: "/rSPw7tgCH9c6NqICZef4kZjFOQ5.jpg",
			Genres:       []string{"Drama", "Crime"},
			Companies:    []string{"Paramount Pictures"},
			SpokenLanguages: []string{"en", "it"},
			Popularity: 94.5, VoteAverage: 9.2, VoteCount: 19200,
			CastMembers: []string{"Al Pacino", "Robert De Niro"},
			Director: "Martin Scorsese", Writer: "Aaron Sorkin", Composer: "John Williams",
			HasTranslations: false,
		},
		{
			Title: "Pulp Fiction", OriginalTitle: "Pulp Fiction", Lang: "en",
			TagLine:  "Just because you are a character doesn't mean you have character.",
			Overview: "A burger-loving hit man, his philosophical partner, a drug-addled gangster's moll and a washed-up boxer converge in this sprawling, comedic crime caper.",
			Release:  "1994-10-14", Runtime: 154, Adult: false,
			Budget: 8000000, Revenue: 213928762,
			Homepage:     "",
			PosterPath:   "/fIE3lAGcZDV1G6XM5KmuWnNsPp1.jpg",
			BackdropPath: "/suaEOtk1N1sgg2MTM7oZd2cfVp3.jpg",
			Genres:       []string{"Thriller", "Crime"},
			Companies:    []string{"A24", "Lionsgate"},
			SpokenLanguages: []string{"en", "es", "fr"},
			Popularity: 88.4, VoteAverage: 8.5, VoteCount: 26300,
			CastMembers: []string{"Brad Pitt", "Scarlett Johansson"},
			Director: "Quentin Tarantino", Writer: "Quentin Tarantino", Composer: "Hans Zimmer",
			HasTranslations: false,
		},
		{
			Title: "Get Out", OriginalTitle: "Get Out", Lang: "en",
			TagLine:  "Just because you're invited, doesn't mean you're welcome.",
			Overview: "Chris and his girlfriend Rose go upstate to visit her parents for the weekend. At first, Chris reads the family's overly accommodating behavior as nervous attempts to deal with their daughter's interracial relationship, but as the weekend progresses, a series of increasingly disturbing discoveries lead him to a truth that he never could have imagined.",
			Release:  "2017-02-24", Runtime: 104, Adult: false,
			Budget: 4500000, Revenue: 255457078,
			Homepage:     "https://www.getoutfilm.com",
			PosterPath:   "/qbaERlsQKIz0LHPZyNFpJT2KlP8.jpg",
			BackdropPath: "/hZkgoQYus5vegHoetLkCJzb17zJ.jpg",
			Genres:       []string{"Horror", "Mystery", "Thriller"},
			Companies:    []string{"Blumhouse Productions", "Universal Pictures"},
			SpokenLanguages: []string{"en"},
			Popularity: 80.2, VoteAverage: 7.6, VoteCount: 12800,
			CastMembers: []string{"Ryan Gosling", "Florence Pugh"},
			Director: "Jordan Peele", Writer: "Jordan Peele", Composer: "Ludwig Goransson",
			HasTranslations: false,
		},
		{
			Title: "Dune", OriginalTitle: "Dune", Lang: "en",
			TagLine:  "It begins.",
			Overview: "Paul Atreides, a brilliant and gifted young man born into a great destiny beyond his understanding, must travel to the most dangerous planet in the universe to ensure the future of his family and his people.",
			Release:  "2021-09-15", Runtime: 155, Adult: false,
			Budget: 165000000, Revenue: 402027830,
			Homepage:     "https://www.dunemovie.com",
			PosterPath:   "/d5NXSklXo0qyIYkgV94XAgMIckC.jpg",
			BackdropPath: "/jYEW5xZkZk2WTrdbMGAPFuBqbDc.jpg",
			Genres:       []string{"Science Fiction", "Adventure", "Drama"},
			Companies:    []string{"Warner Bros. Pictures", "Legendary Entertainment"},
			SpokenLanguages: []string{"en", "ar"},
			Popularity: 93.2, VoteAverage: 8.0, VoteCount: 11200,
			CastMembers: []string{"Timothee Chalamet", "Florence Pugh"},
			Director: "Denis Villeneuve", Writer: "Denis Villeneuve", Composer: "Hans Zimmer",
			HasTranslations: true,
		},
		{
			Title: "Blade Runner 2049", OriginalTitle: "Blade Runner 2049", Lang: "en",
			TagLine:  "The key to the future is finally unearthed.",
			Overview: "Thirty years after the events of the first film, a new blade runner, LAPD Officer K, unearths a long-buried secret that has the potential to plunge what's left of society into chaos.",
			Release:  "2017-10-06", Runtime: 164, Adult: false,
			Budget: 150000000, Revenue: 259239658,
			Homepage:     "https://www.bladerunner2049.com",
			PosterPath:   "/gajva2L0rPYkEWjzgFlBXCAVBE5.jpg",
			BackdropPath: "/ilRyazdMJwN05exqhwK4tMKBYZs.jpg",
			Genres:       []string{"Science Fiction", "Drama", "Mystery"},
			Companies:    []string{"Warner Bros. Pictures", "Columbia Pictures"},
			SpokenLanguages: []string{"en", "ru"},
			Popularity: 84.9, VoteAverage: 7.5, VoteCount: 13400,
			CastMembers: []string{"Ryan Gosling", "Leonardo DiCaprio"},
			Director: "Denis Villeneuve", Writer: "Aaron Sorkin", Composer: "Hans Zimmer",
			HasTranslations: false,
		},
		{
			Title: "The Irishman", OriginalTitle: "The Irishman", Lang: "en",
			TagLine:  "His story changed history.",
			Overview: "Pennsylvania, 1956. Frank Sheeran, a war veteran of Irish origin who works as a truck driver, accidentally meets mobster Russell Bufalino. Once Frank becomes his trusted man, Bufalino sends him to Chicago with the task of helping Jimmy Hoffa, a powerful union leader related to organized crime.",
			Release:  "2019-11-01", Runtime: 209, Adult: false,
			Budget: 159000000, Revenue: 8000000,
			Homepage:     "https://www.netflix.com/title/80175798",
			PosterPath:   "/mbm8k3GFhXS0ROd9AD1gqYbIFbM.jpg",
			BackdropPath: "/4Zay0v3d3YGMJ6ES4FplkrSMPFi.jpg",
			Genres:       []string{"Crime", "Drama"},
			Companies:    []string{"Netflix"},
			SpokenLanguages: []string{"en", "it"},
			Popularity: 75.4, VoteAverage: 7.8, VoteCount: 8200,
			CastMembers: []string{"Robert De Niro", "Al Pacino"},
			Director: "Martin Scorsese", Writer: "Aaron Sorkin", Composer: "Ludwig Goransson",
			HasTranslations: false,
		},
		{
			Title: "Mad Max: Fury Road", OriginalTitle: "Mad Max: Fury Road", Lang: "en",
			TagLine:  "What a lovely day.",
			Overview: "An apocalyptic story set in the furthest reaches of our planet, in a stark desert landscape where humanity is broken, and most everyone is crazed fighting for the necessities of life.",
			Release:  "2015-05-15", Runtime: 120, Adult: false,
			Budget: 150000000, Revenue: 378858340,
			Homepage:     "https://www.warnerbros.com/movies/mad-max-fury-road",
			PosterPath:   "/8tZYtuWezp8JbcsvHYO0O46tFbo.jpg",
			BackdropPath: "/phszHPFVhPHhMZgo0fWTKBDQsJA.jpg",
			Genres:       []string{"Action", "Adventure", "Science Fiction"},
			Companies:    []string{"Warner Bros. Pictures"},
			SpokenLanguages: []string{"en"},
			Popularity: 82.0, VoteAverage: 7.6, VoteCount: 21300,
			CastMembers: []string{"Tom Cruise", "Scarlett Johansson"},
			Director: "Steven Spielberg", Writer: "Aaron Sorkin", Composer: "Hans Zimmer",
			HasTranslations: false,
		},
		{
			Title: "La La Land", OriginalTitle: "La La Land", Lang: "en",
			TagLine:  "Here's to the fools who dream.",
			Overview: "Mia, an aspiring actress, serves lattes to movie stars in between auditions and Sebastian, a jazz musician, scrapes by playing cocktail party gigs in dingy bars, but as success mounts they are faced with decisions that begin to fray the fragile fabric of their love affair.",
			Release:  "2016-12-09", Runtime: 128, Adult: false,
			Budget: 30000000, Revenue: 446092357,
			Homepage:     "https://www.lionsgate.com/movies/la-la-land",
			PosterPath:   "/uDO8zWDhfWwoFdKS4fzkUJt0Rf0.jpg",
			BackdropPath: "/wqtaHWOEZ3rXDJ8c6ZZShulbo18.jpg",
			Genres:       []string{"Comedy", "Drama", "Romance"},
			Companies:    []string{"Lionsgate"},
			SpokenLanguages: []string{"en"},
			Popularity: 83.5, VoteAverage: 7.9, VoteCount: 15600,
			CastMembers: []string{"Ryan Gosling", "Emma Stone"},
			Director: "Steven Spielberg", Writer: "Charlie Kaufman", Composer: "Ludwig Goransson",
			HasTranslations: true,
		},
		{
			Title: "Barbie", OriginalTitle: "Barbie", Lang: "en",
			TagLine:  "She's everything. He's just Ken.",
			Overview: "Barbie and Ken are having the time of their lives in the colorful and seemingly perfect world of Barbie Land. However, when they get a chance to go to the real world, they soon discover the joys and perils of living among humans.",
			Release:  "2023-07-21", Runtime: 114, Adult: false,
			Budget: 145000000, Revenue: 1441866460,
			Homepage:     "https://www.barbie-themovie.com",
			PosterPath:   "/iuFNMS8U5cb6xfzi51Dbkovj7vM.jpg",
			BackdropPath: "/nHf61UzkfFno5X1ofIhugCPus2R.jpg",
			Genres:       []string{"Comedy", "Adventure", "Fantasy"},
			Companies:    []string{"Warner Bros. Pictures"},
			SpokenLanguages: []string{"en"},
			Popularity: 96.8, VoteAverage: 7.0, VoteCount: 8900,
			CastMembers: []string{"Margot Robbie", "Ryan Gosling"},
			Director: "Greta Gerwig", Writer: "Greta Gerwig", Composer: "Ludwig Goransson",
			HasTranslations: true,
		},
		{
			Title: "Oppenheimer", OriginalTitle: "Oppenheimer", Lang: "en",
			TagLine:  "The world forever changes.",
			Overview: "The story of J. Robert Oppenheimer's role in the development of the atomic bomb during World War II.",
			Release:  "2023-07-21", Runtime: 180, Adult: false,
			Budget: 100000000, Revenue: 952000000,
			Homepage:     "https://www.oppenheimermovie.com",
			PosterPath:   "/8Gxv8gSFCU0XGDykEGv7zR1n2ua.jpg",
			BackdropPath: "/rLb2cwF3Pazuxaj0sRXQ037tGI1.jpg",
			Genres:       []string{"Drama", "Thriller", "War"},
			Companies:    []string{"Universal Pictures"},
			SpokenLanguages: []string{"en", "de"},
			Popularity: 97.5, VoteAverage: 8.1, VoteCount: 7800,
			CastMembers: []string{"Leonardo DiCaprio", "Robert De Niro", "Florence Pugh"},
			Director: "Christopher Nolan", Writer: "Christopher Nolan", Composer: "Ludwig Goransson",
			HasTranslations: true,
		},
		{
			Title: "Everything Everywhere All at Once", OriginalTitle: "Everything Everywhere All at Once", Lang: "en",
			TagLine:  "The universe is so much bigger than you realize.",
			Overview: "An aging Chinese immigrant is swept up in an insane adventure, where she alone can save the world by exploring other universes connecting with the lives she could have led.",
			Release:  "2022-03-25", Runtime: 139, Adult: false,
			Budget: 14300000, Revenue: 141287838,
			Homepage:     "https://a24films.com/films/everything-everywhere-all-at-once",
			PosterPath:   "/w3LxiVYdWWRvEVdn5RYq6jIqkb1.jpg",
			BackdropPath: "/fOy2Jurz9k6RnJnMUMRDAgBwru2.jpg",
			Genres:       []string{"Action", "Adventure", "Science Fiction", "Comedy"},
			Companies:    []string{"A24"},
			SpokenLanguages: []string{"en", "zh"},
			Popularity: 91.2, VoteAverage: 7.8, VoteCount: 9100,
			CastMembers: []string{"Scarlett Johansson", "Timothee Chalamet"},
			Director: "Quentin Tarantino", Writer: "Charlie Kaufman", Composer: "Ludwig Goransson",
			HasTranslations: true,
		},
		{
			Title: "The Conjuring", OriginalTitle: "The Conjuring", Lang: "en",
			TagLine:  "Based on the true case files of the Warrens.",
			Overview: "Paranormal investigators Ed and Lorraine Warren work to help a family terrorized by a dark presence in their farmhouse. Forced to confront a powerful entity, the Warrens find themselves caught in the most terrifying case of their lives.",
			Release:  "2013-07-19", Runtime: 112, Adult: false,
			Budget: 20000000, Revenue: 319494638,
			Homepage:     "https://www.warnerbros.com/movies/conjuring",
			PosterPath:   "/wVYREutTvI2tmxr6ujrHT704wGF.jpg",
			BackdropPath: "/o5brynSBJvwJvaExY0DM4xKPTEk.jpg",
			Genres:       []string{"Horror", "Thriller", "Mystery"},
			Companies:    []string{"Warner Bros. Pictures", "Blumhouse Productions"},
			SpokenLanguages: []string{"en"},
			Popularity: 78.0, VoteAverage: 7.5, VoteCount: 11400,
			CastMembers: []string{"Morgan Freeman", "Emma Stone"},
			Director: "Jordan Peele", Writer: "Diablo Cody", Composer: "Hans Zimmer",
			HasTranslations: false,
		},
		{
			Title: "Edge of Tomorrow", OriginalTitle: "Edge of Tomorrow", Lang: "en",
			TagLine:  "Live. Die. Repeat.",
			Overview: "Major Bill Cage is an officer who has never seen a day of combat when he is unceremoniously demoted and dropped into combat. Killed within minutes, Cage now finds himself inexplicably thrown into a time loop—forcing him to live out the same brutal combat over and over.",
			Release:  "2014-06-06", Runtime: 113, Adult: false,
			Budget: 178000000, Revenue: 370541256,
			Homepage:     "https://www.warnerbros.com/movies/edge-tomorrow",
			PosterPath:   "/xjw5trHV8BMGhXPTqzVQFLfF9hc.jpg",
			BackdropPath: "/jA6cU5MYf2I4ejMGfFbVAXBflxJ.jpg",
			Genres:       []string{"Action", "Science Fiction"},
			Companies:    []string{"Warner Bros. Pictures"},
			SpokenLanguages: []string{"en"},
			Popularity: 85.4, VoteAverage: 7.9, VoteCount: 13500,
			CastMembers: []string{"Tom Cruise", "Scarlett Johansson"},
			Director: "Christopher Nolan", Writer: "Aaron Sorkin", Composer: "Hans Zimmer",
			HasTranslations: false,
		},
		{
			Title: "Amélie", OriginalTitle: "Le Fabuleux Destin d'Amélie Poulain", Lang: "fr",
			TagLine:  "She'll change your life.",
			Overview: "At a tiny Parisian café, the adorable yet painfully shy Amélie accidentally discovers a gift for helping others. Soon Amélie is spending her days as a matchmaker, guardian angel, and all-around do-gooder.",
			Release:  "2001-04-25", Runtime: 122, Adult: false,
			Budget: 10000000, Revenue: 174200000,
			Homepage:     "",
			PosterPath:   "/slVnvaH6fpW9A7J5k0FRo7oQ6AG.jpg",
			BackdropPath: "/nWs0auTqn2UaFGfTKtUE5tlTeBu.jpg",
			Genres:       []string{"Comedy", "Romance"},
			Companies:    []string{"Canal+", "Gaumont"},
			SpokenLanguages: []string{"fr"},
			Popularity: 76.3, VoteAverage: 8.0, VoteCount: 12100,
			CastMembers: []string{"Emma Stone", "Scarlett Johansson"},
			Director: "Greta Gerwig", Writer: "Charlie Kaufman", Composer: "Hans Zimmer",
			HasTranslations: true,
		},
		{
			Title: "Oldboy", OriginalTitle: "올드보이", Lang: "ko",
			TagLine:  "15 years of imprisonment, 5 days of vengeance.",
			Overview: "With no clue how he came to be imprisoned, drugged and gaslighted for 15 years, a desperate businessman seeks revenge on his captors.",
			Release:  "2003-11-21", Runtime: 120, Adult: true,
			Budget: 3000000, Revenue: 15000000,
			Homepage:     "",
			PosterPath:   "/pWDtjs568ZfOTMbURQBYuT4Qxka.jpg",
			BackdropPath: "/2t9fTnkwOrLMVdpMFxXPMNSrBvn.jpg",
			Genres:       []string{"Thriller", "Drama", "Mystery", "Action"},
			Companies:    []string{"CJ Entertainment"},
			SpokenLanguages: []string{"ko"},
			Popularity: 72.8, VoteAverage: 8.4, VoteCount: 5600,
			CastMembers: []string{"Song Kang-ho"},
			Director: "Bong Joon-ho", Writer: "Bong Joon-ho", Composer: "Ludwig Goransson",
			HasTranslations: true,
		},
		{
			Title: "Schindler's List", OriginalTitle: "Schindler's List", Lang: "en",
			TagLine:  "The list is life.",
			Overview: "The true story of how businessman Oskar Schindler saved over a thousand Jewish lives from the Nazis while they worked as slaves in his factory during World War II.",
			Release:  "1993-12-15", Runtime: 195, Adult: false,
			Budget: 22000000, Revenue: 321306305,
			Homepage:     "",
			PosterPath:   "/sF1U4EUQS8YHUYjNl3pMGNIQyr0.jpg",
			BackdropPath: "/zb6fM1CX41D9rF9hdgclu0peUmy.jpg",
			Genres:       []string{"Drama", "War"},
			Companies:    []string{"Universal Pictures"},
			SpokenLanguages: []string{"en", "de", "hi"},
			Popularity: 85.0, VoteAverage: 8.9, VoteCount: 14900,
			CastMembers: []string{"Tom Hanks", "Morgan Freeman"},
			Director: "Steven Spielberg", Writer: "Aaron Sorkin", Composer: "John Williams",
			HasTranslations: false,
		},
		{
			Title: "The Grand Budapest Hotel", OriginalTitle: "The Grand Budapest Hotel", Lang: "en",
			TagLine:  "A perfect holiday without the family.",
			Overview: "The Grand Budapest Hotel recounts the adventures of Gustave H, a legendary concierge at a famous European hotel between the wars, and Zero Moustafa, the lobby boy who becomes his most trusted friend.",
			Release:  "2014-02-26", Runtime: 99, Adult: false,
			Budget: 25000000, Revenue: 174600000,
			Homepage:     "https://www.foxsearchlight.com/thegrandbudapesthotel",
			PosterPath:   "/eWdyYQreja6JGCzqHWXpWHDrrPo.jpg",
			BackdropPath: "/nX5XotM9yprCKarRH4fzOq1VM1J.jpg",
			Genres:       []string{"Comedy", "Drama", "Adventure"},
			Companies:    []string{"20th Century Studios", "BBC Films"},
			SpokenLanguages: []string{"en", "fr", "de"},
			Popularity: 79.5, VoteAverage: 8.1, VoteCount: 16200,
			CastMembers: []string{"Brad Pitt", "Meryl Streep", "Scarlett Johansson"},
			Director: "Quentin Tarantino", Writer: "Quentin Tarantino", Composer: "Hans Zimmer",
			HasTranslations: true,
		},
	}

	movieIDs := make([]uuid.UUID, 0, len(movies))

	for _, m := range movies {
		id := newUUID()
		movieIDs = append(movieIDs, id)

		// Insert movie
		mustExec(ctx, tx, `
			INSERT INTO movies (
				id, title, original_title, original_language, tag_line, overview, 
				release_date, runtime, adult, homepage, poster_path, backdrop_path,
				budget, revenue, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $15)
		`, id, m.Title, m.OriginalTitle, m.Lang, m.TagLine, m.Overview,
			m.Release, m.Runtime, m.Adult, m.Homepage, m.PosterPath, m.BackdropPath,
			m.Budget, m.Revenue, now)

		// Insert movie_stats
		mustExec(ctx, tx, `INSERT INTO movie_stats (movie_id, popularity, vote_average, vote_count) VALUES ($1, $2, $3, $4)`,
			id, m.Popularity, m.VoteAverage, m.VoteCount)

		// Insert movie_genres
		for _, genreName := range m.Genres {
			if gid, ok := genreIDs[genreName]; ok {
				mustExec(ctx, tx, `INSERT INTO movie_genres (movie_id, genre_id) VALUES ($1, $2)`, id, gid)
			}
		}

		// Insert movie_companies
		for _, companyName := range m.Companies {
			if cid, ok := companyIDs[companyName]; ok {
				mustExec(ctx, tx, `INSERT INTO movie_companies (movie_id, company_id) VALUES ($1, $2)`, id, cid)
			}
		}

		// Insert movie_spoken_languages
		for _, lang := range m.SpokenLanguages {
			mustExec(ctx, tx, `INSERT INTO movie_spoken_languages (movie_id, iso_639_1) VALUES ($1, $2)`, id, lang)
		}

		// Insert credits - cast
		for i, castName := range m.CastMembers {
			if pid, ok := personIDs[castName]; ok {
				credID := newUUID()
				castOrder := i + 1
				charName := "Character " + castName
				mustExec(ctx, tx, `
					INSERT INTO credits (id, movie_id, person_id, credit_type, character_name, cast_order)
					VALUES ($1, $2, $3, $4, $5, $6)
				`, credID, id, pid, "cast", charName, castOrder)
			}
		}

		// Insert credits - director
		if m.Director != "" {
			if pid, ok := personIDs[m.Director]; ok {
				credID := newUUID()
				mustExec(ctx, tx, `
					INSERT INTO credits (id, movie_id, person_id, credit_type, department, job)
					VALUES ($1, $2, $3, $4, $5, $6)
				`, credID, id, pid, "crew", "Directing", "Director")
			}
		}

		// Insert credits - writer
		if m.Writer != "" && m.Writer != m.Director {
			if pid, ok := personIDs[m.Writer]; ok {
				credID := newUUID()
				mustExec(ctx, tx, `
					INSERT INTO credits (id, movie_id, person_id, credit_type, department, job)
					VALUES ($1, $2, $3, $4, $5, $6)
				`, credID, id, pid, "crew", "Writing", "Screenplay")
			}
		}

		// Insert credits - composer
		if m.Composer != "" {
			if pid, ok := personIDs[m.Composer]; ok {
				credID := newUUID()
				mustExec(ctx, tx, `
					INSERT INTO credits (id, movie_id, person_id, credit_type, department, job)
					VALUES ($1, $2, $3, $4, $5, $6)
				`, credID, id, pid, "crew", "Sound", "Original Music Composer")
			}
		}

		// Insert images - poster
		posterID := newUUID()
		mustExec(ctx, tx, `
			INSERT INTO images (id, movie_id, file_path, type, width, height, language)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, posterID, id, m.PosterPath, "poster", 500, 750, m.Lang)

		// Insert images - backdrop
		backdropID := newUUID()
		mustExec(ctx, tx, `
			INSERT INTO images (id, movie_id, file_path, type, width, height, language)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, backdropID, id, m.BackdropPath, "backdrop", 1920, 1080, nil)

		// Insert translations for movies that have them
		if m.HasTranslations {
			// Japanese translation
			mustExec(ctx, tx, `
				INSERT INTO movie_translations (movie_id, language, title, overview)
				VALUES ($1, $2, $3, $4)
			`, id, "ja", m.Title+" (日本語)", m.Overview+" (日本語翻訳)")

			// Spanish translation
			mustExec(ctx, tx, `
				INSERT INTO movie_translations (movie_id, language, title, overview)
				VALUES ($1, $2, $3, $4)
			`, id, "es", m.Title+" (Español)", m.Overview+" (Traducción al español)")

			// French translation
			mustExec(ctx, tx, `
				INSERT INTO movie_translations (movie_id, language, title, overview)
				VALUES ($1, $2, $3, $4)
			`, id, "fr", m.Title+" (Français)", m.Overview+" (Traduction française)")
		}
	}
	log.Printf("  - Movies seeded (%d movies)", len(movies))

	// ====================
	// ADD STILL IMAGES FOR SOME MOVIES
	// ====================
	// Add extra still images for the first 10 movies
	for i := 0; i < 10 && i < len(movieIDs); i++ {
		for j := 1; j <= 3; j++ {
			stillID := newUUID()
			stillPath := "/stills/" + movies[i].Title + "_still_" + string(rune('0'+j)) + ".jpg"
			mustExec(ctx, tx, `
				INSERT INTO images (id, movie_id, file_path, type, width, height, language)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, stillID, movieIDs[i], stillPath, "still", 1280, 720, nil)
		}
	}
	log.Println("  - Still images seeded")

	// ====================
	// ADD CINEMATOGRAPHER CREDITS
	// ====================
	cinematographers := []struct {
		MovieIndex int
		Person     string
	}{
		{0, "Roger Deakins"},   // Inception
		{2, "Roger Deakins"},   // Interstellar
		{11, "Roger Deakins"},  // Dune
		{12, "Roger Deakins"},  // Blade Runner 2049
		{6, "Emmanuel Lubezki"}, // The Matrix
	}
	for _, c := range cinematographers {
		if c.MovieIndex < len(movieIDs) {
			if pid, ok := personIDs[c.Person]; ok {
				credID := newUUID()
				mustExec(ctx, tx, `
					INSERT INTO credits (id, movie_id, person_id, credit_type, department, job)
					VALUES ($1, $2, $3, $4, $5, $6)
				`, credID, movieIDs[c.MovieIndex], pid, "crew", "Camera", "Director of Photography")
			}
		}
	}
	log.Println("  - Cinematographer credits seeded")

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("tx commit: %v", err)
	}

	log.Println("Database seeded successfully!")
	log.Printf("Summary:")
	log.Printf("  - %d languages", len(languages))
	log.Printf("  - %d genres", len(genreNames))
	log.Printf("  - %d companies", len(companyData))
	log.Printf("  - %d people", len(peopleData))
	log.Printf("  - %d movies", len(movies))
}
