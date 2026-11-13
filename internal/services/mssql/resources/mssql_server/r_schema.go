// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serverconnectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlServerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlServerName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"version": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"2.0",
				"12.0",
			}, false),
		},

		"administrator_login": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			AtLeastOneOf: []string{"administrator_login", "azuread_administrator.0.azuread_authentication_only"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"administrator_login_password": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			AtLeastOneOf:  []string{"administrator_login_password", "administrator_login_password_wo", "azuread_administrator.0.azuread_authentication_only"},
			ConflictsWith: []string{"administrator_login_password_wo"},
		},

		"administrator_login_password_wo": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			WriteOnly:     true,
			AtLeastOneOf:  []string{"administrator_login_password_wo", "administrator_login_password", "azuread_administrator.0.azuread_authentication_only"},
			ConflictsWith: []string{"administrator_login_password"},
			RequiredWith:  []string{"administrator_login_password_wo_version"},
		},

		"administrator_login_password_wo_version": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"administrator_login_password_wo"},
		},

		"azuread_administrator": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"login_username": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"object_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IsUUID,
					},

					"azuread_authentication_only": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
					},
				},
			},
		},

		"connection_policy": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(serverconnectionpolicies.ServerConnectionTypeDefault),
			ValidateFunc: validation.StringInSlice(serverconnectionpolicies.PossibleValuesForServerConnectionType(),
				false),
		},

		"express_vulnerability_assessment_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"transparent_data_encryption_key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},

		"primary_user_assigned_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			RequiredWith: []string{
				"identity",
			},
		},

		"minimum_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "1.2",
			ValidateFunc: validation.StringInSlice([]string{
				"1.2",
			}, false),
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"outbound_network_restriction_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"fully_qualified_domain_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"restorable_dropped_database_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}
