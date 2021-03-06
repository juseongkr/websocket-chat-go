CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL
);

CREATE TABLE chatrooms (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL REFERENCES users ON DELETE CASCADE,
    chatroom_id INT NOT NULL REFERENCES chatrooms ON DELETE CASCADE,
    text VARCHAR NOT NULL,
    sent_on TIMESTAMPTZ NOT NULL DEFAULT NOW()
);