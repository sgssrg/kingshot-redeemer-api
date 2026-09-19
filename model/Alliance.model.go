package model

type GetAllUniqueAllianceFromDBResponse struct {
	ErrorState   bool     `json:"error_state,omitempty" example:"false"`
	Message      string   `json:"message" example:"Fetched Unique Alliances List From DB. || No Alliances or Player in DB."`
	AllianceList []string `json:"alliance_list,omitempty" example:"[\"DES\", \"JPA\", \"HAW\", \"JPA\"]"`
}
