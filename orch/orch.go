package orch

import (
	"github.com/bonedaddy/go-compound/v2/client"
	"github.com/bonedaddy/go-compound/v2/models"
	"log"
	"strconv"
	"sync"
)

var PAGE_SIZE string = "100"
var MAX_GO_ROUTINES int = 40

type Orch struct {
	c client.Client
}

func NewOrch(c client.Client) *Orch {
	return &Orch{c: c}
}

func (o *Orch) RetrieveAllAccounts(wg *sync.WaitGroup) []models.Account {
	var accounts []models.Account
	var pgNum = 1
	numPages := o.RetrieveNumPages()
	wg.Add(numPages)
	acctChan := make(chan []models.Account, numPages)
	limitRoutines := make(chan int, MAX_GO_ROUTINES)
	for pgNum <= numPages {
		limitRoutines <- 1
		go o.RetrieveAccountsFromPage(acctChan, pgNum, wg, limitRoutines)
		log.Println(pgNum)
		pgNum += 1
	}
	wg.Wait()
	close(acctChan)
	i := 1
	for i <= numPages {
		accounts = append(accounts, <-acctChan...)
		i += 1
	}
	return accounts
}

func (o *Orch) RetrieveAccountsFromPage(ch chan []models.Account, pgNum int, wg *sync.WaitGroup, limitRoutines chan int) {
	defer wg.Done()
	acctResp, err := o.c.GetAccounts(PAGE_SIZE, strconv.Itoa(pgNum))
	if err != nil {
		log.Fatalf("error is: %s", err)
	}
	if acctResp != nil {
		ch <- acctResp.Accounts
	}
	<-limitRoutines
}

func (o *Orch) RetrieveNumPages() int {
	// make dummy req to get total pages
	acctResp, _ := o.c.GetAccounts(PAGE_SIZE, strconv.Itoa(0))
	return acctResp.PaginationSummary.TotalPages
}
