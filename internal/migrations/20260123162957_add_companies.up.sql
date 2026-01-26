CREATE TABLE companies (
  id UUID PRIMARY KEY ,
  name TEXT NOT NULL,
  origin_country VARCHAR(10),
  homepage TEXT
);