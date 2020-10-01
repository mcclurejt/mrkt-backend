package api

type AlphaVantageRateExceededError struct {
}

func (e *AlphaVantageRateExceededError) Error() string {
	return "AlphaVantage api rate exceeded, try again in a minute"
}
