// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_microsoft_support_auditing_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlServerMicrosoftSupportAuditingPolicySchema() map[string]*pluginsdk.Schema {
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
			Description:  "The blob storage endpoint (e.g. https://example.blob.core.windows.net). This blob storage will hold all Microsoft support auditing logs.",
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

		"log_monitoring_enabled": {
			Description: "Enable audit events to Azure Monitor? To enable server audit events to Azure Monitor, please enable its main database audit events to Azure Monitor.",
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
	}
}
