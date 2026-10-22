// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/publicmaintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databasesecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/geobackuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/transparentdataencryptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMssqlDatabaseSetFlatten(d *pluginsdk.ResourceData, id *commonids.SqlDatabaseId, model *databases.Database, metaClient *clients.Client) error {
	securityAlertPoliciesClient := metaClient.MSSQL.DatabaseSecurityAlertPoliciesClient

	longTermRetentionClient := metaClient.MSSQL.LongTermRetentionPoliciesClient
	shortTermRetentionClient := metaClient.MSSQL.BackupShortTermRetentionPoliciesClient
	geoBackupPoliciesClient := metaClient.MSSQL.GeoBackupPoliciesClient
	transparentEncryptionClient := metaClient.MSSQL.TransparentDataEncryptionsClient
	ctx, cancel := timeouts.ForRead(metaClient.StopContext, d)
	defer cancel()

	geoBackupPolicy := true
	skuName := ""
	elasticPoolId := ""
	ledgerEnabled := false
	enclaveType := ""

	if model != nil {
		d.Set("name", id.DatabaseName)

		if props := model.Properties; props != nil {
			requestedBackupStorageRedundancy := ""
			if props.RequestedBackupStorageRedundancy != nil {
				requestedBackupStorageRedundancy = string(*props.RequestedBackupStorageRedundancy)
			}

			// A named replica doesn't return props.RequestedBackupStorageRedundancy from the api but it is Geo in the portal regardless of what the parent database has
			// so we'll copy that here to get around a perpetual diff
			if props.SecondaryType != nil && *props.SecondaryType == "Named" {
				requestedBackupStorageRedundancy = string(databases.BackupStorageRedundancyGeo)
			}

			d.Set("auto_pause_delay_in_minutes", pointer.From(props.AutoPauseDelay))
			d.Set("collation", pointer.From(props.Collation))
			d.Set("read_replica_count", pointer.From(props.HighAvailabilityReplicaCount))
			d.Set("storage_account_type", requestedBackupStorageRedundancy)
			d.Set("zone_redundant", pointer.From(props.ZoneRedundant))
			d.Set("read_scale", pointer.From(props.ReadScale) == databases.DatabaseReadScaleEnabled)
			d.Set("secondary_type", pointer.From(props.SecondaryType))

			if props.ElasticPoolId != nil {
				elasticPoolId = pointer.From(props.ElasticPoolId)
			}

			if props.LicenseType != nil {
				d.Set("license_type", pointer.FromEnum(props.LicenseType))
			} else {
				// value not returned, try to set from state
				d.Set("license_type", d.Get("license_type").(string))
			}

			// when `max_size_gb` is below 1GB, the API only accepts 100MB and 500MB. 100MB is 104857600, and 500MB is 524288000.
			// directly set `max_size_gb` when API returns `104857600` or `524288000`.
			size := float64(pointer.From(props.MaxSizeBytes) / int64(1073741824))
			switch pointer.From(props.MaxSizeBytes) {
			case 104857600:
				size = 0.1
			case 524288000:
				size = 0.5
			}
			d.Set("max_size_gb", size)

			if props.CurrentServiceObjectiveName != nil {
				skuName = *props.CurrentServiceObjectiveName
			}

			if props.IsLedgerOn != nil {
				ledgerEnabled = *props.IsLedgerOn
			}

			// NOTE: Always set the PreferredEnclaveType to an empty string
			// if not in the properties that were returned from Azure...
			if v := props.PreferredEnclaveType; v != nil {
				enclaveType = pointer.FromEnum(v)
			}

			configurationName := ""
			if v := props.MaintenanceConfigurationId; v != nil {
				maintenanceConfigId, err := publicmaintenanceconfigurations.ParsePublicMaintenanceConfigurationIDInsensitively(pointer.From(v))
				if err != nil {
					return err
				}
				configurationName = maintenanceConfigId.PublicMaintenanceConfigurationName
			}

			d.Set("elastic_pool_id", elasticPoolId)
			d.Set("min_capacity", pointer.From(props.MinCapacity))
			d.Set("sku_name", skuName)
			d.Set("maintenance_configuration_name", configurationName)
			d.Set("ledger_enabled", ledgerEnabled)
			d.Set("enclave_type", enclaveType)
			d.Set("transparent_data_encryption_key_vault_key_id", props.EncryptionProtector)
			d.Set("transparent_data_encryption_key_automatic_rotation_enabled", pointer.From(props.EncryptionProtectorAutoRotation))

			identity, err := identity.FlattenUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if err := d.Set("identity", identity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return err
			}
		}

		// Determine whether the SKU is for SQL Data Warehouse
		isDwSku := strings.HasPrefix(strings.ToLower(skuName), "dw")

		// Determine whether the SKU is for SQL Database Free tier
		isFreeSku := strings.EqualFold(skuName, "free")

		// DW SKUs and SQL Database Free tier do not currently support LRP and do not honour normal SRP operations
		if !isDwSku && !isFreeSku {
			longTermPolicy, err := longTermRetentionClient.Get(ctx, pointer.From(id))
			if err != nil {
				return fmt.Errorf("retrieving Long Term Retention Policies for %s: %+v", id, err)
			}

			if longTermPolicyModel := longTermPolicy.Model; longTermPolicyModel != nil {
				if err := d.Set("long_term_retention_policy", helper.FlattenLongTermRetentionPolicy(longTermPolicyModel)); err != nil {
					return fmt.Errorf("setting `long_term_retention_policy`: %+v", err)
				}
			}

			shortTermPolicy, err := shortTermRetentionClient.Get(ctx, pointer.From(id))
			if err != nil {
				return fmt.Errorf("retrieving Short Term Retention Policies for %s: %+v", id, err)
			}

			if shortTermPolicyModel := shortTermPolicy.Model; shortTermPolicyModel != nil {
				if err := d.Set("short_term_retention_policy", helper.FlattenShortTermRetentionPolicy(shortTermPolicyModel)); err != nil {
					return fmt.Errorf("setting `short_term_retention_policy`: %+v", err)
				}
			}
		} else {
			// DW SKUs and SQL Database Free tier need the retention policies to be empty for state consistency
			emptySlice := make([]any, 0)
			d.Set("long_term_retention_policy", emptySlice)
			d.Set("short_term_retention_policy", emptySlice)

			geoPoliciesResponse, err := geoBackupPoliciesClient.Get(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(geoPoliciesResponse.HttpResponse) {
					d.SetId("")
					return nil
				}

				return fmt.Errorf("retrieving Geo Backup Policies for %s: %+v", id, err)
			}

			// For Datawarehouse SKUs and SQL Database Free tier, set the geo-backup policy setting
			if geoPolicyModel := geoPoliciesResponse.Model; geoPolicyModel != nil {
				if (isDwSku || isFreeSku) && geoPolicyModel.Properties.State == geobackuppolicies.GeoBackupPolicyStateDisabled {
					geoBackupPolicy = false
				}
			}
		}
	}

	if err := d.Set("geo_backup_enabled", geoBackupPolicy); err != nil {
		return fmt.Errorf("setting `geo_backup_enabled`: %+v", err)
	}

	securityAlertPolicy, err := securityAlertPoliciesClient.Get(ctx, pointer.From(id))
	if err == nil && securityAlertPolicy.Model != nil {
		if err := d.Set("threat_detection_policy", flattenMsSqlServerSecurityAlertPolicy(d, pointer.From(securityAlertPolicy.Model))); err != nil {
			return fmt.Errorf("setting `threat_detection_policy`: %+v", err)
		}
	}

	tde, err := transparentEncryptionClient.Get(ctx, pointer.From(id))
	if err != nil {
		return fmt.Errorf("while retrieving Transparent Data Encryption state for %s: %+v", id, err)
	}

	tdeState := false
	if model := tde.Model; model != nil {
		if props := model.Properties; props != nil {
			tdeState = (props.State == transparentdataencryptions.TransparentDataEncryptionStateEnabled)
		}
	}
	d.Set("transparent_data_encryption_enabled", tdeState)

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenMsSqlServerSecurityAlertPolicy(d *pluginsdk.ResourceData, policy databasesecurityalertpolicies.DatabaseSecurityAlertPolicy) []any {
	// The SQL database security alert API always returns the default value even if never set.
	// If the values are on their default one, threat it as not set.
	properties := policy.Properties
	if properties == nil {
		return []any{}
	}

	securityAlertPolicy := map[string]any{
		"state":                        string(properties.State),
		"email_account_admins_enabled": pointer.From(properties.EmailAccountAdmins),
	}

	if disabledAlerts := properties.DisabledAlerts; disabledAlerts != nil {
		flattenedAlerts := pluginsdk.NewSet(pluginsdk.HashString, []any{})
		for _, a := range *disabledAlerts {
			if a != "" {
				flattenedAlerts.Add(a)
			}
		}
		securityAlertPolicy["disabled_alerts"] = flattenedAlerts
	}
	if emailAddresses := properties.EmailAddresses; emailAddresses != nil {
		flattenedEmails := pluginsdk.NewSet(pluginsdk.HashString, []any{})
		for _, e := range *emailAddresses {
			if e != "" {
				flattenedEmails.Add(e)
			}
		}
		securityAlertPolicy["email_addresses"] = flattenedEmails
	}
	if properties.StorageEndpoint != nil {
		securityAlertPolicy["storage_endpoint"] = *properties.StorageEndpoint
	}
	if properties.RetentionDays != nil {
		securityAlertPolicy["retention_days"] = int(*properties.RetentionDays)
	}

	// If storage account access key is in state read it to the new state, as the API does not return it for security reasons
	if v, ok := d.GetOk("threat_detection_policy.0.storage_account_access_key"); ok {
		securityAlertPolicy["storage_account_access_key"] = v.(string)
	}

	return []any{securityAlertPolicy}
}
