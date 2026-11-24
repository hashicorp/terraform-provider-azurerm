// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_clone_from_database

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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

		// Clone from database specific required fields

		"source_autonomous_database_id": commonschema.ResourceIDReferenceRequiredForceNew(&autonomousdatabases.AutonomousDatabaseId{}),

		"clone_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(autonomousdatabases.PossibleValuesForCloneType(), false),
		},
		//  clone from database specific optional fields

		"refreshable_model": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(autonomousdatabases.PossibleValuesForRefreshableModelType(), false),
		},

		// Required (inherited from base)

		"admin_password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ForceNew:     true,
			ValidateFunc: validate.AutonomousDatabasePassword,
		},

		"backup_retention_period_in_days": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 60),
		},

		"character_set": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"compute_count": {
			Type:         pluginsdk.TypeFloat,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.FloatBetween(2.0, 512.0),
		},

		"compute_model": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AdbsComputeModel,
		},

		"data_storage_size_in_tb": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 384),
		},

		"database_version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"database_workload": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(autonomousdatabases.PossibleValuesForWorkloadType(), false),
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 255),
		},

		"auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
			ForceNew: true,
		},

		"auto_scaling_for_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
			ForceNew: true,
		},

		"mtls_connection_required": {
			Type:     pluginsdk.TypeBool,
			Required: true,
			ForceNew: true,
		},

		"license_model": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(autonomousdatabases.PossibleValuesForLicenseModel(), false),
		},

		"national_character_set": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Optional
		"customer_contacts": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsEmailAddress,
			},
		},
		"allowed_ip_addresses": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			MaxItems:     1024,
			ForceNew:     true,
			ExactlyOneOf: []string{"allowed_ip_addresses", "virtual_network_id"},
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.Any(
					validation.IsIPv4Address,
					validation.IsCIDR,
				),
			},
		},

		"subnet_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"virtual_network_id"},
		},

		"virtual_network_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"allowed_ip_addresses", "virtual_network_id"},
			RequiredWith: []string{"subnet_id"},
		},

		"tags": commonschema.TagsForceNew(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (Resource) CustomizeDiff() sdk.ResourceFunc {
	workloadMatrixForClone := map[string][]string{
		"DW":   {"OLTP", "DW"},
		"OLTP": {"DW", "OLTP"},
		"AJD":  {"OLTP", "DW", "APEX"},
		"APEX": {"AJD", "OLTP", "DW"},
	}
	return sdk.ResourceFunc{
		Timeout: 30 * time.Second,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			sourceId := metadata.ResourceDiff.Get("source_autonomous_database_id").(string)
			dbWorkload := metadata.ResourceDiff.Get("database_workload").(string)

			if sourceId == "" || dbWorkload == "" {
				return nil
			}

			if metadata.ResourceDiff.Id() != "" {
				return nil
			}

			sourceWorkload, err := getSourceWorkloadforClone(ctx, sourceId, metadata)
			if err != nil {
				return err
			}

			targets, exists := workloadMatrixForClone[sourceWorkload]
			if !exists {
				return fmt.Errorf("unsupported source workload: %s", sourceWorkload)
			}

			if slices.Contains(targets, dbWorkload) {
				return nil
			}

			return fmt.Errorf("invalid workload: %s->%s not allowed", sourceWorkload, dbWorkload)
		},
	}
}
