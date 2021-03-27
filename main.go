package main

import (
	"github.com/bonedaddy/go-compound/v2/client"
	"github.com/bonedaddy/go-compound/v2/models"
	"github.com/davidiola/compound_go_monitor/dataaccess"
	"github.com/davidiola/compound_go_monitor/orch"
)

var COMP_V2_URL string = "https://api.compound.finance/api/v2"

func main() {
	var accounts []models.Account

	cl := client.NewClient(COMP_V2_URL)
	o := orch.NewOrch(*cl)

	accounts = o.RetrieveAllAccounts()

	dataAccess := dataaccess.NewDataAccess()

	for _, acct := range accounts {
		dataAccess.WriteAccount(acct)
	}
}
