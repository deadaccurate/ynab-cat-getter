package main

import "errors"

func GetCategory(client ClientInt, budgetName string, catName string) (string, error) {
	budgets, err := client.ListBudgets()
	if err != nil {
		return "", err
	}
	var budgetID string

	for _, b := range budgets {
		if b.Name == budgetName {
			budgetID = b.Id
			break
		}
	}

	if budgetID == "" {
		return "", errors.New("Unable to find the budget: " + budgetName)
	}

}
