package main

import (
	"github.com/bonedaddy/go-compound/v2/client"
	"github.com/bonedaddy/go-compound/v2/models"
	"github.com/davidiola/compound_go_monitor/dataaccess"
	"github.com/davidiola/compound_go_monitor/orch"
)

var COMP_V2_URL string = "https://api.compound.finance/api/v2"

func main() {
	cl := client.NewClient(COMP_V2_URL)
	da := dataaccess.NewDataAccess()

	o := orch.NewOrch(cl)
	acctChan := make(chan []models.Account)

	o.RetrieveAllAccounts(acctChan)
	da.WriteAllAccounts(acctChan)
}
