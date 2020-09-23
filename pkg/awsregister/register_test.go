package awsregister

import (
	"aws-poc/pkg/awssession"
	"aws-poc/pkg/test/integration"
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	tableName = aws.String("ChargebackError_TEST")
	disputeID = aws.String("DisputeID")
	timestamp = aws.String("Timestamp")
)

type (
	Item struct {
		DisputeID int
		Timestamp string
	}
)

func TestPutIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	defaultInput := Item{
		DisputeID: 666,
		Timestamp: "2020-04-17T17:19:19.831Z",
	}
	err := errors.New("parseError")
	cases := []struct {
		name string
		in   interface{}
		want error
		register
	}{
		{"success", defaultInput, nil, newRegister(awssession.NewLocalSession(), tableName)},
		{"parseError", defaultInput, err, register{awssession.NewLocalSession(), tableName, func(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
			return nil, err
		}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.register.put(c.in); !reflect.DeepEqual(c.want, got) {
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

	if out, _ := svc.DeleteTable(input); out != nil {
		log.Printf("table %v deleted", tableName)
	}
}
