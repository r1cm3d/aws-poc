package awsrepository

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func cleanupTable() {
	svc := svc()

	input := &dynamodb.DeleteTableInput{
		TableName: tableName,
	}

	if out, _ := svc.DeleteTable(input); out != nil {
		log.Printf("table %v deleted", tableName)
	}
}
