package common

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2020-04-01/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func ExpandCosmosDBThroughputSettingsUpdateParameters(d *schema.ResourceData) *documentdb.ThroughputSettingsUpdateParameters {
	throughputParameters := documentdb.ThroughputSettingsUpdateParameters{
		ThroughputSettingsUpdateProperties: &documentdb.ThroughputSettingsUpdateProperties{
			Resource: &documentdb.ThroughputSettingsResource{},
		},
	}

	if d.Get("throughput").(int) != 0 {
		throughputParameters.ThroughputSettingsUpdateProperties.Resource.Throughput = ConvertThroughputFromResourceData(d.Get("throughput"))
	}

	if d.HasChange("autoscale_settings") && d.Get("throughput").(int) == 0 {
		log.Printf("[DEBUG] Cosmos DB autoscale settings have changed")
		throughputParameters.ThroughputSettingsUpdateProperties.Resource.Throughput = nil
		throughputParameters.ThroughputSettingsUpdateProperties.Resource.AutoscaleSettings = ExpandCosmosDbAutoscaleSettingsResource(d)
	}

	return &throughputParameters
}

func SetResourceDataThroughputFromResponse(throughputResponse documentdb.ThroughputSettingsGetResults, d *schema.ResourceData) {
	d.Set("throughput", GetThroughputFromResult(throughputResponse))

	autoscaleSettings := FlattenCosmosDbAutoscaleSettings(throughputResponse)
	d.Set("autoscale_settings", autoscaleSettings)
	if len(autoscaleSettings) != 0 {
		d.Set("throughput", nil)
	}
}

func CheckForChangeFromAutoscaleAndManualThroughput(d *schema.ResourceData) error {
	if d.HasChange("throughput") && d.HasChange("autoscale_settings") {
		return fmt.Errorf("Switching between autoscale and manual provisioned throughput is not supported at this time.")
	}

	return nil
}

func HasThroughputChange(d *schema.ResourceData) bool {
	return d.HasChanges("throughput", "autoscale_settings")
}
