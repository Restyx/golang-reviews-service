CREATE TABLE IF NOT EXISTS reviews(
    id serial PRIMARY KEY,
    author VARCHAR (50) NOT NULL,
    rating smallint NOT NULL,
    title VARCHAR (50),
    description TEXT
)