package main

type Config struct {
	Key       string `json:"key"`
	Budget    string `json:"budget"`
	Category  string `json:"category"`
	StartDate string `json:"start_date,omitempty"`
}
