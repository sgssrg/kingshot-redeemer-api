-- name: GetAllPlayers :many
SELECT
    *
FROM
    Players
ORDER BY
    PiD;

-- name: DeletePlayer :exec
DELETE FROM Players
WHERE
    PiD = $1;

-- name: PushPlayer :one
INSERT INTO
    Players (PiD, KiD, dName, PFP, Alliance)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdatePlayer :one
UPDATE Players
SET
    KiD = $2,
    dName = $3,
    PFP = $4,
    Alliance = $5
WHERE PiD = $1
RETURNING *;


-- name: PushGC :one
INSERT INTO
    Giftcode (code)
VALUES
    ($1)
RETURNING *;