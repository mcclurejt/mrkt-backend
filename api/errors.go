package api

// AlphaVantageRateExceededError - Error thrown when the number of calls to alphavantage exceeds 5 per minute
type AlphaVantageRateExceededError struct {
}

func (e *AlphaVantageRateExceededError) Error() string {
	return "AlphaVantage api rate exceeded, try again in a minute"
}
