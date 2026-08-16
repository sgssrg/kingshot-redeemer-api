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
	WrongField string `json:"wrong_field" example:"pid"`
	Message    string `json:"message" example:"Invalid PiD & It should be Interger > 0"`
}
type UpdatePlayerSSEValidator struct {
	WrongField string `json:"wrong_field" example:"code"`
	Message    string `json:"message" example:"Please enter a valid code."`
}
type DeletePlayerValidationError struct {
	WrongField string `json:"wrong_field" example:"pid"`
	Message    string `json:"message" example:"invalid player ID"`
	Value      string `json:"value"`
}

type PlayerInfo struct {
	Pid      uint   `json:"pid" example:"123456"`
	Kid      uint   `json:"kid" example:"123"`
	Dname    string `json:"dName,omitempty" example:"MeowMeow"`
	Pfp      string `json:"pfp,omitempty" example:"<-pfp-url->"`
	Alliance string `json:"alliance,omitempty" example:"XYZ"`
}
type DeletePlayerResponseSuccess struct {
	Pid     string     `json:"pid" example:"123"`
	Deleted bool       `json:"deleted" example:"true"`
	Player  PlayerInfo `json:"player_info"`
	Message string     `json:"message" example:"Player deleted successfully"`
}
type DeletePlayerResponseFailure struct {
	Pid     string     `json:"pid" example:"123"`
	Deleted bool       `json:"deleted" example:"false"`
	Player  PlayerInfo `json:"player_info,omitempty"`
	Message string     `json:"message" example:"Player not found"`
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
