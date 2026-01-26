CREATE TABLE people (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  profile_path TEXT,
  known_for TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);