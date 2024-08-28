-- +goose Up

CREATE TABLE posts
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    publish_date TIMESTAMP WITH TIME ZONE NOT NULL,
    publish_date_raw TEXT NOT NULL,
    description TEXT,
    img_url TEXT,
    categories TEXT,
    url TEXT UNIQUE NOT NULL,
    website TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'Europe/Athens')
);

CREATE INDEX idx_posts_website ON posts(website);


-- +goose Down

DROP TABLE posts;