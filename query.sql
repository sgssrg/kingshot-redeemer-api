-- name: GetAllPlayers :many
SELECT
    *
FROM
    Players
ORDER BY
    PiD;

-- name: DeletePlayer :one
DELETE FROM Players
WHERE PiD = ?
RETURNING *;


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

-- name: GetGC :one
SELECT code, claimedAt 
FROM Giftcode
WHERE code=?
LIMIT 1;

-- name: GetAllUniqueAlliance :many
SELECT DISTINCT Alliance 
FROM Players 
WHERE Alliance IS NOT NULL
AND Alliance != '';