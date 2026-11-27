// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func recoveryServicesVaultSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Recovery Services Vault. Recovery Service Vault name must be 2 - 50 characters long, start with a letter, contain only letters, numbers and hyphens.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"encryption": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			RequiredWith: []string{"identity"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_id": {
						Description:  "The Key Vault Key ID used to encrypt this vault. Key managed by Vault Managed Hardware Security Module is also supported.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
					},
					"infrastructure_encryption_enabled": {
						Description: "Enabling/Disabling the Double Encryption state.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
					"use_system_assigned_identity": {
						Description: "Indicate that System Assigned Identity should be used or not.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"user_assigned_identity_id": {
						Description:  "The User Assigned Identity ID to be used.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"public_network_access_enabled": {
			Description: "Is it enabled to access the vault from public networks?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		// set `immutability` to Computed, because it will start to return from the service once it has been set.
		"immutability": {
			Description:  "Immutability Settings of vault.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(vaults.PossibleValuesForImmutabilityState(), false),
		},

		"tags": commonschema.Tags(),

		"sku": {
			Description:  "Sets the vault's SKU.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(vaults.PossibleValuesForSkuName(), false),
		},

		"storage_mode_type": {
			Description: "The storage type of the Recovery Services Vault.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     vaults.StandardTierStorageRedundancyGeoRedundant,
			ValidateFunc: validation.StringInSlice([]string{
				string(vaults.StandardTierStorageRedundancyGeoRedundant),
				string(vaults.StandardTierStorageRedundancyLocallyRedundant),
				string(vaults.StandardTierStorageRedundancyZoneRedundant),
			}, false),
		},

		"cross_region_restore_enabled": {
			Description: "Is cross region restore enabled for this Vault? Can only be `true`, when `storage_mode_type` is `GeoRedundant`.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"monitoring": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"alerts_for_all_job_failures_enabled": {
						Description: "Enabling/Disabling built-in Azure Monitor alerts for security scenarios and job failure scenarios.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"alerts_for_all_failover_issues_enabled": {
						Description: "Enabling/Disabling built-in Azure Monitor alerts for all failover issues.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"alerts_for_all_replication_issues_enabled": {
						Description: "Enabling/Disabling built-in Azure Monitor alerts for all replication issues.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"alerts_for_critical_operation_failures_enabled": {
						Description: "Enabling/Disabling alerts from the older (classic alerts) solution. More details could be found [here](https://learn.microsoft.com/azure/backup/monitoring-and-alerts-overview).",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"email_notifications_for_site_recovery_enabled": {
						Description: "Enabling/Disabling email notifications for Site Recovery (classic alerts) solution.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
				},
			},
		},

		"classic_vmware_replication_enabled": {
			Description: "Whether to enable the Classic experience for VMware replication. If set to `false` VMware machines will be protected using the new stateless ASR replication appliance.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},
	}
}
