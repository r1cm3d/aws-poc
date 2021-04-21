package awsrepository

import (
	"aws-poc/internal/protocol"
	"aws-poc/pkg/awssession"
	"aws-poc/pkg/test/integration"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func TestPutIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	cases := []struct {
		name string
		in   record
		want error
		dynamoRepository
	}{
		{"success", defaultInput(), nil, newRegister(awssession.NewLocalSession(), tableName)},
		{"parseError", defaultInput(), parserError, dynamoRepository{awssession.NewLocalSession(), tableName, errMarshaller, mockUnmarshaller, mockUnmarshallerListOfMaps, svc()}},
		{"putItemError", defaultInput(), putItemError, dynamoRepository{awssession.NewLocalSession(), tableName, dynamodbattribute.MarshalMap, mockUnmarshaller, mockUnmarshallerListOfMaps, errPutItemMock{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.dynamoRepository.put(c.in); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func TestDeleteIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	cases := []struct {
		name string
		in   record
		want error
		dynamoRepository
	}{
		{"success", defaultInput(), nil, newRegister(awssession.NewLocalSession(), tableName)},
		{"error", defaultInput(), deleteError, dynamoRepository{adapter: errDeleteItemMock{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.want == nil {
				c.dynamoRepository.put(c.in)
			}
			if got := c.dynamoRepository.delete(c.in); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func TestGetIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	cases := []struct {
		name    string
		inRec   record
		inItem  interface{}
		outErr  error
		outItem interface{}
		dynamoRepository
	}{
		{"success", disputeStub, protocol.Dispute{}, nil, disputeStub, newRegister(awssession.NewLocalSession(), tableName)},
		{"unmarshallError", disputeStub, protocol.Dispute{}, unmarshallerError, nil, dynamoRepository{sess: awssession.NewLocalSession(), tableName: tableName, unmarshall: errUnmarshaller, adapter: svc()}},
		{"getError", disputeStub, protocol.Dispute{}, getError, nil, dynamoRepository{sess: awssession.NewLocalSession(), tableName: tableName, adapter: errGetMock{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.outItem != nil {
				c.dynamoRepository.put(c.inRec)
			}
			if gotItem, gotErr := c.dynamoRepository.get(c.inRec, c.inItem); !reflect.DeepEqual(gotItem, c.outItem) && !reflect.DeepEqual(gotErr, c.outErr) {
				t.Errorf("%s, want: %v %v, got: %v %v", c.name, c.outItem, c.outErr, gotItem, gotErr)
			}
		})
	}
}

func TestQueryIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	cases := []struct {
		name      string
		inItem    record
		inField   string
		inValue   string
		emptyItem interface{}
		outErr    error
		outItem   interface{}
		dynamoRepository
	}{
		{"success", disputeStub, "ID", disputeStub.ID(), protocol.Dispute{}, nil, disputeStub, newRegister(awssession.NewLocalSession(), tableName)},
		{"UnmarshallerListOfMapsError", disputeStub, "ID", disputeStub.ID(), protocol.Dispute{}, unmarshallerListOfMapsError, nil, dynamoRepository{sess: awssession.NewLocalSession(), tableName: tableName, unmarshallListOfMaps: errUnmarshallerListOfMaps, adapter: svc()}},
		{"queryError", disputeStub, "ID", disputeStub.ID(), protocol.Dispute{}, queryError, nil, dynamoRepository{sess: awssession.NewLocalSession(), tableName: tableName, unmarshallListOfMaps: dynamodbattribute.UnmarshalListOfMaps, adapter: errQueryMock{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.outItem != nil {
				c.dynamoRepository.put(c.inItem)
			}
			if gotItem, gotErr := c.dynamoRepository.query(c.inField, c.inValue, c.emptyItem); !reflect.DeepEqual(gotItem, c.outItem) && !reflect.DeepEqual(gotErr, c.outErr) {
				t.Errorf("%s, want: %v %v, got: %v %v", c.name, c.outItem, c.outErr, gotItem, gotErr)
			}
		})
	}
}
