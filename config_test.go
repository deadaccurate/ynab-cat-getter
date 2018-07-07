package main

import (
	"encoding/json"
	"testing"
)

var config = `
{
	"key": "FAKE-KEY",
	"budget": "MyBudget",
	"category": "RandomCat",
	"start_date": "2000-01-01",
	"end_date": "2001-01-01"
}
`

func TestConfig(t *testing.T) {
	var c Config
	b := []byte(config)
	if err := json.Unmarshal(b, &c); err != nil {
		t.Errorf("Unmarshall err: %s", err)
	}

	if c.Budget != "MyBudget" {
		t.Errorf("Budget field was wrong: %s", c.Budget)
	}

	if c.Category != "RandomCat" {
		t.Errorf("Category field was wrong: %s", c.Category)
	}

	if c.Key != "FAKE-KEY" {
		t.Errorf("Key field was wrong: %s", c.Key)
	}

	if c.StartDate != "2000-01-01" {
		t.Errorf("StartDate field was wrong: %s", c.StartDate)
	}
}
