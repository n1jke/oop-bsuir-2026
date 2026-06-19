BEGIN;

CREATE TABLE users (
  id            UUID PRIMARY KEY,
  login         TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  rating        DOUBLE PRECISION DEFAULT 0
);

CREATE TABLE books (
  id          UUID PRIMARY KEY,
  title       TEXT NOT NULL,
  author_id   UUID,
  isbn        TEXT UNIQUE NOT NULL,
  description TEXT DEFAULT '',
  topics      TEXT[] DEFAULT '{}'
);

CREATE TABLE owned_books (
  id       UUID PRIMARY KEY,
  book_id  UUID NOT NULL REFERENCES books(id),
  owner_id UUID NOT NULL REFERENCES users(id),
  status   INT DEFAULT 0,
  UNIQUE (owner_id, book_id)
);

CREATE TABLE exchanges (
  id            UUID PRIMARY KEY,
  owned_book_id UUID NOT NULL REFERENCES owned_books(id),
  from_id       UUID NOT NULL REFERENCES users(id),
  to_id         UUID NOT NULL REFERENCES users(id),
  status        TEXT NOT NULL DEFAULT 'pending',
  note          TEXT DEFAULT '',
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  updated_at    TIMESTAMPTZ DEFAULT NOW(),
  expires_at    TIMESTAMPTZ
);

CREATE TABLE book_reviews (
  id      UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  book_id UUID NOT NULL REFERENCES books(id),
  mark    INT NOT NULL,
  report  TEXT DEFAULT ''
);

CREATE TABLE user_reviews (
  id      UUID PRIMARY KEY,
  from_id UUID NOT NULL REFERENCES users(id),
  to_id   UUID NOT NULL REFERENCES users(id),
  mark    INT NOT NULL,
  report  TEXT DEFAULT ''
);

COMMIT;