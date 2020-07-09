package common

import (
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2020-04-01/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func GetThroughputFromResult(throughputResponse documentdb.ThroughputSettingsGetResults) *int32 {
	props := throughputResponse.ThroughputSettingsGetProperties
	if props == nil {
		return nil
	}

	res := props.Resource
	if res == nil {
		return nil
	}

	return res.Throughput
}

func ConvertThroughputFromResourceData(throughput interface{}) *int32 {
	return utils.Int32(int32(throughput.(int)))
}
