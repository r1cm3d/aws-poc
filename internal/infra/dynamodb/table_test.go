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

func setup() {
	env, _ := infra.LoadConf("../../../scripts/env/")
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(env["REGION"]),
		Endpoint: aws.String(env["ENDPOINT"]),
	}))
	svc := dynamodb.New(sess)

	tableName := "Movies"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Year"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Title"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Year"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Title"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	fmt.Println("Created the table", tableName)
}

func TestCreateAndDeleteIntegration(t *testing.T) {
	skipShort(t)
	setup()

	time.Sleep(5 * time.Second)

	teardown()
}


func teardown() {
	env, _ := infra.LoadConf("../../../scripts/env/")
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(env["REGION"]),
		Endpoint: aws.String(env["ENDPOINT"]),
	}))
	svc := dynamodb.New(sess)

	tableName := "Movies"

	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}
	_, err := svc.DeleteTable(input)
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	fmt.Println("Deleted the table", tableName)
}

func skipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
}
