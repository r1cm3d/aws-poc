package awsrepository

import (
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type (
	errPutItemMock    struct{}
	errQueryMock      struct{}
	errGetMock        struct{}
	errDeleteItemMock struct{}
	Item              struct {
		DisputeID int
		Timestamp string
	}
	errRegister struct{}
)

func (e errRegister) put(_ record) error {
	return errPutItemStub
}

func (e errRegister) delete(_ record) error {
	return errDeleteStub
}

func (i Item) ID() string {
	return strconv.Itoa(i.DisputeID)
}

func (e errPutItemMock) PutItem(_ *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, errPutItemStub
}

func (e errPutItemMock) DeleteItem(_ *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return nil, errPutItemStub
}

func (e errPutItemMock) GetItem(_ *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return nil, errPutItemStub
}

func (e errPutItemMock) Query(_ *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return nil, errPutItemStub
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
	return nil, errQueryStub
}

func (e errGetMock) PutItem(_ *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, nil
}

func (e errGetMock) DeleteItem(_ *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return nil, nil
}

func (e errGetMock) GetItem(_ *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return nil, errGetStub
}

func (e errGetMock) Query(_ *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return nil, nil
}

func (e errDeleteItemMock) PutItem(_ *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, nil
}

func (e errDeleteItemMock) DeleteItem(_ *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return nil, errDeleteStub
}

func (e errDeleteItemMock) GetItem(_ *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return nil, nil
}

func (e errDeleteItemMock) Query(_ *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return nil, nil
}

func errMarshaller(_ interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return nil, errParserStub
}

func errUnmarshaller(_ map[string]*dynamodb.AttributeValue, _ interface{}) error {
	return errUnmarshallerStub
}

func mockUnmarshaller(_ map[string]*dynamodb.AttributeValue, out interface{}) error {
	return nil
}

func mockUnmarshallerListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	return nil
}

func errUnmarshallerListOfMaps(_ []map[string]*dynamodb.AttributeValue, _ interface{}) error {
	return errUnmarshallerListOfMapsStub
}
