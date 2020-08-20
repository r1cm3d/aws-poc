package awscli

import (
	"aws-poc/internal/infra"
	"github.com/aws/aws-sdk-go/aws"
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
	marshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error)
}

type persistent interface {
	putItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
}

type awsRepository struct{}

func (m awsRepository) marshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(in)
}

func (m awsRepository) putItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	svc := svc()
	return svc.PutItem(input)
}

type repository struct {
	marshaller
	persistent
}

func newRepository() repository {
	awsRepository	 := awsRepository{}
	return repository{marshaller: awsRepository, persistent: awsRepository}
}

func svc() (svc *dynamodb.DynamoDB) {
	env, _ := infra.LoadDefaultConf()
	sess := newSession(env["REGION"], env["ENDPOINT"])
	svc = dynamodb.New(sess)
	return
}

func (t repository) put(i interface{}) error {
	av, err := t.marshalMap(i)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: tableName,
	}

	_, err = t.putItem(input)
	if err != nil {
		return err
	}

	return nil
}
