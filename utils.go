package main

import (
	"strings"

	"github.com/deadaccurate/ynab-go/ynab"
)

func SumPayees(c ClientInt,
	budgetID string,
	group *ynab.CategoryGroupWithCategories,
	sinceDate string) (map[string]float32, error) {

	payees := make(map[string]float32)
	for _, cat := range group.Categories {
		trans, err := c.GetTransByCat(budgetID, cat.Id, sinceDate)
		if err != nil {
			return nil, err
		}

		for _, t := range trans {
			val, _ := payees[strings.TrimRight(t.PayeeName, " ")]

			payees[t.PayeeName] = val + (float32(t.TransactionSummary.Amount) / 1000)
		}
	}

	return payees, nil
}
