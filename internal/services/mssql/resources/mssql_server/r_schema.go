// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serverconnectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlServerSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Microsoft SQL Server. This needs to be globally unique within Azure.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlServerName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"version": {
			Description: "The version for the new server.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"2.0",
				"12.0",
			}, false),
		},

		"administrator_login": {
			Description:  "The administrator login name for the new server. Required unless `azuread_authentication_only` in the `azuread_administrator` block is `true`. When omitted, Azure will generate a default username which cannot be subsequently changed.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			AtLeastOneOf: []string{"administrator_login", "azuread_administrator.0.azuread_authentication_only"},
		},

		"administrator_login_password": {
			Description:   "The password associated with the `administrator_login` user. Needs to comply with Azure's [Password Policy](https://msdn.microsoft.com/library/ms161959.aspx).",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{"administrator_login_password_wo"},
			AtLeastOneOf:  []string{"administrator_login_password", "administrator_login_password_wo", "azuread_administrator.0.azuread_authentication_only"},
		},

		"administrator_login_password_wo": {
			Description:   "The Password associated with the `administrator_login` user. Needs to comply with Azure's [Password Policy](https://msdn.microsoft.com/library/ms161959.aspx).",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			WriteOnly:     true,
			RequiredWith:  []string{"administrator_login_password_wo_version"},
			ConflictsWith: []string{"administrator_login_password"},
			AtLeastOneOf:  []string{"administrator_login_password_wo", "administrator_login_password", "azuread_administrator.0.azuread_authentication_only"},
		},

		"administrator_login_password_wo_version": {
			Description:  "An integer value used to trigger an update for `administrator_login_password_wo`. This property should be incremented when updating `administrator_login_password_wo`.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
			RequiredWith: []string{"administrator_login_password_wo"},
		},

		"azuread_administrator": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"login_username": {
						Description:  "The login username of the Azure AD Administrator of this SQL Server.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"object_id": {
						Description:  "The object ID of the Azure AD Administrator of this SQL Server.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"tenant_id": {
						Description:  "The tenant ID of the Azure AD Administrator of this SQL Server.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IsUUID,
					},

					"azuread_authentication_only": {
						Description: "Whether only AD Users and administrators (e.g. `azuread_administrator[0].login_username`) can be used to login, or also local database users (e.g. `administrator_login`). When `true`, the `administrator_login` and `administrator_login_password` properties can be omitted.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},
				},
			},
		},

		"connection_policy": {
			Description: "The connection policy the server will use.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(serverconnectionpolicies.ServerConnectionTypeDefault),
			ValidateFunc: validation.StringInSlice(serverconnectionpolicies.PossibleValuesForServerConnectionType(),
				false),
		},

		"express_vulnerability_assessment_enabled": {
			Description: "Whether to enable the Express Vulnerability Assessment Configuration.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"transparent_data_encryption_key_vault_key_id": {
			Description:  "The fully versioned `Key Vault` `Key` URL (e.g. `'https://<YourVaultName>.vault.azure.net/keys/<YourKeyName>/<YourKeyVersion>`) to be used as the `Customer Managed Key` (CMK/BYOK) for the `Transparent Data Encryption` (TDE) layer.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},

		"primary_user_assigned_identity_id": {
			Description:  "The primary user Managed Identity ID.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			RequiredWith: []string{
				"identity",
			},
		},

		"minimum_tls_version": {
			Description: "The Minimum TLS Version for all SQL Database and SQL Data Warehouse databases associated with the server.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "1.2",
			ValidateFunc: validation.StringInSlice([]string{
				"1.2",
			}, false),
		},

		"public_network_access_enabled": {
			Description: "Whether public network access is allowed for this server.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"outbound_network_restriction_enabled": {
			Description: "Whether outbound network traffic is restricted for this server.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"soft_delete_retention_days": {
			Description: "The number of days a deleted server is retained in a soft-deleted state and can be restored. Setting this to `0` disables soft delete.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// Note: O+C because soft delete can be configured outside of Terraform (e.g. Portal/PowerShell) and Azure may apply a default retention period when not specified
			Computed:     true,
			ValidateFunc: validation.IntBetween(0, 7),
		},

		"fully_qualified_domain_name": {
			Description: "The fully qualified domain name of the Azure SQL Server (e.g. myServerName.database.windows.net).",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"restorable_dropped_database_ids": {
			Description: "A list of dropped restorable database IDs on the server.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}

	if !features.SixPointOh() {
		s["administrator_login_password_wo_version"] = &pluginsdk.Schema{
			Description:  "An integer value used to trigger an update for `administrator_login_password_wo`. This property should be incremented when updating `administrator_login_password_wo`.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"administrator_login_password_wo"},
		}
	}

	return s
}
