package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deadaccurate/ynab-go/ynab"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <config file>\n", os.Args[0])
		os.Exit(1)
	}
	configData, err := ioutil.ReadFile(os.Args[1])
	checkErr(err)
	var config Config
	err = json.Unmarshal(configData, &config)
	checkErr(err)

	c := &ClientWrapper{ynab.NewDefaultClient(config.Key)}
	bID, err := FindBudget(c, config.Budget)
	checkErr(err)
	group, err := FindCategoryGroup(c, bID, config.Category)
	checkErr(err)
	payees, err := SumPayees(c, bID, group, config.StartDate)
	checkErr(err)
	for k, v := range payees {
		fmt.Printf("%s, total: $%.2f\n", k, v)
	}

}
