package model

type ScrapePlayerInfo struct {
	Pid      int    `json:"pid" example:"123456"`
	Kid      int    `json:"kid" example:"123"`
	Dname    string `json:"dName,omitempty" example:"MeowMeow"`
	Pfp      string `json:"pfp,omitempty" example:"<-pfp-url->"`
	Alliance string `json:"alliance,omitempty" example:"XYZ"`
}

type CustomScrapePlayerErrInfo struct {
	Pid     int    `json:"pid"`
	Type    int    `json:"type"`
	Message string `json:"message"`
}
