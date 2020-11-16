package awsregister

import (
	"aws-poc/pkg/awssession"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	hashKeyType   = aws.String("HASH")
	rangeKeyType  = aws.String("RANGE")
	numberType    = aws.String("N")
	stringType    = aws.String("S")
	payPerRequest = aws.String("PAY_PER_REQUEST")
)

type (
	mapMarshaller func(in interface{}) (map[string]*dynamodb.AttributeValue, error)
	register      interface {
		PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
		DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	}
	dynamoRegister struct {
		sess      *session.Session
		tableName *string
		mapMarshaller
		register
	}
)

func newRegister(sess *session.Session, tableName *string) dynamoRegister {
	return dynamoRegister{sess, tableName, dynamodbattribute.MarshalMap, svc()}
}

func svc() (svc *dynamodb.DynamoDB) {
	sess := awssession.NewLocalSession()
	svc = dynamodb.New(sess)
	return
}

func (r dynamoRegister) put(i interface{}) error {
	item, err := r.mapMarshaller(i)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: r.tableName,
	}

	if _, err = r.PutItem(input); err != nil {
		return err
	}

	return nil
}

func (r dynamoRegister) delete(key int, timestamp string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"DisputeID": {
				N: aws.String(strconv.Itoa(key)),
			},
			"Timestamp": {
				S: aws.String(timestamp),
			},
		},
	}

	if _, err := r.DeleteItem(input); err != nil {
		return err
	}

	return nil
}
