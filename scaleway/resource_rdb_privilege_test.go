package scaleway

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	mock "github.com/scaleway/terraform-provider-scaleway/scaleway/mocks"
	"github.com/stretchr/testify/assert"
)

/*func init() {
	resource.AddTestSweepers("scaleway_rdb_database", &resource.Sweeper{
		Name: "scaleway_rdb_database",
		F:    testSweepRDBInstance,
	})
}*/

func TestAccScalewayRdbPrivilege_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	instanceName := "TestAccScalewayRdbPrivilege_Basic"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayRdbInstanceDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_rdb_instance" "instance" {
						name = "%s"
						node_type = "db-dev-s"
						engine = "PostgreSQL-12"
						is_ha_cluster = false
						tags = [ "terraform-test", "scaleway_rdb_user", "minimal" ]
					}

					resource "scaleway_rdb_database" "db" {
						instance_id = scaleway_rdb_instance.instance.id
						name = "foo"
					}

					resource "scaleway_rdb_user" "foo" {
						instance_id = scaleway_rdb_instance.instance.id
						name = "foo"
						password = "R34lP4sSw#Rd"
					}

					resource "scaleway_rdb_privilege" "priv" {
						instance_id   = scaleway_rdb_instance.instance.id
						user_name     = scaleway_rdb_user.foo.name
						database_name = scaleway_rdb_database.db.name
						permission    = "all"
					}`, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdbPrivilegeExists(tt, "scaleway_rdb_instance.instance", "scaleway_rdb_database.db", "scaleway_rdb_user.foo"),
					resource.TestCheckResourceAttr("scaleway_rdb_privilege.priv", "permission", "all"),
				),
			},
		},
	})
}

func testAccCheckRdbPrivilegeExists(tt *TestTools, instance string, database string, user string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		instanceResource, ok := state.RootModule().Resources[instance]
		if !ok {
			return fmt.Errorf("resource not found: %s", instance)
		}

		databaseResource, ok := state.RootModule().Resources[database]
		if !ok {
			return fmt.Errorf("resource database not found: %s", database)
		}

		userResource, ok := state.RootModule().Resources[user]
		if !ok {
			return fmt.Errorf("resource not found: %s", user)
		}

		rdbAPI, region, _, err := rdbAPIWithRegionAndID(tt.Meta, instanceResource.Primary.ID)
		if err != nil {
			return err
		}

		_, instanceID, databaseName, err := resourceScalewayRdbDatabaseParseID(databaseResource.Primary.ID)
		if err != nil {
			return err
		}

		_, userName, err := resourceScalewayRdbUserParseID(userResource.Primary.ID)
		if err != nil {
			return err
		}

		databases, err := rdbAPI.ListPrivileges(&rdb.ListPrivilegesRequest{
			Region:       region,
			InstanceID:   instanceID,
			DatabaseName: &databaseName,
			UserName:     &userName,
		})
		if err != nil {
			return err
		}

		if len(databases.Privileges) != 1 {
			return fmt.Errorf("no privilege found")
		}

		return nil
	}
}

func TestValidationPrivilegePermissionWithWrongPermissionReturnError(t *testing.T) {
	assert := assert.New(t)
	warnings, errors := validationPrivilegePermission()("NotAPermission", "key")
	assert.Empty(warnings)
	assert.Len(errors, 1)
	assert.Error(errors[0])
	assert.Equal("'key' is not a valid permission", errors[0].Error())
}

func TestValidationPrivilegePermissionWithoutStringReturnError(t *testing.T) {
	assert := assert.New(t)
	warnings, errors := validationPrivilegePermission()(1, "key")
	assert.Empty(warnings)
	assert.Len(errors, 1)
	assert.Error(errors[0])
	assert.Equal("'key' is not a string", errors[0].Error())
}

func TestValidationPrivilegePermission(t *testing.T) {
	assert := assert.New(t)
	warnings, errors := validationPrivilegePermission()("none", "key")
	assert.Empty(warnings)
	assert.Empty(errors)
}
func TestResourceScalewayRdbPrivilegeRead(t *testing.T) {
	// init
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	meta, rdbAPI := NewMeta(ctrl)
	data := NewTestResourceDataRawForResourceRDBPrivilege(t, "fr-srr/1111-11111111-111111111111")

	// mocking
	rdbAPI.ListPrivilegesReturnPrivilege("fr-srr")

	// run
	diags := resourceScalewayRdbPrivilegeRead(mock.NewMockContext(ctrl), data, meta)

	// assertions
	assert.Len(diags, 0)
	assertResourcePrivilege(assert, data, "fr-srr")
}

func TestResourceScalewayRdbPrivilegeReadWithNonRegionalizedIDUseDefault(t *testing.T) {
	// init
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	meta, rdbAPI := NewMeta(ctrl)
	data := NewTestResourceDataRawForResourceRDBPrivilege(t, "1111-11111111-111111111111")

	// mocking
	rdbAPI.ListPrivilegesReturnPrivilege("fr-par")

	// run
	diags := resourceScalewayRdbPrivilegeRead(mock.NewMockContext(ctrl), data, meta)

	// assertions
	assert.Len(diags, 0)
	assertResourcePrivilege(assert, data, "fr-par")
}

func TestResourceScalewayRdbPrivilegeCreateWithoutRegionalIDReturnsExpectedResource(t *testing.T) {
	// init
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	meta, rdbAPI := NewMeta(ctrl)
	data := NewTestResourceDataRawForResourceRDBPrivilege(t, "1111-11111111-111111111111")

	// mocking
	rdbAPI.SetPrivilegeRequest("fr-par")
	rdbAPI.ListPrivilegesReturnPrivilege("fr-par")

	// run
	diags := resourceScalewayRdbPrivilegeCreate(mock.NewMockContext(ctrl), data, meta)

	// assertions
	assert.Len(diags, 0)
	assertResourcePrivilege(assert, data, "fr-par")
}

func TestResourceScalewayRdbPrivilegeCreateWithRegionalIDReturnsExpectedResource(t *testing.T) {
	// init
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	meta, rdbAPI := NewMeta(ctrl)
	data := NewTestResourceDataRawForResourceRDBPrivilege(t, "fr-srr/1111-11111111-111111111111")

	// mocking
	rdbAPI.SetPrivilegeRequest("fr-srr")
	rdbAPI.ListPrivilegesReturnPrivilege("fr-srr")

	// run
	diags := resourceScalewayRdbPrivilegeCreate(mock.NewMockContext(ctrl), data, meta)

	// assertions
	assert.Len(diags, 0)
	assertResourcePrivilege(assert, data, "fr-srr")
}

func NewTestResourceDataRawForResourceRDBPrivilege(t *testing.T, uuid string) *schema.ResourceData {
	raw := make(map[string]interface{})
	raw["instance_id"] = uuid
	raw["database_name"] = "dbname"
	raw["user_name"] = "dbowner"
	raw["permission"] = "all"
	return schema.TestResourceDataRaw(t, resourceScalewayRdbPrivilege().Schema, raw)
}

func assertResourcePrivilege(assert *assert.Assertions, data *schema.ResourceData, region string) {
	assert.Equal(fmt.Sprintf("%s/1111-11111111-111111111111", region), data.Get("instance_id"))
	assert.Equal("dbname", data.Get("database_name"))
	assert.Equal("dbowner", data.Get("user_name"))
	assert.Equal(rdb.PermissionAll, rdb.Permission(data.Get("permission").(string)))
}
