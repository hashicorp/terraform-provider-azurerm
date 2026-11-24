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
			Description:  "The ID of the SQL database to set the extended auditing policy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSqlDatabaseID,
		},

		"enabled": {
			Description: "Whether to enable the extended auditing policy.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"blob_storage_endpoint": {
			Description:  "The blob storage endpoint (e.g. <https://example.blob.core.windows.net>). This blob storage will hold all extended auditing logs.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"storage_account_access_key": {
			Description:  "The access key to use for the auditing Storage Account.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_account_access_key_is_secondary": {
			Description: "Is `storage_account_access_key` value the storage's secondary key?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"retention_in_days": {
			Description:  "The number of days to retain logs for in the Storage Account.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntBetween(0, 3285),
		},

		"log_monitoring_enabled": {
			Description: "Enable audit events to Azure Monitor?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},
	}
}
