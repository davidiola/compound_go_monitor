package orch

import (
	"github.com/bonedaddy/go-compound/v2/client"
	"github.com/bonedaddy/go-compound/v2/models"
	"log"
	"strconv"
	"sync"
)

var PAGE_SIZE string = "100"
var MAX_COMP_CONCURRENT_REQS int = 40

type Orch struct {
	c *client.Client
}

func NewOrch(c *client.Client) *Orch {
	return &Orch{c: c}
}

func (o *Orch) RetrieveAllAccounts(acctChan chan []models.Account) {
	var wg sync.WaitGroup
	var pgNum = 1
	numPages := o.RetrieveNumPages()
	wg.Add(numPages)
	limitCompoundReqs := make(chan int, MAX_COMP_CONCURRENT_REQS)
	go func() {
		for pgNum <= numPages {
			limitCompoundReqs <- 1
			go o.RetrieveAccountsFromPage(acctChan, pgNum, limitCompoundReqs, &wg)
			log.Println(pgNum)
			pgNum += 1
		}
		wg.Wait()
		close(acctChan)
	}()
}

func (o *Orch) RetrieveAccountsFromPage(ch chan []models.Account, pgNum int, limitCompoundsReqs chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	pgNumStr := strconv.Itoa(pgNum)
	acctResp, err := o.c.GetAccounts(PAGE_SIZE, pgNumStr)
	if err != nil {
		log.Fatalf("error is: %s", err)
	}
	if acctResp != nil {
		ch <- acctResp.Accounts
	}
	<-limitCompoundsReqs
}

func (o *Orch) RetrieveNumPages() int {
	// make dummy req to get total pages
	acctResp, _ := o.c.GetAccounts(PAGE_SIZE, strconv.Itoa(0))
	return acctResp.PaginationSummary.TotalPages
}
