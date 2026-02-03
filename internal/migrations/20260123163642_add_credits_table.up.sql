CREATE TABLE credits (
  id UUID PRIMARY KEY,
  movie_id UUID REFERENCES movies(id) ON DELETE CASCADE,
  person_id UUID REFERENCES people(id),
  credit_type VARCHAR(10) NOT NULL,  -- cast or crew
  department TEXT,
  job TEXT,
  character_name TEXT,
  cast_order INT,
  created_at TIMESTAMPTZ DEFAULT now()
);