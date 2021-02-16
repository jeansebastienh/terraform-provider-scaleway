package scaleway

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	defaultRdbInstanceTimeout = 15 * time.Minute
)

// rdbAPIWithRegion returns a new lb API and the region for a Create request
func rdbAPIWithRegion(d *schema.ResourceData, m interface{}) (RdbAPIInterface, scw.Region, error) {
	meta := m.(*Meta)
	var rdbAPI RdbAPIInterface
	if meta.mockedAPI != nil {
		rdbAPI = meta.mockedAPI.(RdbAPIInterface)
	} else {
		rdbAPI = rdb.NewAPI(meta.scwClient)
	}

	region, err := extractRegion(d, meta)
	if err != nil {
		return nil, "", err
	}
	return rdbAPI, region, nil
}

// rdbAPIWithRegionAndID returns an lb API with region and ID extracted from the state
func rdbAPIWithRegionAndID(m interface{}, id string) (*rdb.API, scw.Region, string, error) {
	meta := m.(*Meta)
	rdbAPI := rdb.NewAPI(meta.scwClient)

	region, ID, err := parseRegionalID(id)
	if err != nil {
		return nil, "", "", err
	}
	return rdbAPI, region, ID, nil
}

func flattenRdbInstanceReadReplicas(readReplicas []*rdb.Endpoint) interface{} {
	replicasI := []map[string]interface{}(nil)
	for _, readReplica := range readReplicas {
		replicasI = append(replicasI, map[string]interface{}{
			"ip":   flattenIPPtr(readReplica.IP),
			"port": int(readReplica.Port),
			"name": flattenStringPtr(readReplica.Name),
		})
	}
	return replicasI
}

func flattenInstanceSettings(settings []*rdb.InstanceSetting) interface{} {
	res := make(map[string]string)
	for _, value := range settings {
		res[value.Name] = value.Value
	}

	return res
}

func expandInstanceSettings(i interface{}) []*rdb.InstanceSetting {
	rawRule := i.(map[string]interface{})
	var res []*rdb.InstanceSetting
	for key, value := range rawRule {
		res = append(res, &rdb.InstanceSetting{
			Name:  key,
			Value: value.(string),
		})
	}

	return res
}
