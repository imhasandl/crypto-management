package handlers

type CoinRequest struct {
	Coin string `json:"coin"`
}

type CoinPriceRequest struct {
	Coin      string `json:"coin"`
	Timestamp int64  `json:"timestamp"`
}
