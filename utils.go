package main

import (
	"fmt"
	"strings"
)

func SumPayees(c ClientInt, budgetID string, categoryID string) (map[string]float32, error) {
	payees, err := FindPayees(c, budgetID)
	if err != nil {
		return nil, err
	}

	trans, err := c.GetTransByCat(budgetID, categoryID)
	if err != nil {
		return nil, err
	}

	for _, t := range trans {
		_, ok := payees[strings.TrimRight(t.PayeeName, " ")]
		if ok {
			payees[t.PayeeName] += (float32(t.TransactionSummary.Amount) / 1000)
		} else {
			fmt.Printf("Found a payee that wasn't in the map! [%s]", t.PayeeName)
		}
	}
	return payees, nil
}
