package dynamodb

import (
	"aws-poc/internal/infra"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"testing"
	"time"
)

var (
	tableName     = aws.String("MasterChargebackError")
	disputeId     = aws.String("dispute_id")
	timestamp     = aws.String("timestamp")
	hashKeyType   = aws.String("HASH")
	rangeKeyType  = aws.String("RANGE")
	numberType    = aws.String("N")
	stringType    = aws.String("S")
	payPerRequest = aws.String("PAY_PER_REQUEST")
)

func setup() {
	env, _ := infra.LoadDefaultConf()
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(env["REGION"]),
		Endpoint: aws.String(env["ENDPOINT"]),
	}))
	svc := dynamodb.New(sess)

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: disputeId,
				AttributeType: numberType,
			},
			{
				AttributeName: timestamp,
				AttributeType: stringType,
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: disputeId,
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

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("created the table", tableName)
}

func TestCreateAndDeleteIntegration(t *testing.T) {
	skipShort(t)
	setup()

	time.Sleep(5 * time.Second)

	teardown()
}

func teardown() {
	env, _ := infra.LoadDefaultConf()
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(env["REGION"]),
		Endpoint: aws.String(env["ENDPOINT"]),
	}))
	svc := dynamodb.New(sess)

	input := &dynamodb.DeleteTableInput{
		TableName: tableName,
	}
	_, err := svc.DeleteTable(input)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("deleted the table", tableName)
}

func skipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}
