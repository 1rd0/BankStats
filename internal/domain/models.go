package domain

import "time"

type RatePoint struct {
	Date     time.Time `json:"date"`
	CharCode string    `json:"char_code"`
	Name     string    `json:"name"`
	PerUnit  float64   `json:"per_unit"`
}

type Extremum struct {
	Value    float64   `json:"value"`
	CharCode string    `json:"char_code"`
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
}

type Stats struct {
	Max     Extremum `json:"max"`
	Min     Extremum `json:"min"`
	Average float64  `json:"average"`
	Count   int      `json:"count"`
	Days    int      `json:"days"`
}
