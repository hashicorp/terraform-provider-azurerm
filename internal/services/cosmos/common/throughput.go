// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

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

func ConvertThroughputFromResourceData(throughput interface{}) *int64 {
	return pointer.To(int64(throughput.(int)))
}

func ExpandCosmosDBThroughputSettingsUpdateParameters(d *pluginsdk.ResourceData) cosmosdb.ThroughputSettingsUpdateParameters {
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

	return throughputParameters
}

func SetResourceDataThroughputFromResponse(throughputResponse cosmosdb.ThroughputSettingsGetResults, d *pluginsdk.ResourceData) {
	d.Set("throughput", GetThroughputFromResult(throughputResponse))

	d.Set("autoscale_settings", FlattenCosmosDbAutoscaleSettings(throughputResponse))
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

type ThroughputMode string

const (
	ThroughputModeNone      ThroughputMode = "none"
	ThroughputModeManual    ThroughputMode = "manual"
	ThroughputModeAutoscale ThroughputMode = "autoscale"
)

// DesiredThroughputMode determines which throughput mode is set in the configuration.
//
// `throughput` and `autoscale_settings.0.max_throughput` are both Optional+Computed, so `d.Get`
// returns the value last read back from the API even once the user has removed it from their
// configuration. The raw configuration is therefore the only reliable signal of intent here.
func DesiredThroughputMode(d *pluginsdk.ResourceData) ThroughputMode {
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return ThroughputModeNone
	}

	configMap := rawConfig.AsValueMap()

	if v, ok := configMap["autoscale_settings"]; ok && !v.IsNull() && v.IsKnown() && len(v.AsValueSlice()) > 0 {
		return ThroughputModeAutoscale
	}

	if v, ok := configMap["throughput"]; ok && !v.IsNull() {
		return ThroughputModeManual
	}

	return ThroughputModeNone
}

// CurrentThroughputMode determines the throughput mode the resource is currently provisioned with.
//
// An autoscale offer reports both `throughput` (the RU/s it is currently scaled to) and
// `autoscaleSettings.maxThroughput`, so the presence of `autoscaleSettings` is what distinguishes
// the two modes.
func CurrentThroughputMode(input cosmosdb.ThroughputSettingsGetResults) ThroughputMode {
	if input.Properties == nil || input.Properties.Resource == nil {
		return ThroughputModeNone
	}

	res := input.Properties.Resource

	if res.AutoScaleSettings != nil {
		return ThroughputModeAutoscale
	}

	if res.Throughput != nil {
		return ThroughputModeManual
	}

	return ThroughputModeNone
}

// ThroughputMigrationRequired reports whether the resource has to be migrated between throughput
// modes before the value in the configuration can be applied.
func ThroughputMigrationRequired(current ThroughputMode, desired ThroughputMode) bool {
	if current == ThroughputModeNone || desired == ThroughputModeNone {
		return false
	}

	return current != desired
}

// ThroughputValueMatchesConfig reports whether the throughput value returned by the API already
// matches the value in the configuration.
//
// The migration APIs accept no request body and assign a system-determined value, so this is used
// to decide whether a follow-up update is needed to reach the configured value.
func ThroughputValueMatchesConfig(d *pluginsdk.ResourceData, input cosmosdb.ThroughputSettingsGetResults) bool {
	if input.Properties == nil || input.Properties.Resource == nil {
		return false
	}

	res := input.Properties.Resource

	switch DesiredThroughputMode(d) {
	case ThroughputModeAutoscale:
		if res.AutoScaleSettings == nil {
			return false
		}

		return res.AutoScaleSettings.MaxThroughput == int64(d.Get("autoscale_settings.0.max_throughput").(int))

	case ThroughputModeManual:
		if res.Throughput == nil {
			return false
		}

		return pointer.From(res.Throughput) == int64(d.Get("throughput").(int))
	}

	return false
}
