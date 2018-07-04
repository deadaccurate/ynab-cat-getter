package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/davidsteinsland/ynab-go/ynab"
)

type FakeClient struct {
	listBudErr bool
	retBud     []ynab.BudgetSummary
	listCatErr bool
	retCat     []ynab.CategoryGroupWithCategories
}

func (f *FakeClient) ListCategories(budgetID string) ([]ynab.CategoryGroupWithCategories, error) {
	if f.listCatErr {
		return nil, errors.New("error categories")
	}
	return f.retCat, nil
}

func (f *FakeClient) ListBudgets() ([]ynab.BudgetSummary, error) {
	if f.listBudErr {
		return nil, errors.New("error budgets")
	}
	return f.retBud, nil
}

func (f *FakeClient) GetBudget(id string) (ynab.BudgetDetail, error) {
	return ynab.BudgetDetail{}, errors.New("error budget")
}

func TestGetCategory(t *testing.T) {
	d := make([]ynab.BudgetSummary, 1)
	d[0] = ynab.BudgetSummary{Name: "my-budget", Id: "b-id"}

	cat := make([]ynab.CategoryGroupWithCategories, 1)
	cat[0] = ynab.CategoryGroupWithCategories{
		CategoryGroup: ynab.CategoryGroup{
			Name: "my-category",
			Id:   "c-id",
		},
	}

	client := &FakeClient{listBudErr: false, listCatErr: false, retBud: d, retCat: cat}

	c, err := GetCategoryID(client, "my-budget", "my-category")
	if err != nil {
		t.Errorf("Error while getting category: %s", err)
	}
	if c != "c-id" {
		t.Errorf("Category didn't equal fake one: %s", c)
	}
}

func TestGetCategoryErr(t *testing.T) {
	cl := &FakeClient{listBudErr: true}

	_, err := GetCategoryID(cl, "blah", "blah")
	if err == nil || err.Error() != "error budgets" {
		t.Error("Expected error budget")
	}
	cl = &FakeClient{
		listCatErr: true,
		retBud: []ynab.BudgetSummary{
			ynab.BudgetSummary{
				Id:   "b-id",
				Name: "my-budget",
			},
		},
	}
	_, err = GetCategoryID(cl, "my-budget", "blah")

	if err == nil || err.Error() != "error categories" {
		t.Errorf("Expected error category: %s", err.Error())
	}
}

func TestGetCategoryCantFind(t *testing.T) {
	cl := &FakeClient{}

	_, err := GetCategoryID(cl, "blah", "blah")
	if err == nil || !strings.Contains(err.Error(), "Unable to find the budget") {
		t.Error("Expected budget error")
	}

	cl = &FakeClient{
		retBud: []ynab.BudgetSummary{
			ynab.BudgetSummary{
				Name: "my-budget",
				Id:   "b-id",
			},
		},
	}

	_, err = GetCategoryID(cl, "my-budget", "blah")
	if err == nil || !strings.Contains(err.Error(), "Unable to find the category") {
		t.Errorf("Expected category error: %s", err.Error())
	}
}
