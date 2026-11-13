// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serverazureadonlyauthentications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serverconnectionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/sqlvulnerabilityassessmentssettings"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlServerUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	connectionClient := meta.(*clients.Client).MSSQL.ServerConnectionPoliciesClient
	adminClient := meta.(*clients.Client).MSSQL.ServerAzureADAdministratorsClient
	aadOnlyAuthenticationsClient := meta.(*clients.Client).MSSQL.ServerAzureADOnlyAuthenticationsClient
	vaClient := meta.(*clients.Client).MSSQL.SqlVulnerabilityAssessmentSettingsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSqlServerID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("azuread_administrator") {
		log.Printf("[INFO] Expanding 'azuread_administrator' to see if we need Create or Delete")
		if adminProps := expandMsSqlServerAdministrator(d.Get("azuread_administrator").([]any)); adminProps != nil {
			if err := adminClient.CreateOrUpdateThenPoll(ctx, *id, pointer.From(adminProps)); err != nil {
				return fmt.Errorf("updating Azure Active Directory Administrator %s: %+v", id, err)
			}
		} else {
			if _, err := adminClient.Get(ctx, *id); err != nil {
				return fmt.Errorf("retrieving Azure Active Directory Administrator %s: %+v", id, err)
			}

			if err = adminClient.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting Azure Active Directory Administrator %s: %+v", id, err)
			}
		}
	}

	// The `AzureADOnlyAuthentication` cannot be updated by `serversClient`,
	// Service return `Invalid value given for parameter AzureADOnlyAuthentication. Specify a valid parameter value.` in that case.
	if d.HasChange("azuread_administrator") && d.HasChange("azuread_administrator.0.azuread_authentication_only") {
		if aadOnlyAuthenticationEnabled := expandMsSqlServerAADOnlyAuthentication(d.Get("azuread_administrator").([]any)); aadOnlyAuthenticationEnabled {
			aadOnlyAuthenticationProps := serverazureadonlyauthentications.ServerAzureADOnlyAuthentication{
				Properties: &serverazureadonlyauthentications.AzureADOnlyAuthProperties{
					AzureADOnlyAuthentication: true,
				},
			}

			if err := aadOnlyAuthenticationsClient.CreateOrUpdateThenPoll(ctx, *id, aadOnlyAuthenticationProps); err != nil {
				return fmt.Errorf("updating Azure Active Directory Only Authentication for %s: %+v", id, err)
			}
		} else {
			resp, err := aadOnlyAuthenticationsClient.Delete(ctx, *id)
			if err != nil {
				log.Printf("[INFO] Deletion of Azure Active Directory Only Authentication failed for %s: %+v", pointer.From(id), err)
				return fmt.Errorf("deleting Azure Active Directory Only Authentications for %s: %+v", pointer.From(id), err)
			}

			// NOTE: This call does not return a future it returns a response, but you will get a future back if the status code is 202...
			// https://learn.microsoft.com/rest/api/sql/server-azure-ad-only-authentications/delete?view=rest-sql-2023-05-01-preview&tabs=HTTP
			if response.WasStatusCode(resp.HttpResponse, 202) {
				// NOTE: It was accepted but not completed, it is now an async operation...
				// create a custom poller and wait for it to complete as 'Succeeded'...
				log.Printf("[INFO] Delete Azure Active Directory Only Administrators response was a 202 WaitForStateContext...")

				initialDelayDuration := 5 * time.Second
				pollerType := custompollers.NewMsSqlServerDeleteServerAzureADOnlyAuthenticationPoller(aadOnlyAuthenticationsClient, pointer.From(id))
				poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("waiting for the deletion of the Azure Active Directory Only Administrator: %+v", err)
				}
			}
		}
	}

	existing, err := client.Get(ctx, *id, servers.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if payload := existing.Model; payload != nil {
		if d.HasChange("tags") {
			payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
		}

		if d.HasChange("identity") {
			expanded, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			payload.Identity = expanded
		}

		if d.HasChange("transparent_data_encryption_key_vault_key_id") {
			keyId, err := keyvault.ParseNestedItemID(d.Get("transparent_data_encryption_key_vault_key_id").(string), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
			if err != nil {
				return err
			}
			payload.Properties.KeyId = pointer.To(keyId.ID())
		}

		// The `primary_user_assigned_identity_id` is an `O+C` property, when it's being removed from configuration.
		// `HasChange()` will be `true`, but `GetOk()` will be `false`. So use `d.Get()` to read the value.
		if d.HasChange("primary_user_assigned_identity_id") {
			payload.Properties.PrimaryUserAssignedIdentityId = pointer.To(d.Get("primary_user_assigned_identity_id").(string))
		}

		payload.Properties.PublicNetworkAccess = pointer.To(servers.ServerPublicNetworkAccessFlagDisabled)
		payload.Properties.RestrictOutboundNetworkAccess = pointer.To(servers.ServerNetworkAccessFlagDisabled)

		if v := d.Get("public_network_access_enabled"); v.(bool) {
			payload.Properties.PublicNetworkAccess = pointer.To(servers.ServerPublicNetworkAccessFlagEnabled)
		}

		if v := d.Get("outbound_network_restriction_enabled"); v.(bool) {
			payload.Properties.RestrictOutboundNetworkAccess = pointer.To(servers.ServerNetworkAccessFlagEnabled)
		}

		if d.HasChange("administrator_login_password") {
			payload.Properties.AdministratorLoginPassword = pointer.To(d.Get("administrator_login_password").(string))
		}

		if d.HasChange("administrator_login_password_wo_version") {
			woAdminLoginPassword, err := pluginsdk.GetWriteOnly(d, "administrator_login_password_wo", cty.String)
			if err != nil {
				return err
			}
			if !woAdminLoginPassword.IsNull() {
				payload.Properties.AdministratorLoginPassword = pointer.To(woAdminLoginPassword.AsString())
			}
		}

		if d.HasChange("minimum_tls_version") {
			payload.Properties.MinimalTlsVersion = pointer.ToEnum[servers.MinimalTlsVersion](d.Get("minimum_tls_version").(string))
		}

		if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	if d.HasChange("connection_policy") {
		connection := serverconnectionpolicies.ServerConnectionPolicy{
			Properties: &serverconnectionpolicies.ServerConnectionPolicyProperties{
				ConnectionType: serverconnectionpolicies.ServerConnectionType(d.Get("connection_policy").(string)),
			},
		}
		if err = connectionClient.CreateOrUpdateThenPoll(ctx, *id, connection); err != nil {
			return fmt.Errorf("updating Connection Policy for %s: %+v", id, err)
		}
	}

	if d.HasChange("express_vulnerability_assessment_enabled") {
		vaState := sqlvulnerabilityassessmentssettings.SqlVulnerabilityAssessmentStateDisabled
		if d.Get("express_vulnerability_assessment_enabled").(bool) {
			vaState = sqlvulnerabilityassessmentssettings.SqlVulnerabilityAssessmentStateEnabled
		}

		va := sqlvulnerabilityassessmentssettings.SqlVulnerabilityAssessment{
			Properties: &sqlvulnerabilityassessmentssettings.SqlVulnerabilityAssessmentPolicyProperties{
				State: &vaState,
			},
		}

		if _, err := vaClient.CreateOrUpdate(ctx, *id, va); err != nil {
			return fmt.Errorf("updating Express Vulnerability Assessment Settings for %s: %+v", *id, err)
		}
	}

	return resourceMsSqlServerRead(d, meta)
}
