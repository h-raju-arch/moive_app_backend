CREATE TABLE movie_companies (
  movie_id UUID REFERENCES movies(id) ON DELETE CASCADE,
  company_id UUID REFERENCES companies(id),
  PRIMARY KEY (movie_id, company_id)
);
