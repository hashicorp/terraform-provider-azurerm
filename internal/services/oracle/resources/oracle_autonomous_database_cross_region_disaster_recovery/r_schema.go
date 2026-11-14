// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_cross_region_disaster_recovery

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
			ForceNew:     true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},

		"source_autonomous_database_id": commonschema.ResourceIDReferenceRequiredForceNew(&autonomousdatabases.AutonomousDatabaseId{}),

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"replicate_automatic_backups_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"tags": commonschema.TagsForceNew(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"auto_scaling_for_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"backup_retention_period_in_days": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"character_set": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"compute_count": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},
		"compute_model": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"customer_contacts": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
		"data_storage_size_in_tb": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"database_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"database_workload": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"license_model": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"mtls_connection_required": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"national_character_set": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"remote_disaster_recovery_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
