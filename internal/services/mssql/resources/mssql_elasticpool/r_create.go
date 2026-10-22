// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_elasticpool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/publicmaintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/elasticpools"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlElasticPoolCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSqlElasticPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_mssql_elasticpool", id.ID())
			}
		}
	}

	location := location.Normalize(d.Get("location").(string))
	sku := expandMsSqlElasticPoolSku(d)

	maintenanceConfigId := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(subscriptionId, d.Get("maintenance_configuration_name").(string))
	elasticPool := elasticpools.ElasticPool{
		Name:     pointer.To(id.ElasticPoolName),
		Location: location,
		Sku:      sku,
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
		Properties: &elasticpools.ElasticPoolProperties{
			LicenseType:                pointer.ToEnum[elasticpools.ElasticPoolLicenseType](d.Get("license_type").(string)),
			PerDatabaseSettings:        expandMsSqlElasticPoolPerDatabaseSettings(d),
			ZoneRedundant:              pointer.To(d.Get("zone_redundant").(bool)),
			MaintenanceConfigurationId: pointer.To(maintenanceConfigId.ID()),
		},
	}

	if !d.GetRawConfig().AsValueMap()["high_availability_replica_count"].IsNull() {
		elasticPool.Properties.HighAvailabilityReplicaCount = pointer.To(int64(d.Get("high_availability_replica_count").(int)))
	}

	// NOTE: The service default is actually nil/empty which indicates enclave is disabled. the value `Default` is NOT the default.
	if v, ok := d.GetOk("enclave_type"); ok && v.(string) != "" {
		elasticPool.Properties.PreferredEnclaveType = pointer.ToEnum[elasticpools.AlwaysEncryptedEnclaveType](v.(string))
	}

	if d.HasChange("max_size_gb") {
		if v, ok := d.GetOk("max_size_gb"); ok {
			maxSizeBytes := v.(float64) * 1073741824
			elasticPool.Properties.MaxSizeBytes = pointer.To(int64(maxSizeBytes))
		}
	} else if v, ok := d.GetOk("max_size_bytes"); ok {
		elasticPool.Properties.MaxSizeBytes = pointer.To(int64(v.(int)))
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, elasticPool, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, elasticPool); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceMsSqlElasticPoolRead(d, meta)
}

func expandMsSqlElasticPoolPerDatabaseSettings(d *pluginsdk.ResourceData) *elasticpools.ElasticPoolPerDatabaseSettings {
	perDatabaseSettings := d.Get("per_database_settings").([]any)
	perDatabaseSetting := perDatabaseSettings[0].(map[string]any)

	minCapacity := perDatabaseSetting["min_capacity"].(float64)
	maxCapacity := perDatabaseSetting["max_capacity"].(float64)

	return &elasticpools.ElasticPoolPerDatabaseSettings{
		MinCapacity: pointer.To(minCapacity),
		MaxCapacity: pointer.To(maxCapacity),
	}
}

func expandMsSqlElasticPoolSku(d *pluginsdk.ResourceData) *elasticpools.Sku {
	skus := d.Get("sku").([]any)
	sku := skus[0].(map[string]any)

	name := sku["name"].(string)
	tier := sku["tier"].(string)
	family := sku["family"].(string)
	capacity := sku["capacity"].(int)

	return &elasticpools.Sku{
		Name:     name,
		Tier:     pointer.To(tier),
		Family:   pointer.To(family),
		Capacity: pointer.To(int64(capacity)),
	}
}
