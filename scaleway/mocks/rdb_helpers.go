package mocks

import (
	"errors"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

var (
	databaseName  string = "dbname"
	databaseOwner string = "dbowner"
	instanceID    string = "1111-11111111-111111111111"
)

type ListDatabasesRequestMatcher struct {
	ExpectedRegion       string
	ExpectedInstanceID   string
	ExpectedDatabaseName string
	errorString          string
}

func (m ListDatabasesRequestMatcher) Matches(x interface{}) bool {
	req := x.(*rdb.ListDatabasesRequest)

	if req.Region.String() != m.ExpectedRegion {
		return false
	}
	if req.InstanceID != m.ExpectedInstanceID {
		return false
	}
	if fmt.Sprintf("%s", *req.Name) != m.ExpectedDatabaseName {
		return false
	}
	return true
}

func (m ListDatabasesRequestMatcher) String() string {
	return fmt.Sprintf("is equal to (%s, %s, %s)", m.ExpectedRegion, m.ExpectedInstanceID, m.ExpectedDatabaseName)
}

func NewTestDatabase() *rdb.Database {
	db := rdb.Database{
		Name:    databaseName,
		Owner:   databaseOwner,
		Managed: true,
		Size:    42,
	}
	return &db
}

func (m *MockRdbApiInterface) CreateDatabaseMustReturnError() {
	m.EXPECT().CreateDatabase(gomock.Any(), gomock.Any()).Return(nil, errors.New("Error"))
}

func (m *MockRdbApiInterface) CreateDatabaseMustReturnDb(expectedRegion string) {
	matcher := CreateDatabaseRequestMatcher{
		ExpectedRegion:       expectedRegion,
		ExpectedInstanceID:   instanceID,
		ExpectedDatabaseName: databaseName,
	}
	m.EXPECT().CreateDatabase(matcher, gomock.Any()).Return(NewTestDatabase(), nil)
}
func (m *MockRdbApiInterface) ListDatabasesMustReturnError() {
	m.EXPECT().ListDatabases(gomock.Any(), gomock.Any()).Return(nil, errors.New("Error"))
}
func (m *MockRdbApiInterface) ListDatabasesMustReturnDb(expectedRegion string) {
	matcher := ListDatabasesRequestMatcher{
		ExpectedRegion:       expectedRegion,
		ExpectedInstanceID:   instanceID,
		ExpectedDatabaseName: databaseName,
	}
	dbs := make([]*rdb.Database, 0)
	dbs = append(dbs, NewTestDatabase())
	resp := rdb.ListDatabasesResponse{
		Databases:  dbs,
		TotalCount: 1,
	}
	m.EXPECT().ListDatabases(matcher, gomock.Any()).Return(&resp, nil)

}

type CreateDatabaseRequestMatcher struct {
	ExpectedRegion       string
	ExpectedInstanceID   string
	ExpectedDatabaseName string
	errorString          string
}

func (m CreateDatabaseRequestMatcher) Matches(x interface{}) bool {
	req := x.(*rdb.CreateDatabaseRequest)

	if req.Region.String() != m.ExpectedRegion {
		return false
	}
	if req.InstanceID != m.ExpectedInstanceID {
		return false
	}
	if req.Name != m.ExpectedDatabaseName {
		return false
	}
	return true
}

func (m CreateDatabaseRequestMatcher) String() string {
	return fmt.Sprintf("is equal to (%s, %s, %s)", m.ExpectedRegion, m.ExpectedInstanceID, m.ExpectedDatabaseName)
}