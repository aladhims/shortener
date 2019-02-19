DROP TABLE IF EXISTS users;
CREATE TABLE users (
id serial NOT NULL,
email character varying(100) NOT NULL UNIQUE,
fullname character varying(100) NOT NULL,
CONSTRAINT userinfo_pkey PRIMARY KEY (id)
)
WITH (OIDS=FALSE);