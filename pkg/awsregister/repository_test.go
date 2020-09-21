package awsdynamo

import (
	"aws-poc/pkg/awssession"
	"aws-poc/pkg/test/integration"
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Item struct {
	DisputeID int
	Timestamp string
}

func TestPutIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	i := Item{
		DisputeID: 666, Timestamp: "2020-04-17T17:19:19.831Z",
	}
	table := newRepository(awssession.NewLocalSession())

	if err := table.put(i); err != nil {
		t.Errorf("put fails: %d", err)
	}
}

func setupTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: disputeID,
				AttributeType: numberType,
			},
			{
				AttributeName: timestamp,
				AttributeType: stringType,
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: disputeID,
				KeyType:       hashKeyType,
			},
			{
				AttributeName: timestamp,
				KeyType:       rangeKeyType,
			},
		},
		BillingMode: payPerRequest,
		TableName:   tableName,
	}

	r := newRepository(awssession.NewLocalSession())
	svc := r.svc()
	if _, err := svc.CreateTable(input); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("created the repository", tableName)
}

func cleanupTable() {
	r := newRepository(awssession.NewLocalSession())
	svc := r.svc()

	input := &dynamodb.DeleteTableInput{
		TableName: tableName,
	}
	if _, err := svc.DeleteTable(input); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("deleted the repository", tableName)
}
