package model

type ScrapePlayerInfo struct {
	Pid      int    `json:"pid"`
	Kid      int    `json:"kid"`
	Dname    string `json:"dName,omitempty"`
	Pfp      string `json:"pfp,omitempty"`
	Alliance string `json:"alliance,omitempty"`
}
type CustomScrapePlayerErrInfo struct {
	Pid     int    `json:"pid"`
	Type    int    `json:"type"`
	Message string `json:"message"`
}
