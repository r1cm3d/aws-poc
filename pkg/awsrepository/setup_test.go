package awsrepository

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

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
