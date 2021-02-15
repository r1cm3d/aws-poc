package awsrepository

import (
	"aws-poc/internal/protocol"
	"aws-poc/pkg/awssession"
	"aws-poc/pkg/test/integration"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestPutIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	cases := []struct {
		name string
		in   record
		want error
		dynamoRepository
	}{
		{"success", defaultInput(), nil, newRegister(awssession.NewLocalSession(), tableName)},
		{"parseError", defaultInput(), errParser, dynamoRepository{awssession.NewLocalSession(), tableName, errMarshaller, mockUnmarshaller, svc()}},
		{"putItemError", defaultInput(), errPutItem, dynamoRepository{awssession.NewLocalSession(), tableName, dynamodbattribute.MarshalMap, mockUnmarshaller, errPutItemMock{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.dynamoRepository.put(c.in); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func TestDeleteIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	cases := []struct {
		name string
		in   record
		want error
		dynamoRepository
	}{
		{"success", defaultInput(), nil, newRegister(awssession.NewLocalSession(), tableName)},
		//TODO: add error case
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.dynamoRepository.put(c.in)
			if got := c.dynamoRepository.delete(c.in); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func TestGetIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	cases := []struct {
		name    string
		inRec   record
		inItem  interface{}
		outErr  error
		outItem interface{}
		dynamoRepository
	}{
		{"success", disputeStub, protocol.Dispute{}, nil, disputeStub, newRegister(awssession.NewLocalSession(), tableName)},
		//TODO: add error case
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.dynamoRepository.put(c.inRec)
			if gotItem, gotErr := c.dynamoRepository.get(c.inRec, c.inItem); !reflect.DeepEqual(gotItem, c.outItem) && !reflect.DeepEqual(gotErr, c.outErr) {
				t.Errorf("%s, want: %v %v, got: %v %v", c.name, c.outItem, c.outErr, gotItem, gotErr)
			}
		})
	}
}

func setupTable() {
	cleanupTable()
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: id,
				AttributeType: stringType,
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: id,
				KeyType:       hashKeyType,
			},
		},
		BillingMode: payPerRequest,
		TableName:   tableName,
	}

	svc := svc()
	if _, err := svc.CreateTable(input); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("created the dynamoRepository", tableName)
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
