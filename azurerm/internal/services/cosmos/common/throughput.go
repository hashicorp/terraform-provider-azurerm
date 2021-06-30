package common

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
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

func ExpandCosmosDBThroughputSettingsUpdateParameters(d *pluginsdk.ResourceData) *documentdb.ThroughputSettingsUpdateParameters {
	throughputParameters := documentdb.ThroughputSettingsUpdateParameters{
		ThroughputSettingsUpdateProperties: &documentdb.ThroughputSettingsUpdateProperties{
			Resource: &documentdb.ThroughputSettingsResource{},
		},
	}

	if v, exists := d.GetOk("throughput"); exists {
		throughputParameters.ThroughputSettingsUpdateProperties.Resource.Throughput = ConvertThroughputFromResourceData(v)
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		// If updating the autoscale throughput, set the manual throughput to nil to ensure the autoscale throughput is applied
		throughputParameters.ThroughputSettingsUpdateProperties.Resource.Throughput = nil
		throughputParameters.ThroughputSettingsUpdateProperties.Resource.AutoscaleSettings = ExpandCosmosDbAutoscaleSettingsResource(d)
	}

	return &throughputParameters
}

func SetResourceDataThroughputFromResponse(throughputResponse documentdb.ThroughputSettingsGetResults, d *pluginsdk.ResourceData) {
	d.Set("throughput", GetThroughputFromResult(throughputResponse))

	autoscaleSettings := FlattenCosmosDbAutoscaleSettings(throughputResponse)
	d.Set("autoscale_settings", autoscaleSettings)
}

func CheckForChangeFromAutoscaleAndManualThroughput(d *pluginsdk.ResourceData) error {
	if d.HasChange("throughput") && d.HasChange("autoscale_settings") {
		return fmt.Errorf("switching between autoscale and manually provisioned throughput via Terraform is not supported at this time")
	}

	return nil
}

func HasThroughputChange(d *pluginsdk.ResourceData) bool {
	return d.HasChanges("throughput", "autoscale_settings")
}
