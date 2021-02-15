package awsrepository

import (
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type (
	errPutItemMock struct{}
	Item           struct {
		DisputeID int
		Timestamp string
	}
	errRegister struct{}
)

func (e errRegister) put(_ record) error {
	return errPutItem
}

func (e errRegister) delete(_ record) error {
	return errDelete
}

func (i Item) ID() string {
	return strconv.Itoa(i.DisputeID)
}

func (e errPutItemMock) PutItem(_ *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, errPutItem
}

func (e errPutItemMock) DeleteItem(_ *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return nil, errPutItem
}

func (e errPutItemMock) GetItem(_ *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return nil, errPutItem
}

func (e errPutItemMock) Query(_ *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return nil, errPutItem
}

func errMarshaller(_ interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return nil, errParser
}

func mockUnmarshaller(_ map[string]*dynamodb.AttributeValue, out interface{}) error {
	return nil
}

func mockUnmarshallerListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	return nil
}
