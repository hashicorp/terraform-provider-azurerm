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
			Description:  "The name for this Cross Region Disaster Recovery Autonomous Database.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"display_name": {
			Description:  "The user-friendly name for the Autonomous Database.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},

		"source_autonomous_database_id": commonschema.ResourceIDReferenceRequiredForceNew(&autonomousdatabases.AutonomousDatabaseId{}, "source (primary) Autonomous Database"),

		"subnet_id": {
			Description:  "The ID of the subnet in the target region.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"replicate_automatic_backups_enabled": {
			Description: "If true, 7 days of backups are replicated across regions.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"tags": commonschema.TagsForceNew(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"auto_scaling_enabled": {
			Description: "Whether auto-scaling is enabled for the Autonomous Database CPU core count.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},
		"auto_scaling_for_storage_enabled": {
			Description: "Whether auto-scaling is enabled for the Autonomous Database storage.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},
		"backup_retention_period_in_days": {
			Description: "The backup retention period in days.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},
		"character_set": {
			Description: "The character set for the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"compute_count": {
			Description: "The compute amount available to the database.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},
		"compute_model": {
			Description: "The compute model of the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"customer_contacts": {
			Description: "A list of customer contact email addresses.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
		"data_storage_size_in_tb": {
			Description: "The maximum storage that can be allocated for the database in terabytes.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},
		"database_version": {
			Description: "The Oracle Database version for the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"database_workload": {
			Description: "The Autonomous Database workload type.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"license_model": {
			Description: "The Oracle license model that applies to the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"mtls_connection_required": {
			Description: "Whether the Autonomous Database requires mTLS connections.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},
		"national_character_set": {
			Description: "The national character set for the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"remote_disaster_recovery_type": {
			Description: "Type of recovery.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
