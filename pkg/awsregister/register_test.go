package awsregister

import (
	"aws-poc/pkg/awssession"
	"aws-poc/pkg/test/integration"
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	tableName  = aws.String("ChargebackError_TEST")
	disputeID  = aws.String("DisputeID")
	timestamp  = aws.String("Timestamp")
	errParser  = errors.New("parseError")
	errPutItem = errors.New("putItemError")
)

type (
	errPutItemMock struct{}
	Item           struct {
		DisputeID int
		Timestamp string
	}
)

func (e errPutItemMock) PutItem(_ *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, errPutItem
}

func (e errPutItemMock) DeleteItem(_ *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return nil, errPutItem
}

func errMarshaller(_ interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return nil, errParser
}

func TestPutIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	defaultInput := Item{
		DisputeID: 666,
		Timestamp: "2020-04-17T17:19:19.831Z",
	}
	cases := []struct {
		name string
		in   interface{}
		want error
		dynamoRegister
	}{
		{"success", defaultInput, nil, newRegister(awssession.NewLocalSession(), tableName)},
		{"parseError", defaultInput, errParser, dynamoRegister{awssession.NewLocalSession(), tableName, errMarshaller, svc()}},
		{"putItemError", defaultInput, errPutItem, dynamoRegister{awssession.NewLocalSession(), tableName, dynamodbattribute.MarshalMap, errPutItemMock{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.dynamoRegister.put(c.in); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func TestDeleteIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	defaultInput := Item{
		DisputeID: 666,
		Timestamp: "2020-04-17T17:19:19.831Z",
	}
	cases := []struct {
		name string
		in   Item
		want error
		dynamoRegister
	}{
		{"success", defaultInput, nil, newRegister(awssession.NewLocalSession(), tableName)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.dynamoRegister.put(c.in)
			if got := c.dynamoRegister.delete(c.in.DisputeID, c.in.Timestamp); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func setupTable() {
	cleanupTable()
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

	svc := svc()
	if _, err := svc.CreateTable(input); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("created the dynamoRegister", tableName)
}

func cleanupTable() {
	svc := svc()

	input := &dynamodb.DeleteTableInput{
		TableName: tableName,
	}

	if out, _ := svc.DeleteTable(input); out != nil {
		log.Printf("table %v deleted", tableName)
	}
}
