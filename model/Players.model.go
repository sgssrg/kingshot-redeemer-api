package model

import (
	"gitlab.com/ribonin/apis/kingshot-redeem/db"
)

type PlayerResponse struct {
	Message string    `json:"message"`
	Type    int       `json:"type"`
	Player  db.Player `json:"player"`
}

type PlayerAPIResponse struct {
	Message string     `json:"message"`
	Type    int        `json:"type"`
	Player  PlayerInfo `json:"player"`
}

type PlayerInfo struct {
	Pid      int32  `json:"pid"`
	Kid      int32  `json:"kid"`
	Dname    string `json:"dname,omitempty"`
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
