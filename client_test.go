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
	listPayErr bool
	retPay     []ynab.Payee
	transErr   bool
	retTrans   []ynab.HybridTransaction
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

func (f *FakeClient) ListPayees(budgetID string) ([]ynab.Payee, error) {
	if f.listPayErr {
		return nil, errors.New("error payees")
	}
	return f.retPay, nil
}

func (f *FakeClient) GetTransByCat(budgetID string, categoryID string) ([]ynab.HybridTransaction, error) {
	if f.transErr {
		return nil, errors.New("error transactions")
	}
	return f.retTrans, nil
}

func createBudgets() []ynab.BudgetSummary {
	return []ynab.BudgetSummary{
		ynab.BudgetSummary{
			Name: "my-budget",
			Id:   "b-id",
		},
	}
}

func createCategories() []ynab.CategoryGroupWithCategories {
	return []ynab.CategoryGroupWithCategories{
		ynab.CategoryGroupWithCategories{
			CategoryGroup: ynab.CategoryGroup{
				Name: "my-category",
				Id:   "c-id",
			},
			Categories: []ynab.Category{
				ynab.Category{
					Id:   "internal-id",
					Name: "internal",
				},
			},
		},
		ynab.CategoryGroupWithCategories{
			CategoryGroup: ynab.CategoryGroup{
				Name: "my-category2",
				Id:   "c-id2",
			},
			Categories: []ynab.Category{
				ynab.Category{
					Id:   "internal-id2",
					Name: "internal2",
				},
			},
		},
	}
}

func createPayees() []ynab.Payee {
	return []ynab.Payee{
		ynab.Payee{
			Id:   "pay-1",
			Name: "my-payee",
		},
	}
}

func createTransactions() []ynab.HybridTransaction {
	return createTransName("my-payee")
}

func createTransName(payee string) []ynab.HybridTransaction {
	return []ynab.HybridTransaction{
		ynab.HybridTransaction{
			TransactionSummary: ynab.TransactionSummary{
				Amount: -10000,
			},
			PayeeName: payee,
		},
	}
}

func TestFindBudget(t *testing.T) {
	c := &FakeClient{
		retBud: createBudgets(),
	}
	id, err := FindBudget(c, "my-budget")

	if err != nil {
		t.Error("FindBudget returned an error")
	}

	if id != "b-id" {
		t.Errorf("FindBudget returned: %s", id)
	}
}

func TestFindBudgetErr(t *testing.T) {
	c := &FakeClient{}
	_, err := FindBudget(c, "blah")
	if err == nil || !strings.Contains(err.Error(), "Unable to find the budget") {
		t.Error("Expected budget error")
	}

	c = &FakeClient{listBudErr: true}
	_, err = FindBudget(c, "blah")
	if err == nil || err.Error() != "error budgets" {
		t.Error("Expected list budget error")
	}
}

func TestGetCategory(t *testing.T) {
	client := &FakeClient{
		listBudErr: false,
		listCatErr: false,
		retBud:     createBudgets(),
		retCat:     createCategories(),
	}

	c, err := FindCategoryGroup(client, "my-budget", "my-category")
	if err != nil {
		t.Errorf("Error while getting category: %s", err)
	}
	if c.Id != "c-id" {
		t.Errorf("Category didn't equal fake one: %s", c.Id)
	}
}

func TestGetCategoryErr(t *testing.T) {
	cl := &FakeClient{
		listCatErr: true,
		retBud:     createBudgets(),
	}
	_, err := FindCategoryGroup(cl, "my-budget", "blah")

	if err == nil || err.Error() != "error categories" {
		t.Errorf("Expected error category: %s", err.Error())
	}
}

func TestGetCategoryCantFind(t *testing.T) {
	cl := &FakeClient{
		retBud: createBudgets(),
	}

	_, err := FindCategoryGroup(cl, "my-budget", "blah")
	if err == nil || !strings.Contains(err.Error(), "Unable to find the category") {
		t.Errorf("Expected category error: %s", err.Error())
	}
}

func TestGetPayees(t *testing.T) {
	c := &FakeClient{listPayErr: true}
	_, err := FindPayees(c, "blah")
	if err == nil || err.Error() != "error payees" {
		t.Error("Expected an error for FindPayees")
	}

	c = &FakeClient{retPay: createPayees()}
	p, err := FindPayees(c, "my-budget")
	if err != nil {
		t.Error("FindPayees returned an error")
		t.FailNow()
	}
	if _, ok := p["my-payee"]; !ok {
		t.Error("Map did not contain payee")
	}

	if _, ok := p["blah"]; ok {
		t.Error("Should not have found payee in map")
	}
}

func TestSumPayees(t *testing.T) {
	c := &FakeClient{
		transErr: true,
	}
	cat := &createCategories()[0]
	_, err := SumPayees(c, "b-id", cat)
	if err == nil || err.Error() != "error transactions" {
		t.Error("Expected transaction error")
	}

	c = &FakeClient{
		retPay:   createPayees(),
		retTrans: createTransactions(),
	}

	trans, err := SumPayees(c, "b-id", cat)
	if err != nil {
		t.Errorf("Received error from SumPayees: %s", err.Error())
	}

	val, ok := trans["my-payee"]
	if !ok {
		t.Error("Didn't find payee in transaction map!")
	}
	if val != -10.00 {
		t.Errorf("Unexpected sum: %f", val)
	}
	c = &FakeClient{
		retPay: createPayees(),
		// This won't show up in the payees
		retTrans: createTransName("crazy-payee"),
	}

	trans, err = SumPayees(c, "b-id", cat)
	if err != nil {
		t.Errorf("Received error for crazy transaction: %s", err.Error())
	}
	if _, ok := trans["my-payee"]; ok {
		t.Error("Payee was found for crazy transaction test")
	}
}
