// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_elasticpool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/publicmaintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/elasticpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceMssqlElasticPoolSetFlatten(d *pluginsdk.ResourceData, id *commonids.SqlElasticPoolId, model *elasticpools.ElasticPool) error {
	if model != nil {
		d.Set("name", model.Name)
		d.Set("resource_group_name", id.ResourceGroupName)
		d.Set("location", model.Location)
		d.Set("server_name", id.ServerName)

		if err := d.Set("sku", flattenMsSqlElasticPoolSku(model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		if props := model.Properties; props != nil {
			enclaveType := ""
			if v := props.PreferredEnclaveType; v != nil {
				enclaveType = pointer.FromEnum(v)
			}
			d.Set("enclave_type", enclaveType)
			d.Set("max_size_gb", pointer.To(float64(*props.MaxSizeBytes)/1073741824))
			d.Set("max_size_bytes", pointer.To(props.MaxSizeBytes))
			d.Set("zone_redundant", pointer.From(props.ZoneRedundant))

			licenseType := string(elasticpools.ElasticPoolLicenseTypeLicenseIncluded)
			if props.LicenseType != nil {
				licenseType = string(*props.LicenseType)
			}
			d.Set("license_type", licenseType)
			d.Set("high_availability_replica_count", pointer.From(props.HighAvailabilityReplicaCount))

			if err := d.Set("per_database_settings", flattenMsSqlElasticPoolPerDatabaseSettings(props.PerDatabaseSettings)); err != nil {
				return fmt.Errorf("setting `per_database_settings`: %+v", err)
			}

			maintenanceConfigId, err := publicmaintenanceconfigurations.ParsePublicMaintenanceConfigurationIDInsensitively(*props.MaintenanceConfigurationId)
			if err != nil {
				return err
			}
			d.Set("maintenance_configuration_name", maintenanceConfigId.PublicMaintenanceConfigurationName)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenMsSqlElasticPoolSku(input *elasticpools.Sku) []any {
	if input == nil {
		return []any{}
	}

	values := map[string]any{}

	if name := input.Name; name != "" {
		values["name"] = name
	}

	if tier := input.Tier; tier != nil {
		values["tier"] = *tier
	}

	if family := input.Family; family != nil {
		values["family"] = *family
	}

	if capacity := input.Capacity; capacity != nil {
		values["capacity"] = *capacity
	}

	return []any{values}
}

func flattenMsSqlElasticPoolPerDatabaseSettings(resp *elasticpools.ElasticPoolPerDatabaseSettings) []any {
	perDatabaseSettings := map[string]any{}

	if minCapacity := resp.MinCapacity; minCapacity != nil {
		perDatabaseSettings["min_capacity"] = *minCapacity
	}

	if maxCapacity := resp.MaxCapacity; maxCapacity != nil {
		perDatabaseSettings["max_capacity"] = *maxCapacity
	}

	return []any{perDatabaseSettings}
}
