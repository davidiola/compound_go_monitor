package dataaccess

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bonedaddy/go-compound/v2/models"
	"log"
	"time"
)

var COMP_TABLE_NAME string = "compound"
var REGION string = "us-east-1"
var MAX_DDB_CONCURRENT_REQS int = 20
var DEFAULT_RETRYER = client.DefaultRetryer{
	NumMaxRetries:    5,
	MinThrottleDelay: time.Second * 5,
	MaxThrottleDelay: time.Second * 60,
	MinRetryDelay:    time.Millisecond * 400,
	MaxRetryDelay:    time.Second * 2,
}

type DataAccess struct {
	ddb *dynamodb.DynamoDB
}

func NewDataAccess() *DataAccess {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(REGION)},
	)
	cfg := request.WithRetryer(aws.NewConfig(), DEFAULT_RETRYER)
	client := dynamodb.New(sess, cfg)

	return &DataAccess{ddb: client}
}

type Entry struct {
	Address string
	Account models.Account
}

func (da *DataAccess) WriteAllAccounts(ch chan []models.Account) {
	limitDDBReqs := make(chan int, MAX_DDB_CONCURRENT_REQS)
	for acctList := range ch {
		for _, acct := range acctList {
			limitDDBReqs <- 1
			go da.WriteAccount(acct, limitDDBReqs)
		}
	}
}

func (da *DataAccess) WriteAccount(acct models.Account, limit chan int) {
	entry := Entry{
		Address: acct.Address,
		Account: acct,
	}
	ddbEntry, err := dynamodbattribute.MarshalMap(entry)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      ddbEntry,
		TableName: aws.String(COMP_TABLE_NAME),
	}
	_, err = da.ddb.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
	<-limit
}
