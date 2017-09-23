-- create the UUID extension
-- this allows for default UUID values using uuid_generate
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- User information
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  username VARCHAR(128) NOT NULL
);

-- User passwords
--
-- Note this is not part of the users table to support different
-- auth schemes in the future where users may not have a password
CREATE TABLE passwords (
  user_id UUID REFERENCES users PRIMARY KEY,
  hash BYTEA NOT NULL
);