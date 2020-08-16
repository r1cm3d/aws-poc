package dynamodb

import (
	"aws-poc/internal/infra"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	tableName     = aws.String("MasterChargebackError")
	disputeId     = aws.String("DisputeId")
	timestamp     = aws.String("Timestamp")
	hashKeyType   = aws.String("HASH")
	rangeKeyType  = aws.String("RANGE")
	numberType    = aws.String("N")
	stringType    = aws.String("S")
	payPerRequest = aws.String("PAY_PER_REQUEST")
)

type marshaller interface {
	MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error)
}

type awsMarshaller struct{}

func (m awsMarshaller) MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(in)
}

type table struct {
	marshaller
}

func newTable() table {
	return table{
		awsMarshaller{},
	}
}

func svc() (svc *dynamodb.DynamoDB) {
	env, _ := infra.LoadDefaultConf()
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(env["REGION"]),
		Endpoint: aws.String(env["ENDPOINT"]),
	}))
	svc = dynamodb.New(sess)
	return
}

func (t table) put(i interface{}) error {
	svc := svc()

	av, err := t.MarshalMap(i)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: tableName,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
