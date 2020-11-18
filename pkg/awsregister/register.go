package awsregister

import (
	"aws-poc/pkg/awssession"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	hashKeyType   = aws.String("HASH")
	stringType    = aws.String("S")
	payPerRequest = aws.String("PAY_PER_REQUEST")
)

type (
	record interface {
		ID() string
	}
	mapMarshaller   func(in interface{}) (map[string]*dynamodb.AttributeValue, error)
	registerAdapter interface {
		PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
		DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	}
	register interface {
		put(rec record) error
		delete(rec record) error
	}
	dynamoRegister struct {
		sess      *session.Session
		tableName *string
		mapMarshaller
		registerAdapter
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

func (r dynamoRegister) put(rec record) error {
	item, err := r.mapMarshaller(rec)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: r.tableName,
	}

	item["ID"] = &dynamodb.AttributeValue{
		S: aws.String(rec.ID()),
	}

	if _, err = r.PutItem(input); err != nil {
		return err
	}

	return nil
}

func (r dynamoRegister) delete(rec record) error {
	input := &dynamodb.DeleteItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(rec.ID()),
			},
		},
	}

	if _, err := r.DeleteItem(input); err != nil {
		return err
	}

	return nil
}
