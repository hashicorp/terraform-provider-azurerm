// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database_extended_auditing_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlDatabaseExtendedAuditingPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSqlDatabaseID,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"blob_storage_endpoint": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"storage_account_access_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_account_access_key_is_secondary": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"retention_in_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntBetween(0, 3285),
		},

		"log_monitoring_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
	}
}
