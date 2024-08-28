-- name: CreatePost :one
INSERT  INTO posts
    (title , publish_date , publish_date_raw,description , img_url , categories , url , website )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;