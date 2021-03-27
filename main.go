package main

import (
	"fmt"
	"github.com/bonedaddy/go-compound/v2/client"
	"github.com/bonedaddy/go-compound/v2/models"
	"github.com/davidiola/compound_go_monitor/orch"
	"sync"
)

var COMP_V2_URL string = "https://api.compound.finance/api/v2"

func main() {
	var accounts []models.Account
	var wg sync.WaitGroup

	cl := client.NewClient(COMP_V2_URL)
	o := orch.NewOrch(*cl)
	numPages := o.RetrieveNumPages()
	wg.Add(numPages)
	acctChan := make(chan []models.Account, numPages)

	o.RetrieveAllAccounts(acctChan, &wg)
	wg.Wait()
	close(acctChan)
	i := 1
	for i <= numPages {
		accounts = append(accounts, <-acctChan...)
		i += 1
	}

	for _, acct := range accounts {
		printAddress(acct)
	}
}

func printAddress(acct models.Account) {
	fmt.Printf("Address: %s\n", acct.Address)
}
