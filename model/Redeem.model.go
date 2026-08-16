package model

type RedeemResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	ErrCode int    `json:"err_code"`
	Data    []any  `json:"data"`
}

type RedeemSSE struct {
	Fid           string                `json:"fid"`
	Code          string                `json:"code"`
	RawResult     RedeemResponse        `json:"raw_result"`
	RefinedResult RedeemRefinedResponse `json:"result"`
	Success       bool                  `json:"success"`
	Time          string                `json:"time"`
}
type RedeemRefinedResponse struct {
	Message          string `json:"message"`
	Redeemed         bool   `json:"redeemed"`
	TypeRedeemedCode uint   `json:"type_redeem_code"`
}

type RedeemBodyKS struct {
	Fid  int    `json:"fid"`
	Kid  int    `json:"kid"`
	Cdk  string `json:"cdk"`
	Time string `json:"time"`
}

type PlayerMetaRedeemResponse struct {
	Redeemed    int `json:"redeemed"`
	Manual      int `json:"manual"`
	CodeExpired int `json:"code_expired"`
	PlayerDead  int `json:"PlayerDead"`
}


