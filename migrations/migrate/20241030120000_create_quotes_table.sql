-- +goose Up
CREATE TABLE quotes (
                        id SERIAL PRIMARY KEY,
                        quote TEXT NOT NULL,
                        author VARCHAR(255)
);

-- +goose Down
DROP TABLE quotes;