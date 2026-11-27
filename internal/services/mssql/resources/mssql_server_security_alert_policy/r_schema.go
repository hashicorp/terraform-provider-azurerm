// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_security_alert_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlServerSecurityAlertPolicySchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"server_name": {
			Description:  "The name of the MS SQL Server.",
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
				}, false),
			},
		},

		"email_account_admins_enabled": {
			Description: "Whether the alerts are sent to the account administrators.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"email_addresses": {
			Description: "An array of email addresses to which the alert is sent.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Set:         pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"retention_days": {
			Description:  "The number of days to keep the Threat Detection audit logs.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"enabled": {
			Description: "Whether the policy is enabled.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},

		"storage_account_access_key": {
			Description:  "The primary access key of the Threat Detection audit logs blob storage endpoint.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"storage_endpoint"},
		},

		"storage_endpoint": {
			Description:  "The blob storage endpoint that will hold all Threat Detection audit logs (e.g., `https://example.blob.core.windows.net`).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"storage_account_access_key"},
		},
	}

	if !features.SixPointOh() {
		s["state"] = &pluginsdk.Schema{
			Description:  "The state of the policy.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(serversecurityalertpolicies.PossibleValuesForSecurityAlertsPolicyState(), false),
			ExactlyOneOf: []string{"enabled", "state"},
			Deprecated:   "`state` has been deprecated in favour of `enabled` and will be removed in v6.0 of the AzureRM Provider",
		}

		s["enabled"] = &pluginsdk.Schema{
			Description:  "Whether the policy is enabled.",
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			Computed:     true,
			ExactlyOneOf: []string{"enabled", "state"},
		}
	}

	return s
}
