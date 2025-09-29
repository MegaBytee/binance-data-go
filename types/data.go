package types

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"
)

const (
	BUCKET_URL = "https://s3-ap-northeast-1.amazonaws.com/data.binance.vision"
	baseUrl    = "https://data.binance.vision"
	Spot       = "spot"
	Futures    = "futures"
	Daily      = "daily"
	Monthly    = "monthly"
	Klines     = "klines"
)

type DataParams struct {
	From      string //spot or futures
	DateTime  string // daily or monthly
	Symbol    string
	TimeFrame string //1s, 1m, 5min ... (time frame of klines)
}

func NewDataParamsFromCli(from, dt, symbol, tf string) *DataParams {

	return &DataParams{
		From:      flagFrom(from),
		DateTime:  flagDateTime(dt),
		Symbol:    symbol,
		TimeFrame: tf,
	}

}

// validate and sanitize inputs minimally, fallback to defaults if empty
func (dp *DataParams) normalize() {
	if dp.From == "" {
		dp.From = Spot
	}
	if dp.DateTime == "" {
		dp.DateTime = Daily
	}
	if dp.TimeFrame == "" {
		dp.TimeFrame = "1m"
	}
	// symbol usually uppercase with no spaces; keep as supplied but trim spaces
	dp.Symbol = strings.TrimSpace(dp.Symbol)
}

// GetPathUrl builds a full URL according to the pattern:
// https://data.binance.vision/?prefix=data/{from}/{datetime}/klines/{symbol}/{timeframe}/
func (dp *DataParams) GetPathUrl() string {
	prefix := dp.GetPrefix()
	// build final URL with query parameter ?prefix=...
	return fmt.Sprintf("%s/?prefix=%s", baseUrl, prefix)
}

func (dp *DataParams) GetBucketUrl() string {

	prefix := dp.GetPrefix()
	log.Println(prefix)
	u, err := url.Parse(BUCKET_URL)
	if err != nil {
		log.Printf("something went wrong:in:GetBucketUrl  %v", err)
		return ""
	}
	// escape everything except '/'
	encodedPrefix := strings.ReplaceAll(url.QueryEscape(prefix), "%2F", "/")
	u.RawQuery = "delimiter=/&prefix=" + encodedPrefix
	return u.String()
}

func (dp *DataParams) GetPrefix() string {
	dp.normalize()

	// ensure klines segment is always used
	segments := []string{
		"data",
		dp.From,
		dp.DateTime,
		Klines,
		dp.Symbol,
		dp.TimeFrame,
		"",
	}

	// path.Join would clean slashes but also remove trailing slash; we want the trailing slash before query parameter
	prefix := path.Join(segments...)

	// build final URL with query parameter prefix
	return fmt.Sprint(prefix + "/")
}
