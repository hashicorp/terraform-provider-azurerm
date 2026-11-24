// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_security_alert_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlManagedInstanceSecurityAlertPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"managed_instance_name": {
			Description:  "The name of the MS SQL Managed Instance.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlServerName,
		},

		"disabled_alerts": {
			Description: "An array of alerts that are disabled.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Set:         pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"Sql_Injection",
					"Sql_Injection_Vulnerability",
					"Access_Anomaly",
					"Data_Exfiltration",
					"Unsafe_Action",
					"Brute_Force",
				}, false),
			},
		},

		"email_account_admins_enabled": {
			Description: "Boolean flag which specifies if the alert is sent to the account administrators.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"email_addresses": {
			Description: "An array of email addresses to which the alert is sent.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Set: pluginsdk.HashString,
		},

		"retention_days": {
			Description:  "The number of days to keep in the Threat Detection audit logs.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"enabled": {
			Description: "The state of the Security Alert Policy, whether it is enabled or disabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"storage_account_access_key": {
			Description:  "The identifier key of the Threat Detection audit Storage Account. This is mandatory when you use `storage_endpoint` to specify a Storage Account blob endpoint.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_endpoint": {
			Description:  "The blob storage endpoint (e.g. https://example.blob.core.windows.net). This blob storage will hold all Threat Detection audit logs.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
