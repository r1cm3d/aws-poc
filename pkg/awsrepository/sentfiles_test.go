package awsrepository

//import (
//	"aws-poc/internal/protocol"
//	"aws-poc/pkg/awssession"
//	"aws-poc/pkg/test/integration"
//	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
//	"reflect"
//	"testing"
//)
//
//func TestUnsentFilesIntegration(t *testing.T) {
//	integration.SkipShort(t)
//	setupTable()
//	defer cleanupTable()
//	cases := []struct {
//		name      string
//		inItem    record
//		inField   string
//		inValue   string
//		emptyItem interface{}
//		outErr    error
//		outItem   interface{}
//		dynamoRepository
//	}{
//		{"success", disputeStub, "ID", disputeStub.ID(), protocol.Dispute{}, nil, disputeStub, newRegister(awssession.NewLocalSession(), tableName)},
//		{"UnmarshallerListOfMapsError", disputeStub, "ID", disputeStub.ID(), protocol.Dispute{}, errUnmarshallerListOfMapsStub, nil, dynamoRepository{sess: awssession.NewLocalSession(), tableName: tableName, unmarshallListOfMaps: errUnmarshallerListOfMaps, adapter: svc()}},
//		{"errQueryStub", disputeStub, "ID", disputeStub.ID(), protocol.Dispute{}, errQueryStub, nil, dynamoRepository{sess: awssession.NewLocalSession(), tableName: tableName, unmarshallListOfMaps: dynamodbattribute.UnmarshalListOfMaps, adapter: errQueryMock{}}},
//	}
//
//	for _, c := range cases {
//		t.Run(c.name, func(t *testing.T) {
//			if c.outItem != nil {
//				c.dynamoRepository.put(c.inItem)
//			}
//			if gotItem, gotErr := c.dynamoRepository.query(c.inField, c.inValue, c.emptyItem); !reflect.DeepEqual(gotItem, c.outItem) && !reflect.DeepEqual(gotErr, c.outErr) {
//				t.Errorf("%s, want: %v %v, got: %v %v", c.name, c.outItem, c.outErr, gotItem, gotErr)
//			}
//		})
//	}
//}
