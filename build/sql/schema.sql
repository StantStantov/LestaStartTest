CREATE SCHEMA IF NOT EXISTS lesta_start;
SET search_path TO lesta_start;

CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS collections (
  id TEXT PRIMARY KEY,
  user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
  collection TEXT NOT NULL,
  UNIQUE (user_id, collection)
);

CREATE TABLE IF NOT EXISTS documents (
  id TEXT PRIMARY KEY,
  user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
  document TEXT NOT NULL,
  UNIQUE (user_id, document)
);

CREATE TABLE IF NOT EXISTS collection_documents (
  collection_id TEXT REFERENCES collections (id) ON DELETE CASCADE,
  document_id TEXT REFERENCES documents (id) ON DELETE CASCADE,
  PRIMARY KEY (collection_id, document_id)
);

CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
  user_id TEXT NULL REFERENCES users(id) ON DELETE CASCADE,
	csrf_token TEXT NOT NULL UNIQUE,
	expire_on TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS metrics (
  timestamp TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  value DOUBLE PRECISION NOT NULL,
  PRIMARY KEY (timestamp, name)
);

CREATE INDEX sessions_expire_on_idx ON sessions (expire_on);
