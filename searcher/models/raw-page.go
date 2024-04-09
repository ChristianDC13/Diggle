package models

type RawPage struct {
	PageId             int64   `json:"pageId"`
	HitsCount          int     `json:"hitsCount"`
	HitsSum            int     `json:"hitsSum"`
	Rank               float64 `json:"rank"`
	PositionalDistance int     `json:"positionalDistance"`
}
