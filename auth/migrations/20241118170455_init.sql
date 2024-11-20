-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id serial PRIMARY KEY,
  login varchar(100),
  passwordHash varchar(100)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
