// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_extended_auditing_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlServerExtendedAuditingPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"server_id": {
			Description:  "The ID of the SQL Server to set the extended auditing policy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlServerIDInsensitively),
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
			Description: "Whether the `storage_account_access_key` value is the storage's secondary key.",
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
			Description: "Whether to enable audit events to Azure Monitor. To enable server audit events to Azure Monitor, enable its main database audit events to Azure Monitor.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"storage_account_subscription_id": {
			Description:  "The ID of the Subscription containing the Storage Account.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.IsUUID,
		},

		"predicate_expression": {
			Description:  "The condition of the where clause when creating an audit.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"audit_actions_and_groups": {
			Description: "A list of Actions-Groups and Actions to audit.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			// Note: O+C because `audit_actions_and_groups` seems to be pre-populated with values ["SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP", "FAILED_DATABASE_AUTHENTICATION_GROUP", "BATCH_COMPLETED_GROUP"],
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}
