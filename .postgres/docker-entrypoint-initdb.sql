CREATE USER pizza_user
  WITH ENCRYPTED PASSWORD 'pizza_pass';

CREATE DATABASE pizza_db;

GRANT ALL PRIVILEGES ON DATABASE pizza_db TO pizza_user;

ALTER DATABASE pizza_db OWNER TO pizza_user;