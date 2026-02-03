CREATE TABLE movie_spoken_languages (
  movie_id UUID REFERENCES movies(id) ON DELETE CASCADE,
  iso_639_1 VARCHAR(10) REFERENCES languages(iso_639_1),
  PRIMARY KEY (movie_id, iso_639_1)
);