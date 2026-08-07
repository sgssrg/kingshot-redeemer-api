package player

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/ribonin/apis/kingshot-redeem/db"
)

type App struct {
	Queries *db.Queries
	Pool    *pgxpool.Pool
}

type PlayerReqBody struct {
	Pid      string `json:"pid" validate:"required"`
	Kid      int    `json:"kid" validate:"required,gte=1"`
	Alliance string `json:"alliance" `
}
