CREATE TABLE movies (
  id UUID PRIMARY KEY,                                     
  title TEXT NOT NULL,
  original_title TEXT,
  original_language VARCHAR(10),
  tag_line TEXT,
  overview TEXT,
  release_date DATE,
  runtime INT,
  adult BOOLEAN DEFAULT FALSE,
  homepage TEXT,
  poster_path TEXT,
  backdrop_path TEXT,
  budget BIGINT DEFAULT 0,
  revenue BIGINT DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);