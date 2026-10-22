// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/restorabledroppeddatabases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/sqlvulnerabilityassessmentssettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceMssqlServerSetFlatten(ctx context.Context, d *pluginsdk.ResourceData, id *commonids.SqlServerId, model *servers.Server, metaClient *clients.Client) error {
	connectionClient := metaClient.MSSQL.ServerConnectionPoliciesClient
	restorableDroppedDatabasesClient := metaClient.MSSQL.RestorableDroppedDatabasesClient
	vaClient := metaClient.MSSQL.SqlVulnerabilityAssessmentSettingsClient

	d.Set("name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		identity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("version", props.Version)
			d.Set("administrator_login", props.AdministratorLogin)
			d.Set("administrator_login_password_wo_version", d.Get("administrator_login_password_wo_version").(int))
			d.Set("fully_qualified_domain_name", props.FullyQualifiedDomainName)

			// todo remove `|| *v == "None"` when https://github.com/Azure/azure-rest-api-specs/issues/24348 is addressed
			if v := props.MinimalTlsVersion; v == nil || *v == "None" {
				d.Set("minimum_tls_version", "Disabled")
			} else {
				d.Set("minimum_tls_version", pointer.FromEnum(props.MinimalTlsVersion))
			}

			d.Set("public_network_access_enabled", pointer.From(props.PublicNetworkAccess) == servers.ServerPublicNetworkAccessFlagEnabled)
			d.Set("outbound_network_restriction_enabled", pointer.From(props.RestrictOutboundNetworkAccess) == servers.ServerNetworkAccessFlagEnabled)

			primaryUserAssignedIdentityID := ""
			if props.PrimaryUserAssignedIdentityId != nil && pointer.From(props.PrimaryUserAssignedIdentityId) != "" {
				parsedPrimaryUserAssignedIdentityID, err := commonids.ParseUserAssignedIdentityIDInsensitively(pointer.From(props.PrimaryUserAssignedIdentityId))
				if err != nil {
					return err
				}
				primaryUserAssignedIdentityID = parsedPrimaryUserAssignedIdentityID.ID()
			}

			d.Set("primary_user_assigned_identity_id", primaryUserAssignedIdentityID)
			d.Set("transparent_data_encryption_key_vault_key_id", props.KeyId)

			if props.Administrators != nil {
				d.Set("azuread_administrator", flattenMsSqlServerAdministrators(*props.Administrators))
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	connection, err := connectionClient.Get(ctx, pointer.From(id))
	if err != nil {
		return fmt.Errorf("retrieving SQL Server Blob Connection Policy %s: %v ", id, err)
	}

	if model := connection.Model; model != nil && model.Properties != nil {
		d.Set("connection_policy", string(model.Properties.ConnectionType))
	}

	restorableListPage, err := restorableDroppedDatabasesClient.ListByServerComplete(ctx, pointer.From(id))
	if err != nil {
		return fmt.Errorf("listing SQL Server Restorable Dropped Databases %s: %v", id, err)
	}

	if err := d.Set("restorable_dropped_database_ids", flattenSqlServerRestorableDatabases(restorableListPage)); err != nil {
		return fmt.Errorf("setting `restorable_dropped_database_ids`: %+v", err)
	}

	va, err := vaClient.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Express Vulnerability Assessment Settings for %s: %+v", *id, err)
	}

	if model := va.Model; model != nil && model.Properties != nil {
		d.Set("express_vulnerability_assessment_enabled", pointer.From(model.Properties.State) == sqlvulnerabilityassessmentssettings.SqlVulnerabilityAssessmentStateEnabled)
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenMsSqlServerAdministrators(admin servers.ServerExternalAdministrator) []any {
	var login, sid, tid string
	if admin.Login != nil {
		login = *admin.Login
	}

	if admin.Sid != nil {
		sid = pointer.From(admin.Sid)
	}

	if admin.TenantId != nil {
		tid = pointer.From(admin.TenantId)
	}

	var aadOnlyAuthenticationEnabled bool
	if admin.AzureADOnlyAuthentication != nil {
		aadOnlyAuthenticationEnabled = pointer.From(admin.AzureADOnlyAuthentication)
	}

	return []any{
		map[string]any{
			"login_username":              login,
			"object_id":                   sid,
			"tenant_id":                   tid,
			"azuread_authentication_only": aadOnlyAuthenticationEnabled,
		},
	}
}

func flattenSqlServerRestorableDatabases(resp restorabledroppeddatabases.ListByServerCompleteResult) []string {
	if len(resp.Items) == 0 {
		return []string{}
	}

	res := make([]string, 0)
	for _, r := range resp.Items {
		var id string
		if r.Id != nil {
			id = *r.Id
			res = append(res, id)
		}
	}

	return res
}
