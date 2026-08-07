package model

type RedeemResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	ErrCode int    `json:"err_code"`
	Data    []any  `json:"data"`
}

type RedeemSSE struct {
	Fid     string         `json:"fid"`
	Code    string         `json:"code"`
	Result  RedeemResponse `json:"result"`
	Success bool           `json:"success"`
	Time    string         `json:"time"`
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
	UnkRetry    int `json:"unk_retry"`
	CodeExpired int `json:"code_expired"`
	PlayerDead  int `json:"PlayerDead"`
}
