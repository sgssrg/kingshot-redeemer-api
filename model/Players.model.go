package model

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/ribonin/apis/kingshot-redeem/db"
)

type PlayerResponse struct {
	Message string     `json:"message"`
	Type    int        `json:"type"`
	Player  PlayerInfo `json:"player"`
}

type PlayerAPIResponse struct {
	Message string     `json:"message"`
	Type    int        `json:"type"`
	Player  PlayerInfo `json:"player"`
}

type PlayerInfo struct {
	Pid      int32  `json:"pid"`
	Kid      int32  `json:"kid"`
	Dname    string `json:"dName,omitempty"`
	Pfp      string `json:"pfp,omitempty"`
	Alliance string `json:"alliance,omitempty"`
}

type KSFetchPlayerResponse struct {
	Status    string `json:"status"`
	Data      pData  `json:"data"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
type pData struct {
	PlayerId              string `json:"playerId"`
	Name                  string `json:"name"`
	Kingdom               int    `json:"kingdom"`
	Level                 int    `json:"level"`
	LevelRendered         string `json:"levelRendered"`
	LevelRenderedDetailed string `json:"LevelRenderedDetailed"`
	LevelImage            string `json:"levelImage"`
	ProfilePhoto          string `json:"profilePhoto"`
}

type PlayerReqBody struct {
	Pid      string `json:"pid" validate:"required"`
	Kid      int    `json:"kid" validate:"required,gte=1"`
	Alliance string `json:"alliance" `
}
type DeletePlayerResponse struct {
	Pid     string `json:"pid"`
	Deleted bool   `json:"deleted"`
}
type DeletePlayerErrorResponse struct {
	Pid     string `json:"pid"`
	Deleted bool   `json:"deleted"`
	Error   string `json:"error"`
}

type App struct {
	Queries *db.Queries
	Pool    *pgxpool.Pool
}
