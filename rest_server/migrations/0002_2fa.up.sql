CREATE TABLE login_sessions (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  created_timestamp TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT now(),
  expiration_timestamp TIMESTAMP(0) WITH TIME ZONE NOT NULL,

  UNIQUE (user_id)
);

CREATE TABLE two_factor_setup_sessions (
  user_id INTEGER REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  secret_string TEXT NOT NULL,
  created_timestamp TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT now(),
  expiration_timestamp TIMESTAMP(0) WITH TIME ZONE NOT NULL,
  UNIQUE (user_id)
);

CREATE TABLE recovery_codes (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  is_redeemed BOOLEAN DEFAULT false,

  code TEXT NOT NULL
);
