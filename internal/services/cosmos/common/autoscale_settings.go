// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func ExpandCosmosDbAutoscaleSettingsLegacy(d *pluginsdk.ResourceData) *documentdb.AutoscaleSettings {
	i := d.Get("autoscale_settings").([]interface{})
	if len(i) == 0 || i[0] == nil {
		log.Printf("[DEBUG] Cosmos DB autoscale settings are not set on the resource")
		return nil
	}
	input := i[0].(map[string]interface{})

	autoscaleSettings := documentdb.AutoscaleSettings{}

	if maxThroughput, ok := input["max_throughput"].(int); ok {
		autoscaleSettings.MaxThroughput = utils.Int32(int32(maxThroughput))
	}

	return &autoscaleSettings
}

func ExpandCosmosDbAutoscaleSettings(d *pluginsdk.ResourceData) *cosmosdb.AutoScaleSettings {
	i := d.Get("autoscale_settings").([]interface{})
	if len(i) == 0 || i[0] == nil {
		log.Printf("[DEBUG] Cosmos DB autoscale settings are not set on the resource")
		return nil
	}
	input := i[0].(map[string]interface{})

	autoscaleSettings := cosmosdb.AutoScaleSettings{}

	if maxThroughput, ok := input["max_throughput"].(int); ok {
		autoscaleSettings.MaxThroughput = utils.Int64(int64(maxThroughput))
	}

	return &autoscaleSettings
}

func FlattenCosmosDbAutoscaleSettingsLegacy(throughputResponse documentdb.ThroughputSettingsGetResults) []interface{} {
	results := make([]interface{}, 0)

	props := throughputResponse.ThroughputSettingsGetProperties
	if props == nil {
		return results
	}

	res := props.Resource
	if res == nil {
		return results
	}

	autoscaleSettings := res.AutoscaleSettings
	if autoscaleSettings == nil {
		log.Printf("[DEBUG] Cosmos DB autoscale settings are not set on the throughput response")
		return results
	}

	result := make(map[string]interface{})

	if autoscaleSettings.MaxThroughput != nil {
		result["max_throughput"] = autoscaleSettings.MaxThroughput
	}

	return append(results, result)
}

func FlattenCosmosDbAutoscaleSettings(throughputResponse cosmosdb.ThroughputSettingsGetResults) []interface{} {
	results := make([]interface{}, 0)

	props := throughputResponse.Properties
	if props == nil {
		return results
	}

	res := props.Resource
	if res == nil {
		return results
	}

	autoscaleSettings := res.AutoScaleSettings
	if autoscaleSettings == nil {
		log.Printf("[DEBUG] Cosmos DB autoscale settings are not set on the throughput response")
		return results
	}

	result := make(map[string]interface{})

	result["max_throughput"] = autoscaleSettings.MaxThroughput

	return append(results, result)
}

func ExpandCosmosDbAutoscaleSettingsResourceLegacy(d *pluginsdk.ResourceData) *documentdb.AutoscaleSettingsResource {
	autoscaleSettings := ExpandCosmosDbAutoscaleSettingsLegacy(d)

	if autoscaleSettings == nil {
		return nil
	}

	return &documentdb.AutoscaleSettingsResource{
		MaxThroughput: autoscaleSettings.MaxThroughput,
	}
}

func ExpandCosmosDbAutoscaleSettingsResource(d *pluginsdk.ResourceData) *cosmosdb.AutoscaleSettingsResource {
	autoscaleSettings := ExpandCosmosDbAutoscaleSettings(d)

	if autoscaleSettings == nil {
		return nil
	}

	return &cosmosdb.AutoscaleSettingsResource{
		MaxThroughput: *autoscaleSettings.MaxThroughput,
	}
}
