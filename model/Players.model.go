package model

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/ribonin/apis/kingshot-redeem/db"
)

type PlayerResponse struct {
	Message string     `json:"message"`
	Type    uint       `json:"type"`
	Player  PlayerInfo `json:"player"`
}

type PlayerAPIResponse struct {
	Message string     `json:"message"`
	Type    uint       `json:"type"`
	Player  PlayerInfo `json:"player"`
}
type AddPlayerValidator struct {
	WrongField string `json:"wrong_field"`
	Message    string `json:"message"`
}
type PlayerInfo struct {
	Pid      uint   `json:"pid"`
	Kid      uint   `json:"kid"`
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
	Kingdom               uint   `json:"kingdom"`
	Level                 uint   `json:"level"`
	LevelRendered         string `json:"levelRendered"`
	LevelRenderedDetailed string `json:"LevelRenderedDetailed"`
	LevelImage            string `json:"levelImage"`
	ProfilePhoto          string `json:"profilePhoto"`
}

type PlayerReqBody struct {
	Pid      string `json:"pid" validate:"required"`
	Kid      *uint  `json:"kid"`
	Alliance string `json:"alliance"`
}
type GiftCodeReqBody struct {
	Pid      string `json:"pid" validate:"required"`
	Kid      *uint  `json:"kid"`
	Giftcode string `json:"giftcode" validate:"required"`
}
type DeletePlayerResponse struct {
	Pid     string `json:"pid"`
	Deleted bool   `json:"deleted"`
	Message string `json:"message"`
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

type UpdatePlayerSSE struct {
	Updated bool       `json:"updated"`
	Player  PlayerInfo `json:"player"`
}

type UpdateMetaSSE struct {
	Count uint `json:"count"`
}
