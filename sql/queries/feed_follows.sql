-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)

SELECT iff.*, u.name AS user_name, f.name AS feed_name
FROM inserted_feed_follow iff
INNER JOIN users u ON iff.user_id = u.id
INNER JOIN feeds f ON iff.feed_id = f.id;

-- name: GetFeedFollowsForUser :many
SELECT feeds.name AS feed_name, users.name AS creator_name
FROM feed_follows
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
INNER JOIN users
ON feeds.user_id = users.id
WHERE feed_follows.user_id = $1;

-- name: Unfollow :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1
AND feed_follows.feed_id = $2;