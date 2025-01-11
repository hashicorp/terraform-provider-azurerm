// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func GetThroughputFromResultLegacy(throughputResponse documentdb.ThroughputSettingsGetResults) *int32 {
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

func GetThroughputFromResult(throughputResponse cosmosdb.ThroughputSettingsGetResults) *int64 {
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

func ConvertThroughputFromResourceDataLegacy(throughput interface{}) *int32 {
	return utils.Int32(int32(throughput.(int)))
}

func ConvertThroughputFromResourceData(throughput interface{}) *int64 {
	return utils.Int64(int64(throughput.(int)))
}

func ExpandCosmosDBThroughputSettingsUpdateParametersLegacy(d *pluginsdk.ResourceData) *documentdb.ThroughputSettingsUpdateParameters {
	throughputParameters := documentdb.ThroughputSettingsUpdateParameters{
		ThroughputSettingsUpdateProperties: &documentdb.ThroughputSettingsUpdateProperties{
			Resource: &documentdb.ThroughputSettingsResource{},
		},
	}

	if v, exists := d.GetOk("throughput"); exists {
		throughputParameters.ThroughputSettingsUpdateProperties.Resource.Throughput = ConvertThroughputFromResourceDataLegacy(v)
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		// If updating the autoscale throughput, set the manual throughput to nil to ensure the autoscale throughput is applied
		throughputParameters.ThroughputSettingsUpdateProperties.Resource.Throughput = nil
		throughputParameters.ThroughputSettingsUpdateProperties.Resource.AutoscaleSettings = ExpandCosmosDbAutoscaleSettingsResourceLegacy(d)
	}

	return &throughputParameters
}

func ExpandCosmosDBThroughputSettingsUpdateParameters(d *pluginsdk.ResourceData) *cosmosdb.ThroughputSettingsUpdateParameters {
	throughputParameters := cosmosdb.ThroughputSettingsUpdateParameters{
		Properties: cosmosdb.ThroughputSettingsUpdateProperties{
			Resource: cosmosdb.ThroughputSettingsResource{},
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

	return &throughputParameters
}

func SetResourceDataThroughputFromResponseLegacy(throughputResponse documentdb.ThroughputSettingsGetResults, d *pluginsdk.ResourceData) {
	d.Set("throughput", GetThroughputFromResultLegacy(throughputResponse))

	autoscaleSettings := FlattenCosmosDbAutoscaleSettingsLegacy(throughputResponse)
	d.Set("autoscale_settings", autoscaleSettings)
}

func SetResourceDataThroughputFromResponse(throughputResponse cosmosdb.ThroughputSettingsGetResults, d *pluginsdk.ResourceData) {
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
