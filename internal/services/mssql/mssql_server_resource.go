// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/restorabledroppeddatabases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverazureadadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverazureadonlyauthentications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverconnectionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/sqlvulnerabilityassessmentssettings"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	keyVaultParser "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlServer() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceMsSqlServerCreate,
		Read:   resourceMsSqlServerRead,
		Update: resourceMsSqlServerUpdate,
		Delete: resourceMsSqlServerDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ServerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
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
				Computed:     true,
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
							Computed:     true,
							ValidateFunc: validation.IsUUID,
						},

						"azuread_authentication_only": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Computed: true,
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
				ValidateFunc: keyVaultValidate.NestedItemId,
			},

			"primary_user_assigned_identity_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
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
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(msSqlMinimumTLSVersionDiff),

			pluginsdk.CustomizeDiffShim(msSqlPasswordChangeWhenAADAuthOnly),
		),
	}

	if !features.FivePointOh() {
		resource.Schema["minimum_tls_version"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "1.2",
			ValidateFunc: validation.StringInSlice([]string{
				"1.0",
				"1.1",
				"1.2",
				"Disabled",
			}, false),
		}
	}

	return resource
}

func resourceMsSqlServerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	connectionClient := meta.(*clients.Client).MSSQL.ServerConnectionPoliciesClient
	vaClient := meta.(*clients.Client).MSSQL.SqlVulnerabilityAssessmentSettingsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSqlServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, servers.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_mssql_server", id.ID())
	}

	props := servers.Server{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
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
		props.Properties.Administrators = expandMsSqlServerAdministrators(azureADAdministrator.([]interface{}))
	}

	if v, ok := d.GetOk("identity"); ok {
		expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		props.Identity = expandedIdentity
	}

	if v, ok := d.GetOk("transparent_data_encryption_key_vault_key_id"); ok {
		keyVaultKeyId := v.(string)

		keyId, err := keyVaultParser.ParseNestedItemID(keyVaultKeyId)
		if err != nil {
			return fmt.Errorf("unable to parse key: %q: %+v", keyVaultKeyId, err)
		}

		if keyId.NestedItemType == keyVaultParser.NestedItemTypeKey {
			// NOTE: msSql requires the versioned key URL...
			props.Properties.KeyId = pointer.To(keyId.ID())
		} else {
			return fmt.Errorf("key vault key id must be a reference to a key, got %s", keyId.NestedItemType)
		}
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
		props.Properties.MinimalTlsVersion = pointer.To(servers.MinimalTlsVersion(v.(string)))
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

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

func resourceMsSqlServerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	existing, err := client.Get(ctx, *id, servers.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if payload := existing.Model; payload != nil {
		if d.HasChange("tags") {
			payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
		}

		if d.HasChange("identity") {
			expanded, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			payload.Identity = expanded
		}

		if d.HasChange("transparent_data_encryption_key_vault_key_id") {
			keyVaultKeyId := d.Get(("transparent_data_encryption_key_vault_key_id")).(string)

			keyId, err := keyVaultParser.ParseNestedItemID(keyVaultKeyId)
			if err != nil {
				return fmt.Errorf("unable to parse key: %q: %+v", keyVaultKeyId, err)
			}

			if keyId.NestedItemType == keyVaultParser.NestedItemTypeKey {
				payload.Properties.KeyId = pointer.To(keyId.ID())
			} else {
				return fmt.Errorf("key vault key id must be a reference to a key, got %s", keyId.NestedItemType)
			}
		}

		if primaryUserAssignedIdentityID, ok := d.GetOk("primary_user_assigned_identity_id"); ok {
			payload.Properties.PrimaryUserAssignedIdentityId = pointer.To(primaryUserAssignedIdentityID.(string))
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
			payload.Properties.MinimalTlsVersion = pointer.To(servers.MinimalTlsVersion(d.Get("minimum_tls_version").(string)))
		}

		err := client.CreateOrUpdateThenPoll(ctx, *id, *payload)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	if d.HasChange("azuread_administrator") {
		// need to check if aadOnly is enabled or not before calling delete, else you will get the following error:
		// InvalidServerAADOnlyAuthNoAADAdminPropertyName: AAD Admin is not configured, AAD Admin must be set
		// before enabling/disabling AAD Only Authentication.
		log.Printf("[INFO] Checking if Azure Active Directory Administrators exist")
		aadOnlyAdmin := false

		resp, err := adminClient.Get(ctx, pointer.From(id))
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("retrieving Azure Active Directory Administrators %s: %+v", pointer.From(id), err)
			}
		} else {
			aadOnlyAdmin = true
		}

		if aadOnlyAdmin {
			resp, err := aadOnlyAuthenticationsClient.Delete(ctx, *id)
			if err != nil {
				log.Printf("[INFO] Deletion of Azure Active Directory Only Authentication failed for %s: %+v", pointer.From(id), err)
				return fmt.Errorf("deleting Azure Active Directory Only Authentications for %s: %+v", pointer.From(id), err)
			}

			// NOTE: This call does not return a future it returns a response, but you will get a future back if the status code is 202...
			// https://learn.microsoft.com/en-us/rest/api/sql/server-azure-ad-only-authentications/delete?view=rest-sql-2023-05-01-preview&tabs=HTTP
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

		log.Printf("[INFO] Expanding 'azuread_administrator' to see if we need Create or Delete")
		if adminProps := expandMsSqlServerAdministrator(d.Get("azuread_administrator").([]interface{})); adminProps != nil {
			err := adminClient.CreateOrUpdateThenPoll(ctx, *id, pointer.From(adminProps))
			if err != nil {
				return fmt.Errorf("creating Azure Active Directory Administrator %s: %+v", id, err)
			}
		} else {
			_, err := adminClient.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving Azure Active Directory Administrator %s: %+v", id, err)
			}

			err = adminClient.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting Azure Active Directory Administrator %s: %+v", id, err)
			}
		}
	}

	if aadOnlyAuthenticationEnabled := expandMsSqlServerAADOnlyAuthentication(d.Get("azuread_administrator").([]interface{})); d.HasChange("azuread_administrator") && aadOnlyAuthenticationEnabled {
		aadOnlyAuthenticationProps := serverazureadonlyauthentications.ServerAzureADOnlyAuthentication{
			Properties: &serverazureadonlyauthentications.AzureADOnlyAuthProperties{
				AzureADOnlyAuthentication: aadOnlyAuthenticationEnabled,
			},
		}

		err := aadOnlyAuthenticationsClient.CreateOrUpdateThenPoll(ctx, *id, aadOnlyAuthenticationProps)
		if err != nil {
			return fmt.Errorf("updating Azure Active Directory Only Authentication for %s: %+v", id, err)
		}
	}

	connection := serverconnectionpolicies.ServerConnectionPolicy{
		Properties: &serverconnectionpolicies.ServerConnectionPolicyProperties{
			ConnectionType: serverconnectionpolicies.ServerConnectionType(d.Get("connection_policy").(string)),
		},
	}

	if err = connectionClient.CreateOrUpdateThenPoll(ctx, *id, connection); err != nil {
		return fmt.Errorf("updating Connection Policy for %s: %+v", id, err)
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

func resourceMsSqlServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	connectionClient := meta.(*clients.Client).MSSQL.ServerConnectionPoliciesClient
	restorableDroppedDatabasesClient := meta.(*clients.Client).MSSQL.RestorableDroppedDatabasesClient
	vaClient := meta.(*clients.Client).MSSQL.SqlVulnerabilityAssessmentSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSqlServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, servers.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Error retrieving SQL Server %s - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving SQL Server %s: %v", id, err)
	}

	d.Set("name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
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
				d.Set("minimum_tls_version", string(pointer.From(props.MinimalTlsVersion)))
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

	return nil
}

func resourceMsSqlServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSqlServerID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, pointer.From(id))
	if err != nil {
		return fmt.Errorf("deleting SQL Server %s: %+v", id, err)
	}

	return nil
}

func expandMsSqlServerAADOnlyAuthentication(input []interface{}) bool {
	if len(input) == 0 || input[0] == nil {
		return false
	}

	admin := input[0].(map[string]interface{})

	if v, ok := admin["azuread_authentication_only"]; ok && v != nil {
		return v.(bool)
	}

	return false
}

func expandMsSqlServerAdministrator(input []interface{}) *serverazureadadministrators.ServerAzureADAdministrator {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	adminProps := serverazureadadministrators.ServerAzureADAdministrator{
		Properties: &serverazureadadministrators.AdministratorProperties{
			AdministratorType: serverazureadadministrators.AdministratorType(servers.AdministratorTypeActiveDirectory),
			Login:             v["login_username"].(string),
			Sid:               v["object_id"].(string),
		},
	}

	if t, ok := v["tenant_id"]; ok && t != "" {
		adminProps.Properties.TenantId = pointer.To(t.(string))
	}

	return pointer.To(adminProps)
}

func expandMsSqlServerAdministrators(input []interface{}) *servers.ServerExternalAdministrator {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	admin := input[0].(map[string]interface{})
	sid := admin["object_id"].(string)

	adminParams := servers.ServerExternalAdministrator{
		AdministratorType: pointer.To(servers.AdministratorTypeActiveDirectory),
		Login:             pointer.To(admin["login_username"].(string)),
		Sid:               pointer.To(sid),
	}

	if v, ok := admin["tenant_id"]; ok && v != "" {
		adminParams.TenantId = pointer.To(v.(string))
	}

	if v, ok := admin["azuread_authentication_only"]; ok && v != "" {
		adOnlyAuthentication := v.(bool)
		adminParams.AzureADOnlyAuthentication = &adOnlyAuthentication
	}

	return &adminParams
}

func flattenMsSqlServerAdministrators(admin servers.ServerExternalAdministrator) []interface{} {
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

	return []interface{}{
		map[string]interface{}{
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

func msSqlMinimumTLSVersionDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) (err error) {
	old, new := d.GetChange("minimum_tls_version")
	// todo remove `old != "None"` when https://github.com/Azure/azure-rest-api-specs/issues/24348 is addressed
	if old != "" && old != "None" && old != "Disabled" && new == "Disabled" {
		err = fmt.Errorf("`minimum_tls_version` cannot be removed once set, please set a valid value for this property")
	}
	return
}

func msSqlPasswordChangeWhenAADAuthOnly(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) (err error) {
	old, _ := d.GetChange("azuread_administrator.0.azuread_authentication_only")
	if old.(bool) && d.HasChange("administrator_login_password") {
		err = fmt.Errorf("`administrator_login_password` cannot be changed once `azuread_administrator.0.azuread_authentication_only = true`")
	}
	return
}
