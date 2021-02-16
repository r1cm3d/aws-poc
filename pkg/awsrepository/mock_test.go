package awsrepository

import (
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type (
	errPutItemMock struct{}
	errQueryMock struct{}
	Item           struct {
		DisputeID int
		Timestamp string
	}
	errRegister struct{}
)

func (e errRegister) put(_ record) error {
	return putItemError
}

func (e errRegister) delete(_ record) error {
	return deleteError
}

func (i Item) ID() string {
	return strconv.Itoa(i.DisputeID)
}

func (e errPutItemMock) PutItem(_ *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, putItemError
}

func (e errPutItemMock) DeleteItem(_ *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return nil, putItemError
}

func (e errPutItemMock) GetItem(_ *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return nil, putItemError
}

func (e errPutItemMock) Query(_ *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return nil, putItemError
}

func (e errQueryMock) PutItem(_ *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, nil
}

func (e errQueryMock) DeleteItem(_ *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return nil, nil
}

func (e errQueryMock) GetItem(_ *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return nil, nil
}

func (e errQueryMock) Query(_ *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return nil, queryError
}


func errMarshaller(_ interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return nil, parserError
}

func mockUnmarshaller(_ map[string]*dynamodb.AttributeValue, out interface{}) error {
	return nil
}

func mockUnmarshallerListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	return nil
}

func errUnmarshallerListOfMaps(_ []map[string]*dynamodb.AttributeValue, _ interface{}) error {
	return unmarshallerListOfMapsError
}
