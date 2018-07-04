package main

import "errors"

func GetCategoryID(client ClientInt, budgetName string, catName string) (string, error) {
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

	catWithGroup, err := client.ListCategories(budgetID)
	if err != nil {
		return "", err
	}

	var catID string
	for _, c := range catWithGroup {
		if c.Name == catName {
			catID = c.Id
		}
	}

	if catID == "" {
		return "", errors.New("Unable to find the category: " + catName)
	}

	return catID, nil
}
