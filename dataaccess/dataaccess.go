package dataaccess

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bonedaddy/go-compound/v2/models"
	"log"
)

var COMP_TABLE_NAME string = "compound"
var REGION string = "us-east-1"

type DataAccess struct {
	ddb *dynamodb.DynamoDB
}

func NewDataAccess() *DataAccess {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(REGION)},
	)
	client := dynamodb.New(sess)

	return &DataAccess{ddb: client}
}

type Entry struct {
	Address string
	Account models.Account
}

func (da *DataAccess) WriteAccount(acct models.Account) {
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
}
