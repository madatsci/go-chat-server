-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email character varying(128) NOT NULL UNIQUE,
    password character varying(64) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE chat_messages (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id),
    receiver_id integer REFERENCES users(id),
    text text NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE INDEX chat_messages_user_id_receiver_id_idx ON chat_messages(user_id int4_ops, receiver_id int4_ops);
CREATE INDEX chat_messages_receiver_id_idx ON chat_messages(receiver_id int4_ops);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE users;
DROP TABLE chat_messages;
