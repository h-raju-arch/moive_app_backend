DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'image_type') THEN
    CREATE TYPE image_type AS ENUM ('poster', 'backdrop', 'still');
  END IF;
END$$;

CREATE TABLE images (
  id UUID PRIMARY KEY,
  movie_id UUID REFERENCES movies(id) ON DELETE CASCADE,
  file_path TEXT,
  type image_type,
  width INT,
  height INT,
  language VARCHAR(10),
  created_at TIMESTAMPTZ DEFAULT now()
);