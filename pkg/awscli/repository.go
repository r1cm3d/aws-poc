package awscli

import (
	"aws-poc/pkg/infra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	tableName     = aws.String("MasterChargebackError")
	disputeID     = aws.String("DisputeID")
	timestamp     = aws.String("Timestamp")
	hashKeyType   = aws.String("HASH")
	rangeKeyType  = aws.String("RANGE")
	numberType    = aws.String("N")
	stringType    = aws.String("S")
	payPerRequest = aws.String("PAY_PER_REQUEST")
)

type repository struct {
	sess *session.Session
}

func newRepository(sess *session.Session) repository {
	return repository{sess: sess}
}

func (r repository) svc() (svc *dynamodb.DynamoDB) {
	env, _ := infra.LoadDefaultConf()
	sess := newSession(env["REGION"], env["ENDPOINT"])
	svc = dynamodb.New(sess)
	return
}

func (r repository) put(i interface{}) error {
	av, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: tableName,
	}

	svc := r.svc()
	if _, err = svc.PutItem(input); err != nil {
		return err
	}

	return nil
}
