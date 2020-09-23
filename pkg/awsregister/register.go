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
	rangeKeyType  = aws.String("RANGE")
	numberType    = aws.String("N")
	stringType    = aws.String("S")
	payPerRequest = aws.String("PAY_PER_REQUEST")
)

type register struct {
	sess      *session.Session
	tableName *string
}

func newRegister(sess *session.Session, tableName *string) register {
	return register{sess, tableName}
}

func (r register) svc() (svc *dynamodb.DynamoDB) {
	sess := awssession.NewLocalSession()
	svc = dynamodb.New(sess)
	return
}

func (r register) put(i interface{}) error {
	av, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: r.tableName,
	}

	svc := r.svc()
	if _, err = svc.PutItem(input); err != nil {
		return err
	}

	return nil
}
