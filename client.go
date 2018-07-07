package main

import (
	"errors"
	"strings"

	"github.com/davidsteinsland/ynab-go/ynab"
)

type ClientInt interface {
	ListBudgets() ([]ynab.BudgetSummary, error)
	ListCategories(budgetID string) ([]ynab.CategoryGroupWithCategories, error)
	ListPayees(budgetID string) ([]ynab.Payee, error)
	GetTransByCat(budgetID string, categoryID string) ([]ynab.HybridTransaction, error)
}

type ClientWrapper struct {
	Client *ynab.Client
}

func (c *ClientWrapper) ListBudgets() ([]ynab.BudgetSummary, error) {
	return c.Client.BudgetService.List()
}

func (c *ClientWrapper) ListCategories(budgetID string) ([]ynab.CategoryGroupWithCategories, error) {
	return c.Client.CategoriesService.List(budgetID)
}

func (c *ClientWrapper) ListPayees(budgetID string) ([]ynab.Payee, error) {
	return c.Client.PayeesService.List(budgetID)
}

func (c *ClientWrapper) GetTransByCat(budgetID string, categoryID string) ([]ynab.HybridTransaction, error) {
	return c.Client.TransactionsService.GetByCategory(budgetID, categoryID)
}

func FindBudget(client ClientInt, budget string) (string, error) {
	budgets, err := client.ListBudgets()
	if err != nil {
		return "", err
	}
	var budgetID string

	for _, b := range budgets {
		if b.Name == budget {
			budgetID = b.Id
			break
		}
	}

	if budgetID == "" {
		return "", errors.New("Unable to find the budget: " + budget)
	}
	return budgetID, nil
}

func FindCategoryGroup(client ClientInt,
	budgetID string,
	catName string) (*ynab.CategoryGroupWithCategories, error) {
	catWithGroup, err := client.ListCategories(budgetID)
	if err != nil {
		return nil, err
	}

	var group *ynab.CategoryGroupWithCategories

	for _, c := range catWithGroup {

		if strings.TrimRight(c.Name, " ") == catName {
			group = &c
			break
		}
	}

	if group == nil {
		return nil, errors.New("Unable to find the category: " + catName)
	}

	return group, nil
}

func FindPayees(c ClientInt, budgetID string) (map[string]float32, error) {
	p, err := c.ListPayees(budgetID)
	if err != nil {
		return nil, err
	}
	pMap := make(map[string]float32)
	for _, payee := range p {
		pMap[strings.TrimRight(payee.Name, " ")] = 0
	}
	return pMap, nil
}
