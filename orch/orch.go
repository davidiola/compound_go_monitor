package orch

import (
	"fmt"
	"github.com/bonedaddy/go-compound/v2/client"
	"github.com/bonedaddy/go-compound/v2/models"
	"strconv"
	"sync"
)

var DEF_PAGE_SIZE string = "100"
var MAX_GO_ROUTINES int = 40

type Orch struct {
	c client.Client
}

func NewOrch(c client.Client) *Orch {
	return &Orch{c: c}
}

func (o *Orch) RetrieveAllAccounts(ch chan []models.Account, wg *sync.WaitGroup) {
	var pgNum = 1
	numPgs := o.RetrieveNumPages()
	limitRoutines := make(chan int, MAX_GO_ROUTINES)
	for pgNum <= numPgs {
		limitRoutines <- 1
		go o.RetrieveAccountsFromPage(ch, pgNum, wg, limitRoutines)
		fmt.Println(pgNum)
		pgNum += 1
	}
}

func (o *Orch) RetrieveAccountsFromPage(ch chan []models.Account, pgNum int, wg *sync.WaitGroup, limitRoutines chan int) {
	defer wg.Done()
	params := buildParams(pgNum)
	acctResp, err := o.c.GetAccount("", params)
	if err != nil {
		fmt.Printf("error is: %s", err)
	}
	if acctResp != nil {
		ch <- acctResp.Accounts
	}
	<-limitRoutines
}

func (o *Orch) RetrieveNumPages() int {
	// make dummy req to get total pages
	params := buildParams(1)
	acctResp, _ := o.c.GetAccount("", params)
	return acctResp.PaginationSummary.TotalPages
}

func buildParams(pgNum int) map[string]string {
	params := make(map[string]string)
	params["page_size"] = DEF_PAGE_SIZE
	params["page_number"] = strconv.Itoa(pgNum)
	return params
}
