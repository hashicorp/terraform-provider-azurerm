// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serverconnectionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/sqlvulnerabilityassessmentssettings"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlServerCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	connectionClient := meta.(*clients.Client).MSSQL.ServerConnectionPoliciesClient
	vaClient := meta.(*clients.Client).MSSQL.SqlVulnerabilityAssessmentSettingsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSqlServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, servers.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mssql_server", id.ID())
		}
	}

	props := servers.Server{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
		Properties: &servers.ServerProperties{
			Version:                       pointer.To(d.Get("version").(string)),
			PublicNetworkAccess:           pointer.To(servers.ServerPublicNetworkAccessFlagEnabled),
			RestrictOutboundNetworkAccess: pointer.To(servers.ServerNetworkAccessFlagDisabled),
		},
	}

	woAdminLoginPassword, err := pluginsdk.GetWriteOnly(d, "administrator_login_password_wo", cty.String)
	if err != nil {
		return err
	}

	adminLogin := d.Get("administrator_login").(string)
	adminLoginPassword := d.Get("administrator_login_password").(string)

	if adminLogin != "" {
		if adminLoginPassword == "" && woAdminLoginPassword.IsNull() {
			return fmt.Errorf("expected `administrator_login_password` or `administrator_login_password_wo` to be set when `administrator_login` is specified")
		}

		props.Properties.AdministratorLogin = pointer.To(adminLogin)
	}

	if adminLoginPassword != "" {
		props.Properties.AdministratorLoginPassword = pointer.To(adminLoginPassword)
	}

	if !woAdminLoginPassword.IsNull() {
		props.Properties.AdministratorLoginPassword = pointer.To(woAdminLoginPassword.AsString())
	}

	// NOTE: You must set the admin before setting the values of the admin...
	if azureADAdministrator, ok := d.GetOk("azuread_administrator"); ok {
		props.Properties.Administrators = expandMsSqlServerAdministrators(azureADAdministrator.([]any))
	}

	if v, ok := d.GetOk("identity"); ok {
		expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(v.([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		props.Identity = expandedIdentity
	}

	if v, ok := d.GetOk("transparent_data_encryption_key_vault_key_id"); ok {
		keyId, err := keyvault.ParseNestedItemID(v.(string), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
		if err != nil {
			return err
		}
		props.Properties.KeyId = pointer.To(keyId.ID())
	}

	if primaryUserAssignedIdentityID, ok := d.GetOk("primary_user_assigned_identity_id"); ok {
		props.Properties.PrimaryUserAssignedIdentityId = pointer.To(primaryUserAssignedIdentityID.(string))
	}

	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		props.Properties.PublicNetworkAccess = pointer.To(servers.ServerPublicNetworkAccessFlagDisabled)
	}

	if v := d.Get("outbound_network_restriction_enabled"); v.(bool) {
		props.Properties.RestrictOutboundNetworkAccess = pointer.To(servers.ServerNetworkAccessFlagEnabled)
	}

	if v := d.Get("minimum_tls_version"); v.(string) != "Disabled" {
		props.Properties.MinimalTlsVersion = pointer.ToEnum[servers.MinimalTlsVersion](v.(string))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, props, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	connection := serverconnectionpolicies.ServerConnectionPolicy{
		Properties: &serverconnectionpolicies.ServerConnectionPolicyProperties{
			ConnectionType: serverconnectionpolicies.ServerConnectionType(d.Get("connection_policy").(string)),
		},
	}

	if err = connectionClient.CreateOrUpdateThenPoll(ctx, id, connection); err != nil {
		return fmt.Errorf("creating Connection Policy for %s: %+v", id, err)
	}

	vaState := sqlvulnerabilityassessmentssettings.SqlVulnerabilityAssessmentStateDisabled
	if d.Get("express_vulnerability_assessment_enabled").(bool) {
		vaState = sqlvulnerabilityassessmentssettings.SqlVulnerabilityAssessmentStateEnabled
	}

	va := sqlvulnerabilityassessmentssettings.SqlVulnerabilityAssessment{
		Properties: &sqlvulnerabilityassessmentssettings.SqlVulnerabilityAssessmentPolicyProperties{
			State: &vaState,
		},
	}

	if _, err := vaClient.CreateOrUpdate(ctx, id, va); err != nil {
		return fmt.Errorf("creating Express Vulnerability Assessment Settings for %s: %+v", id, err)
	}

	return resourceMsSqlServerRead(d, meta)
}
