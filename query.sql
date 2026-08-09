-- name: GetAllPlayers :many
SELECT
    *
FROM
    Players
ORDER BY
    PiD;

-- name: DeletePlayer :execrows
DELETE FROM Players
WHERE
    PiD = ?;

-- name: PushPlayer :one
INSERT INTO
    Players (PiD, KiD, dName, PFP, Alliance)
VALUES
    (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdatePlayer :one
UPDATE Players
SET
    KiD = ?,
    dName = ?,
    PFP = ?,
    Alliance = ?
WHERE PiD = ?
RETURNING *;


-- name: PushGC :one
INSERT INTO
    Giftcode (code)
VALUES
    (?)
RETURNING *;