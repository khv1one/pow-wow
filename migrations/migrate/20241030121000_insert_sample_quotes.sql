-- +goose Up
INSERT INTO quotes (quote, author) VALUES
                                       ('The only limit to our realization of tomorrow is our doubts of today.', 'Franklin D. Roosevelt'),
                                       ('Do not wait to strike till the iron is hot; but make it hot by striking.', 'William Butler Yeats'),
                                       ('The greatest glory in living lies not in never falling, but in rising every time we fall.', 'Nelson Mandela');

-- +goose Down
DELETE FROM quotes WHERE author IN ('Franklin D. Roosevelt', 'William Butler Yeats', 'Nelson Mandela');