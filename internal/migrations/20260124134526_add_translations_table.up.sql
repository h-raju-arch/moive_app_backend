CREATE TABLE movie_translations (
  movie_id   UUID REFERENCES movies(id) ON DELETE CASCADE,
  language   VARCHAR(10) NOT NULL,  -- e.g. en-US
  title      TEXT,
  overview   TEXT,
  PRIMARY KEY (movie_id, language)
);