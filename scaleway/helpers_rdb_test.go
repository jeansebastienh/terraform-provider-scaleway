package scaleway

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func TestSOMENAMEWithRegionalizedIdReturnProperRegionAndId(t *testing.T) {
	assert := assert.New(t)
	meta, _ := buildMeta(&MetaConfig{
		terraformVersion: "terraform-test-unit",
	})
	data := newTestResourceDataRaw(t)
	_ = data.Set("instance_id", "fr-srr/1111-111111111111-11111111")

	region, id := extractRegionAndInstanceID(data, meta)

	assert.Equal("fr-srr", region.String())
	assert.Equal("1111-111111111111-11111111", id)
}

func TestSOMENAMEWithoutRegionalizedIdReturnDefaultRegionAndId(t *testing.T) {
	assert := assert.New(t)
	meta, _ := buildMeta(&MetaConfig{
		terraformVersion: "terraform-test-unit",
	})
	data := newTestResourceDataRaw(t)
	_ = data.Set("instance_id", "1111-111111111111-11111111")

	region, id := extractRegionAndInstanceID(data, meta)

	assert.Equal("fr-par", region.String())
	assert.Equal("1111-111111111111-11111111", id)
}
func newTestResourceDataRaw(t *testing.T) *schema.ResourceData {
	raw := make(map[string]interface{})
	return schema.TestResourceDataRaw(t, resourceScalewayRdbPrivilege().Schema, raw)
}
