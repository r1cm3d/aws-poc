package awsregister

import (
	"aws-poc/pkg/awssession"
	"aws-poc/pkg/test/integration"
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	tableName = aws.String("ChargebackError_TEST")
	disputeID = aws.String("DisputeID")
	timestamp = aws.String("Timestamp")
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
	table := newRegister(awssession.NewLocalSession(), tableName)

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

	r := newRegister(awssession.NewLocalSession(), tableName)
	svc := r.svc()
	if _, err := svc.CreateTable(input); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("created the register", tableName)
}

func cleanupTable() {
	r := newRegister(awssession.NewLocalSession(), tableName)
	svc := r.svc()

	input := &dynamodb.DeleteTableInput{
		TableName: tableName,
	}
	if _, err := svc.DeleteTable(input); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("deleted the register", tableName)
}
