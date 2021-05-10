package mocks

import (
	"errors"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

var (
	DatabaseName  string = "dbname"
	DatabaseOwner string = "dbowner"
	InstanceID    string = "1111-11111111-111111111111"
)

type ListDatabasesRequestMatcher struct {
	ExpectedRegion       string
	ExpectedInstanceID   string
	ExpectedDatabaseName string
}

func (m ListDatabasesRequestMatcher) Matches(x interface{}) bool {
	req := x.(*rdb.ListDatabasesRequest)

	if req.Region.String() != m.ExpectedRegion {
		return false
	}
	if req.InstanceID != m.ExpectedInstanceID {
		return false
	}
	if *req.Name != m.ExpectedDatabaseName {
		return false
	}
	return true
}

func (m ListDatabasesRequestMatcher) String() string {
	return fmt.Sprintf("is equal to (%s, %s, %s)", m.ExpectedRegion, m.ExpectedInstanceID, m.ExpectedDatabaseName)
}

type DeleteDatabaseRequestMatcher struct {
	ExpectedRegion       string
	ExpectedInstanceID   string
	ExpectedDatabaseName string
}

func (m DeleteDatabaseRequestMatcher) Matches(x interface{}) bool {
	req := x.(*rdb.DeleteDatabaseRequest)

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

func (m DeleteDatabaseRequestMatcher) String() string {
	return fmt.Sprintf("is equal to (%s, %s, %s)", m.ExpectedRegion, m.ExpectedInstanceID, m.ExpectedDatabaseName)
}

func NewTestDatabase() *rdb.Database {
	db := rdb.Database{
		Name:    DatabaseName,
		Owner:   DatabaseOwner,
		Managed: true,
		Size:    42,
	}
	return &db
}

func (m *MockRdbAPIInterface) CreateDatabaseMustReturnError() {
	m.EXPECT().CreateDatabase(gomock.Any(), gomock.Any()).Return(nil, errors.New("Error"))
}

func (m *MockRdbAPIInterface) CreateDatabaseMustReturnDB(expectedRegion string) {
	matcher := CreateDatabaseRequestMatcher{
		ExpectedRegion:       expectedRegion,
		ExpectedInstanceID:   InstanceID,
		ExpectedDatabaseName: DatabaseName,
	}
	m.EXPECT().CreateDatabase(matcher, gomock.Any()).Return(NewTestDatabase(), nil)
}
func (m *MockRdbAPIInterface) ListDatabasesMustReturnError() {
	m.EXPECT().ListDatabases(gomock.Any(), gomock.Any()).Return(nil, errors.New("Error"))
}
func (m *MockRdbAPIInterface) ListDatabasesMustReturnDB(expectedRegion string) {
	matcher := ListDatabasesRequestMatcher{
		ExpectedRegion:       expectedRegion,
		ExpectedInstanceID:   InstanceID,
		ExpectedDatabaseName: DatabaseName,
	}
	dbs := make([]*rdb.Database, 0)
	dbs = append(dbs, NewTestDatabase())
	resp := rdb.ListDatabasesResponse{
		Databases:  dbs,
		TotalCount: 1,
	}
	m.EXPECT().ListDatabases(matcher, gomock.Any()).Return(&resp, nil)
}

func (m *MockRdbAPIInterface) DeleteDatabaseMustReturnError() {
	m.EXPECT().DeleteDatabase(gomock.Any(), gomock.Any()).Return(errors.New("Error"))
}

func (m *MockRdbAPIInterface) DeleteDatabaseReturnNil(expectedRegion string) {
	matcher := DeleteDatabaseRequestMatcher{
		ExpectedRegion:       expectedRegion,
		ExpectedInstanceID:   InstanceID,
		ExpectedDatabaseName: DatabaseName,
	}
	m.EXPECT().DeleteDatabase(matcher, gomock.Any()).Return(nil)
}

type CreateDatabaseRequestMatcher struct {
	ExpectedRegion       string
	ExpectedInstanceID   string
	ExpectedDatabaseName string
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

type ListPrivilegesRequestMatcher struct {
	ExpectedRegion       string
	ExpectedInstanceID   string
	ExpectedDatabaseName string
	ExpectedUserName     string
}

func (m ListPrivilegesRequestMatcher) Matches(x interface{}) bool {
	req := x.(*rdb.ListPrivilegesRequest)

	if req.Region.String() != m.ExpectedRegion {
		return false
	}
	if req.InstanceID != m.ExpectedInstanceID {
		return false
	}
	if *req.DatabaseName != m.ExpectedDatabaseName {
		return false
	}
	if *req.UserName != m.ExpectedUserName {
		return false
	}
	return true
}

func (m ListPrivilegesRequestMatcher) String() string {
	return fmt.Sprintf("is equal to (%s, %s, %s, %s)", m.ExpectedRegion, m.ExpectedInstanceID, m.ExpectedDatabaseName, m.ExpectedUserName)
}

type SetPrivilegeRequestMatcher struct {
	ExpectedRegion       string
	ExpectedInstanceID   string
	ExpectedDatabaseName string
	ExpectedUserName     string
	ExpectedPermission   string
}

func (m SetPrivilegeRequestMatcher) Matches(x interface{}) bool {
	req := x.(*rdb.SetPrivilegeRequest)

	if req.Region.String() != m.ExpectedRegion {
		return false
	}
	if req.InstanceID != m.ExpectedInstanceID {
		return false
	}
	if req.DatabaseName != m.ExpectedDatabaseName {
		return false
	}
	if req.UserName != m.ExpectedUserName {
		return false
	}
	if req.Permission.String() != m.ExpectedPermission {
		return false
	}
	return true
}

func (m SetPrivilegeRequestMatcher) String() string {
	return fmt.Sprintf("is equal to (%s, %s, %s, %s, %s)", m.ExpectedRegion, m.ExpectedInstanceID, m.ExpectedDatabaseName, m.ExpectedUserName, m.ExpectedPermission)
}

func (m *MockRdbAPIInterface) SetPrivilegeRequest(expectedRegion string) {
	matcher := SetPrivilegeRequestMatcher{
		ExpectedRegion:       expectedRegion,
		ExpectedInstanceID:   InstanceID,
		ExpectedDatabaseName: DatabaseName,
		ExpectedUserName:     DatabaseOwner,
		ExpectedPermission:   "all",
	}
	m.EXPECT().SetPrivilege(matcher, gomock.Any()).Return(nil, nil)
}
func (m *MockRdbAPIInterface) ListPrivilegesReturnPrivilege(expectedRegion string) {
	matcher := ListPrivilegesRequestMatcher{
		ExpectedRegion:       expectedRegion,
		ExpectedInstanceID:   InstanceID,
		ExpectedDatabaseName: DatabaseName,
		ExpectedUserName:     DatabaseOwner,
	}
	privs := make([]*rdb.Privilege, 0)
	privs = append(privs, NewTestPrivilege())
	resp := rdb.ListPrivilegesResponse{
		Privileges: privs,
		TotalCount: 1,
	}
	m.EXPECT().ListPrivileges(matcher, gomock.Any()).Return(&resp, nil)
}
func NewTestPrivilege() *rdb.Privilege {
	priv := rdb.Privilege{
		Permission:   rdb.PermissionAll,
		DatabaseName: DatabaseName,
		UserName:     DatabaseOwner,
	}
	return &priv
}
