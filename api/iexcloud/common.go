package iexcloud

// common - Models shared across the iexcloud API

type BidAsk struct {
	Price     float64 `json:"price"`
	Size      int     `json:"size"`
	Timestamp int     `json:"timestamp"`
}

// Quote models the data returned from the IEX Cloud /quote endpoint.
type Quote struct {
	Symbol                 string  `json:"symbol"`
	CompanyName            string  `json:"companyName"`
	Exchange               string  `json:"primaryExchange"`
	CalculationPrice       string  `json:"calculationPrice"`
	Open                   float64 `json:"open"`
	OpenTime               int     `json:"openTime"`
	OpenSource             string  `json:"openSource"`
	Close                  float64 `json:"close"`
	CloseTime              int     `json:"closeTime"`
	CloseSource            string  `json:"closeSource"`
	High                   float64 `json:"high"`
	Low                    float64 `json:"low"`
	LowTime                float64 `json:"lowTime"`
	LowSource              string  `json:"lowSource"`
	LatestPrice            float64 `json:"latestPrice"`
	LatestSource           string  `json:"latestSource"`
	LatestTime             string  `json:"latestTime"`
	LatestUpdate           int     `json:"latestUpdate"`
	LatestVolume           int     `json:"latestVolume"`
	IexRealtimePrice       float64 `json:"iexRealtimePrice"`
	IexRealtimeSize        int     `json:"iexRealtimeSize"`
	IexLastUpdated         int     `json:"iexLastUpdated"`
	DelayedPrice           float64 `json:"delayedPrice"`
	DelayedPriceTime       int     `json:"delayedPriceTime"`
	OddLotDelayedPrice     int     `json:"oddLotDelayedPrice"`
	OddLotDelayedPriceTime int     `json:"oddLotDelayedPriceTime"`
	ExtendedPrice          float64 `json:"extendedPrice"`
	ExtendedChange         float64 `json:"extendedChange"`
	ExtendedChangePercent  float64 `json:"extendedChangePercent"`
	ExtendedPriceTime      int     `json:"extendedPriceTime"`
	PreviousClose          float64 `json:"previousClose"`
	Change                 float64 `json:"change"`
	ChangePercent          float64 `json:"changePercent"`
	Volume                 int     `json:"volume"`
	IexMarketPercent       float64 `json:"iexMarketPercent"`
	IexVolume              int     `json:"iexVolume"`
	AvgTotalVolume         int     `json:"avgTotalVolume"`
	IexBidPrice            float64 `json:"iexBidPrice"`
	IexBidSize             int     `json:"iexBidSize"`
	IexAskPrice            float64 `json:"iexAskPrice"`
	IexAskSize             int     `json:"iexAskSize"`
	IexOpen                float64 `json:"iexOpen"`
	IexOpenTime            int     `json:"iexOpenTime"`
	IexClose               float64 `json:"iexClose"`
	IexCloseTime           int     `json:"iexCloseTime"`
	MarketCap              int     `json:"marketCap"`
	PERatio                float64 `json:"peRatio"`
	Week52High             float64 `json:"week52High"`
	Week52Low              float64 `json:"week52Low"`
	YTDChange              float64 `json:"ytdChange"`
	LastTradeTime          int     `json:"lastTradeTime"`
}

// Trade models a trade for a quote.
type Trade struct {
	Price                 float64 `json:"price"`
	Size                  int     `json:"size"`
	TradeID               int     `json:"tradeId"`
	IsISO                 bool    `json:"isISO"`
	IsOddLot              bool    `json:"isOddLot"`
	IsOutsideRegularHours bool    `json:"isOutsideRegularHours"`
	IsSinglePriceCross    bool    `json:"isSinglePriceCross"`
	IsTradeThroughExempt  bool    `json:"isTradeThroughExempt"`
	Timestamp             int     `json:"timestamp"`
}
