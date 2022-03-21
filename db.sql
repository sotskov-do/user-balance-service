CREATE TABLE IF NOT EXISTS users
(
    id integer PRIMARY KEY NOT NULL,
    balance integer
);

CREATE TABLE IF NOT EXISTS history
(
    user_id integer NOT NULL,
    type character(20),
    amount integer,
    datetime timestamp without time zone,
    id serial PRIMARY KEY NOT NULL,
    CONSTRAINT user_id FOREIGN KEY (user_id)
        REFERENCES users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);