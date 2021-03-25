package main

import (
	"fmt"
	"github.com/bonedaddy/go-compound/v2/client"
	"github.com/bonedaddy/go-compound/v2/models"
)

var COMP_V2_URL string = "https://api.compound.finance/api/v2"
var ADDRESS string = ""

func main() {
	cl := client.NewClient(COMP_V2_URL)
	resp, _ := cl.GetAccount(ADDRESS)
	for _, acct := range resp.Accounts {
		printAddress(acct)
	}
}

func printAddress(acct models.Account) {
	fmt.Printf("Address: %s\n", acct.Address)
}
