CREATE TABLE movie_stats (
  movie_id UUID PRIMARY KEY REFERENCES movies(id) ON DELETE CASCADE,
  popularity DOUBLE PRECISION DEFAULT 0,
  vote_average DOUBLE PRECISION DEFAULT 0,
  vote_count INT DEFAULT 0,
  scores_updated_at TIMESTAMPTZ DEFAULT now()
);