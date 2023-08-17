CREATE TABLE person (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  country VARCHAR(255) NOT NULL,
  full_name VARCHAR(255) NOT NULL,
  date_of_birth DATE NULL,
  age INT NULL
);

INSERT INTO person (id, country, full_name, date_of_birth) VALUES
  ('a58933a1-c24f-43d9-bb53-6a1aa3170a12', 'jp', 'a', '1987-01-01'),
  ('b5653a6d-fa58-4a91-9328-3b8c162123d5', 'de', 'b', '1973-01-01'),
  ('c9661cf0-e0e8-4ddb-9bb7-bfcda1aec90f', 'fr', 'c', '2018-01-01'),
  ('d86506a2-e186-4b89-97be-3294cb86d53a', 'us', 'd', '2003-01-01');

CREATE OR REPLACE FUNCTION get_oldest_person() RETURNS VARCHAR
LANGUAGE SQL
AS $$
  SELECT full_name
  FROM person
  ORDER BY date_of_birth DESC
  LIMIT 1;
$$;
SELECT get_oldest_person();

CREATE OR REPLACE FUNCTION people_born_on(d DATE) RETURNS INT
LANGUAGE SQL
AS $$
  SELECT count(*)
  FROM person
  WHERE date_of_birth = d;
$$;
SELECT people_born_on('1987-01-01');

CREATE OR REPLACE FUNCTION people_between(id_from UUID, id_to UUID) RETURNS SETOF RECORD
LANGUAGE SQL
AS $$
  SELECT id, country, full_name, date_of_birth
  FROM person
  WHERE id BETWEEN id_from AND id_to;
$$;
SELECT people_between('a58933a1-c24f-43d9-bb53-6a1aa3170a12', 'c9661cf0-e0e8-4ddb-9bb7-bfcda1aec90f');

CREATE OR REPLACE FUNCTION add_person(full_name VARCHAR(255), date_of_birth DATE, country VARCHAR(255)) RETURNS VOID
LANGUAGE SQL
AS $$
  INSERT INTO person (full_name, date_of_birth, country) VALUES
    (full_name, date_of_birth, country);
$$;
SELECT add_person('new person', '2000-01-01', 'us');