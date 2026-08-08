// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2026-03-15/openapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func GetThroughputFromResult(throughputResponse openapis.ThroughputSettingsGetResults) *int64 {
	props := throughputResponse.Properties
	if props == nil {
		return nil
	}

	res := props.Resource
	if res == nil {
		return nil
	}

	return res.Throughput
}

func ConvertThroughputFromResourceData(throughput interface{}) *int64 {
	return pointer.To(int64(throughput.(int)))
}

func ExpandCosmosDBThroughputSettingsUpdateParameters(d *pluginsdk.ResourceData) openapis.ThroughputSettingsUpdateParameters {
	throughputParameters := openapis.ThroughputSettingsUpdateParameters{
		Properties: openapis.ThroughputSettingsUpdateProperties{
			Resource: openapis.ThroughputSettingsResource{},
		},
	}

	if v, exists := d.GetOk("throughput"); exists {
		throughputParameters.Properties.Resource.Throughput = ConvertThroughputFromResourceData(v)
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		// If updating the autoscale throughput, set the manual throughput to nil to ensure the autoscale throughput is applied
		throughputParameters.Properties.Resource.Throughput = nil
		throughputParameters.Properties.Resource.AutoScaleSettings = ExpandCosmosDbAutoscaleSettingsResource(d)
	}

	return throughputParameters
}

func SetResourceDataThroughputFromResponse(throughputResponse openapis.ThroughputSettingsGetResults, d *pluginsdk.ResourceData) {
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
