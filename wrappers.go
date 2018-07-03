package main

import (
	"github.com/davidsteinsland/ynab-go/ynab"
)

type ClientInt interface {
	ListBudgets() ([]ynab.BudgetSummary, error)
	GetBudget(id string) (ynab.BudgetDetail, error)
	ListCategories(budgetID string) ([]ynab.CategoryGroupWithCategories, error)
}

type ClientWrapper struct {
	Client *ynab.Client
}

func (c *ClientWrapper) ListBudgets() ([]ynab.BudgetSummary, error) {
	return c.Client.BudgetService.List()
}

func (c *ClientWrapper) GetBudget(id string) (ynab.BudgetDetail, error) {
	return c.Client.BudgetService.Get(id)
}

func (c *ClientWrapper) ListCategories(budgetID string) ([]ynab.CategoryGroupWithCategories, error) {
	return c.Client.CategoriesService.List(budgetID)
}
