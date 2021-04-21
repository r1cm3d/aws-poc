package awsrepository

import (
	"aws-poc/pkg/awssession"
	"fmt"

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
	marshall             func(in interface{}) (map[string]*dynamodb.AttributeValue, error)
	unmarshall           func(m map[string]*dynamodb.AttributeValue, out interface{}) error
	unmarshallListOfMaps func(l []map[string]*dynamodb.AttributeValue, out interface{}) error
	adapter              interface {
		PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
		DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
		GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
		Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	}
	repository interface {
		put(rec record) error
		delete(rec record) error
	}
	dynamoRepository struct {
		sess      *session.Session
		tableName *string
		marshall
		unmarshall
		unmarshallListOfMaps
		adapter
	}
)

func newRegister(sess *session.Session, tableName *string) dynamoRepository {
	return dynamoRepository{sess, tableName, dynamodbattribute.MarshalMap, dynamodbattribute.UnmarshalMap, dynamodbattribute.UnmarshalListOfMaps, svc()}
}

func svc() (svc *dynamodb.DynamoDB) {
	sess := awssession.NewLocalSession()
	svc = dynamodb.New(sess)
	return
}

func (r dynamoRepository) put(rec record) error {
	item, err := r.marshall(rec)
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

func (r dynamoRepository) delete(rec record) error {
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

func (r dynamoRepository) get(rec record, item interface{}) (interface{}, error) {
	input := &dynamodb.GetItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(rec.ID()),
			},
		},
	}

	ri, err := r.GetItem(input)
	if err != nil {
		return nil, err
	}

	if err := r.unmarshall(ri.Item, &item); err != nil {
		return nil, err
	} else {
		return item, err
	}
}

func (r dynamoRepository) query(field string, value string, items interface{}) (interface{}, error) {
	input := &dynamodb.QueryInput{
		TableName: r.tableName,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(value),
			},
		},
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :v1", field)),
	}

	ri, err := r.Query(input)
	if err != nil {
		return nil, err
	}

	if err := r.unmarshallListOfMaps(ri.Items, &items); err != nil {
		return nil, err
	} else {
		return items, err
	}
}
