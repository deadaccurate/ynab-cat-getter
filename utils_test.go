package main

import (
	"errors"
	"testing"

	"github.com/davidsteinsland/ynab-go/ynab"
)

type FakeClient struct {
}

func (f *FakeClient) ListCategories(budgetId string) ([]ynab.CategoryGroupWithCategories, error) {
	return nil, errors.New("blah")
}

func (f *FakeClient) ListBudgets() ([]ynab.BudgetSummary, error) {
	return nil, errors.New("error budgets")
}

func (f *FakeClient) GetBudget(id string) (ynab.BudgetDetail, error) {
	return ynab.BudgetDetail{}, errors.New("error budget")
}

func TestGetCategory(t *testing.T) {
	client := &FakeClient{}

	c, err := GetCategory(client, "Blah")
	if err != nil {
		t.Errorf("Error while getting category: %s", err)
	}
	if c != "fake-id" {
		t.Errorf("Category didn't equal fake one: %s", c)
	}
}
