// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultParser "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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
				RequiredWith: []string{"administrator_login", "administrator_login_password"},
			},

			"administrator_login_password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				AtLeastOneOf: []string{"administrator_login_password", "azuread_administrator.0.azuread_authentication_only"},
				RequiredWith: []string{"administrator_login", "administrator_login_password"},
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
				Default:  string(sql.ServerConnectionTypeDefault),
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.ServerConnectionTypeDefault),
					string(sql.ServerConnectionTypeProxy),
					string(sql.ServerConnectionTypeRedirect),
				}, false),
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
					"1.0",
					"1.1",
					"1.2",
					"Disabled",
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

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(msSqlMinimumTLSVersionDiff),

			pluginsdk.CustomizeDiffShim(msSqlPasswordChangeWhenAADAuthOnly),
		),
	}
}

func resourceMsSqlServerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	connectionClient := meta.(*clients.Client).MSSQL.ServerConnectionPoliciesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	location := azure.NormalizeLocation(d.Get("location").(string))
	version := d.Get("version").(string)

	t := d.Get("tags").(map[string]interface{})
	metadata := tags.Expand(t)

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_mssql_server", id.ID())
	}

	props := sql.Server{
		Location: utils.String(location),
		Tags:     metadata,
		ServerProperties: &sql.ServerProperties{
			Version:                       utils.String(version),
			PublicNetworkAccess:           sql.ServerNetworkAccessFlagEnabled,
			RestrictOutboundNetworkAccess: sql.ServerNetworkAccessFlagDisabled,
		},
	}

	if v := d.Get("administrator_login"); v.(string) != "" {
		props.ServerProperties.AdministratorLogin = utils.String(v.(string))
	}

	if v := d.Get("administrator_login_password"); v.(string) != "" {
		props.ServerProperties.AdministratorLoginPassword = utils.String(v.(string))
	}

	if azureADAdministrator, ok := d.GetOk("azuread_administrator"); ok {
		props.ServerProperties.Administrators = expandMsSqlServerAdministrators(azureADAdministrator.([]interface{}))
	}

	if v, ok := d.GetOk("identity"); ok {
		expandedIdentity, err := expandSqlServerIdentity(v.([]interface{}))
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
			// msSql requires the versioned key URL...
			props.KeyID = pointer.To(keyId.ID())
		} else {
			return fmt.Errorf("key vault key id must be a reference to a key, got %s", keyId.NestedItemType)
		}
	}

	if primaryUserAssignedIdentityID, ok := d.GetOk("primary_user_assigned_identity_id"); ok {
		props.PrimaryUserAssignedIdentityID = utils.String(primaryUserAssignedIdentityID.(string))
	}

	// if you pass the Key ID you must also define the PrimaryUserAssignedIdentityID...
	if props.KeyID != nil && props.PrimaryUserAssignedIdentityID == nil {
		return fmt.Errorf("the `primary_user_assigned_identity_id` field must be specified to use the 'transparent_data_encryption_key_vault_key_id' in %s", id)
	}

	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		props.ServerProperties.PublicNetworkAccess = sql.ServerNetworkAccessFlagDisabled
	}

	if v := d.Get("outbound_network_restriction_enabled"); v.(bool) {
		props.ServerProperties.RestrictOutboundNetworkAccess = sql.ServerNetworkAccessFlagEnabled
	}

	if v := d.Get("minimum_tls_version"); v.(string) != "Disabled" {
		props.ServerProperties.MinimalTLSVersion = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, props)
	if err != nil {
		return fmt.Errorf("issuing create request for %s: %+v", id.String(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasConflict(future.Response()) {
			return fmt.Errorf("SQL Server names need to be globally unique and %q is already in use", id.Name)
		}

		return fmt.Errorf("waiting for creation/update of %s: %+v", id.String(), err)
	}

	d.SetId(id.ID())

	connection := sql.ServerConnectionPolicy{
		ServerConnectionPolicyProperties: &sql.ServerConnectionPolicyProperties{
			ConnectionType: sql.ServerConnectionType(d.Get("connection_policy").(string)),
		},
	}
	if _, err = connectionClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, connection); err != nil {
		return fmt.Errorf("issuing create request for Connection Policy %s: %+v", id.String(), err)
	}

	return resourceMsSqlServerRead(d, meta)
}

func resourceMsSqlServerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	connectionClient := meta.(*clients.Client).MSSQL.ServerConnectionPoliciesClient
	adminClient := meta.(*clients.Client).MSSQL.ServerAzureADAdministratorsClient
	aadOnlyAuthenticationsClient := meta.(*clients.Client).MSSQL.ServerAzureADOnlyAuthenticationsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	version := d.Get("version").(string)

	t := d.Get("tags").(map[string]interface{})
	metadata := tags.Expand(t)

	props := sql.Server{
		Location: utils.String(location),
		Tags:     metadata,
		ServerProperties: &sql.ServerProperties{
			Version:                       utils.String(version),
			PublicNetworkAccess:           sql.ServerNetworkAccessFlagEnabled,
			RestrictOutboundNetworkAccess: sql.ServerNetworkAccessFlagDisabled,
		},
	}

	if v, ok := d.GetOk("identity"); ok {
		expandedIdentity, err := expandSqlServerIdentity(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		props.Identity = expandedIdentity
	} else {
		props.Identity = existing.Identity
	}

	if d.HasChange("key_vault_key_id") {
		keyVaultKeyId := d.Get(("transparent_data_encryption_key_vault_key_id")).(string)

		keyId, err := keyVaultParser.ParseNestedItemID(keyVaultKeyId)
		if err != nil {
			return fmt.Errorf("unable to parse key: %q: %+v", keyVaultKeyId, err)
		}

		if keyId.NestedItemType == keyVaultParser.NestedItemTypeKey {
			props.KeyID = pointer.To(keyId.ID())
		} else {
			return fmt.Errorf("key vault key id must be a reference to a key, got %s", keyId.NestedItemType)
		}
	}

	if primaryUserAssignedIdentityID, ok := d.GetOk("primary_user_assigned_identity_id"); ok {
		props.PrimaryUserAssignedIdentityID = utils.String(primaryUserAssignedIdentityID.(string))
	}

	// if you pass the Key ID you must also define the PrimaryUserAssignedIdentityID...
	if props.KeyID != nil && props.PrimaryUserAssignedIdentityID == nil {
		return fmt.Errorf("the `primary_user_assigned_identity_id` field must be specified to use the 'transparent_data_encryption_key_vault_key_id' in %s", id)
	}

	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		props.ServerProperties.PublicNetworkAccess = sql.ServerNetworkAccessFlagDisabled
	}

	if v := d.Get("outbound_network_restriction_enabled"); v.(bool) {
		props.ServerProperties.RestrictOutboundNetworkAccess = sql.ServerNetworkAccessFlagEnabled
	}

	if d.HasChange("administrator_login_password") {
		adminPassword := d.Get("administrator_login_password").(string)
		props.ServerProperties.AdministratorLoginPassword = utils.String(adminPassword)
	}

	if v := d.Get("minimum_tls_version"); v.(string) != "Disabled" {
		props.ServerProperties.MinimalTLSVersion = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, props)
	if err != nil {
		return fmt.Errorf("issuing update request for %s: %+v", id.String(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasConflict(future.Response()) {
			return fmt.Errorf("SQL Server names need to be globally unique and %q is already in use", id.Name)
		}

		return fmt.Errorf("waiting for update of %s: %+v", id.String(), err)
	}

	d.SetId(id.ID())

	if d.HasChange("azuread_administrator") {
		aadOnlyDeleteFuture, err := aadOnlyAuthenticationsClient.Delete(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if aadOnlyDeleteFuture.Response() == nil || aadOnlyDeleteFuture.Response().StatusCode != http.StatusBadRequest {
				return fmt.Errorf("deleting AD Only Authentications %s: %+v", id.String(), err)
			}
			log.Printf("[INFO] AD Only Authentication is not removed as AD Admin is not set for %s: %+v", id.String(), err)
		} else if err = aadOnlyDeleteFuture.WaitForCompletionRef(ctx, adminClient.Client); err != nil {
			return fmt.Errorf("waiting for deletion of AD Only Authentications %s: %+v", id.String(), err)
		}

		if adminParams := expandMsSqlServerAdministrator(d.Get("azuread_administrator").([]interface{})); adminParams != nil {
			adminFuture, err := adminClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, *adminParams)
			if err != nil {
				return fmt.Errorf("creating AAD admin %s: %+v", id.String(), err)
			}

			if err = adminFuture.WaitForCompletionRef(ctx, adminClient.Client); err != nil {
				return fmt.Errorf("waiting for creation of AAD admin %s: %+v", id.String(), err)
			}
		} else {
			adminDelFuture, err := adminClient.Delete(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("deleting AAD admin  %s: %+v", id.String(), err)
			}

			if err = adminDelFuture.WaitForCompletionRef(ctx, adminClient.Client); err != nil {
				return fmt.Errorf("waiting for deletion of AAD admin %s: %+v", id.String(), err)
			}
		}
	}

	if aadOnlyAuthentictionsEnabled := expandMsSqlServerAADOnlyAuthentictions(d.Get("azuread_administrator").([]interface{})); d.HasChange("azuread_administrator") && aadOnlyAuthentictionsEnabled {
		aadOnlyAuthentictionsParams := sql.ServerAzureADOnlyAuthentication{
			AzureADOnlyAuthProperties: &sql.AzureADOnlyAuthProperties{
				AzureADOnlyAuthentication: utils.Bool(aadOnlyAuthentictionsEnabled),
			},
		}
		aadOnlyEnabledFuture, err := aadOnlyAuthenticationsClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, aadOnlyAuthentictionsParams)
		if err != nil {
			return fmt.Errorf("setting AAD only authentication for %s: %+v", id.String(), err)
		}

		if err = aadOnlyEnabledFuture.WaitForCompletionRef(ctx, adminClient.Client); err != nil {
			return fmt.Errorf("waiting for setting of AAD only authentication for %s: %+v", id.String(), err)
		}
	}

	connection := sql.ServerConnectionPolicy{
		ServerConnectionPolicyProperties: &sql.ServerConnectionPolicyProperties{
			ConnectionType: sql.ServerConnectionType(d.Get("connection_policy").(string)),
		},
	}
	if _, err = connectionClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, connection); err != nil {
		return fmt.Errorf("issuing update request for Connection Policy %s: %+v", id.String(), err)
	}

	return resourceMsSqlServerRead(d, meta)
}

func resourceMsSqlServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	connectionClient := meta.(*clients.Client).MSSQL.ServerConnectionPoliciesClient
	restorableDroppedDatabasesClient := meta.(*clients.Client).MSSQL.RestorableDroppedDatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL Server %s - removing from state", id.String())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading SQL Server %s: %v", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	identity, err := flattenSqlServerIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := resp.ServerProperties; props != nil {
		d.Set("version", props.Version)
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("fully_qualified_domain_name", props.FullyQualifiedDomainName)
		// todo remove `|| *v == "None"` when https://github.com/Azure/azure-rest-api-specs/issues/24348 is addressed
		if v := props.MinimalTLSVersion; v == nil || *v == "None" {
			d.Set("minimum_tls_version", "Disabled")
		} else {
			d.Set("minimum_tls_version", props.MinimalTLSVersion)
		}
		d.Set("public_network_access_enabled", props.PublicNetworkAccess == sql.ServerNetworkAccessFlagEnabled)
		d.Set("outbound_network_restriction_enabled", props.RestrictOutboundNetworkAccess == sql.ServerNetworkAccessFlagEnabled)
		primaryUserAssignedIdentityID := ""
		if props.PrimaryUserAssignedIdentityID != nil && *props.PrimaryUserAssignedIdentityID != "" {
			parsedPrimaryUserAssignedIdentityID, err := commonids.ParseUserAssignedIdentityIDInsensitively(*props.PrimaryUserAssignedIdentityID)
			if err != nil {
				return err
			}
			primaryUserAssignedIdentityID = parsedPrimaryUserAssignedIdentityID.ID()
		}
		d.Set("primary_user_assigned_identity_id", primaryUserAssignedIdentityID)
		d.Set("transparent_data_encryption_key_vault_key_id", props.KeyID)
		if props.Administrators != nil {
			d.Set("azuread_administrator", flatternMsSqlServerAdministrators(*props.Administrators))
		}
	}

	connection, err := connectionClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("reading SQL Server %s Blob Connection Policy: %v ", id.Name, err)
	}

	if props := connection.ServerConnectionPolicyProperties; props != nil {
		d.Set("connection_policy", string(props.ConnectionType))
	}

	restorableListPage, err := restorableDroppedDatabasesClient.ListByServer(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("listing SQL Server %s Restorable Dropped Databases: %v", id.Name, err)
	}
	if err := d.Set("restorable_dropped_database_ids", flattenSqlServerRestorableDatabases(restorableListPage.Response())); err != nil {
		return fmt.Errorf("setting `restorable_dropped_database_ids`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMsSqlServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting SQL Server %s: %+v", id.Name, err)
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}

func expandSqlServerIdentity(input []interface{}) (*sql.ResourceIdentity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := sql.ResourceIdentity{
		Type: sql.IdentityType(string(expanded.Type)),
	}
	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.UserAssignedIdentities = make(map[string]*sql.UserIdentity)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &sql.UserIdentity{
				// intentionally empty
			}
		}
	}

	return &out, nil
}

func flattenSqlServerIdentity(input *sql.ResourceIdentity) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = input.PrincipalID.String()
		}
		if input.TenantID != nil {
			transform.TenantId = input.TenantID.String()
		}
		for k, v := range input.UserAssignedIdentities {
			details := identity.UserAssignedIdentityDetails{}
			if v.ClientID != nil {
				details.ClientId = utils.String(v.ClientID.String())
			}
			if v.PrincipalID != nil {
				details.PrincipalId = utils.String(v.PrincipalID.String())
			}
			transform.IdentityIds[k] = details
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func expandMsSqlServerAADOnlyAuthentictions(input []interface{}) bool {
	if len(input) == 0 || input[0] == nil {
		return false
	}
	admin := input[0].(map[string]interface{})
	if v, ok := admin["azuread_authentication_only"]; ok && v != nil {
		return v.(bool)
	}
	return false
}

func expandMsSqlServerAdministrator(input []interface{}) *sql.ServerAzureADAdministrator {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	admin := input[0].(map[string]interface{})
	sid, _ := uuid.FromString(admin["object_id"].(string))

	adminParams := sql.ServerAzureADAdministrator{
		AdministratorProperties: &sql.AdministratorProperties{
			AdministratorType: utils.String("ActiveDirectory"),
			Login:             utils.String(admin["login_username"].(string)),
			Sid:               &sid,
		},
	}

	if v, ok := admin["tenant_id"]; ok && v != "" {
		tid, _ := uuid.FromString(v.(string))
		adminParams.TenantID = &tid
	}

	return &adminParams
}

func expandMsSqlServerAdministrators(input []interface{}) *sql.ServerExternalAdministrator {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	admin := input[0].(map[string]interface{})
	sid, _ := uuid.FromString(admin["object_id"].(string))

	adminParams := sql.ServerExternalAdministrator{
		AdministratorType: sql.AdministratorTypeActiveDirectory,
		Login:             utils.String(admin["login_username"].(string)),
		Sid:               &sid,
	}

	if v, ok := admin["tenant_id"]; ok && v != "" {
		tid, _ := uuid.FromString(v.(string))
		adminParams.TenantID = &tid
	}

	if v, ok := admin["azuread_authentication_only"]; ok && v != "" {
		adOnlyAuthentication := v.(bool)
		adminParams.AzureADOnlyAuthentication = &adOnlyAuthentication
	}

	return &adminParams
}

func flatternMsSqlServerAdministrators(admin sql.ServerExternalAdministrator) []interface{} {
	var login, sid, tid string
	if admin.Login != nil {
		login = *admin.Login
	}

	if admin.Sid != nil {
		sid = admin.Sid.String()
	}

	if admin.TenantID != nil {
		tid = admin.TenantID.String()
	}

	var aadOnlyAuthentictionsEnabled bool
	if admin.AzureADOnlyAuthentication != nil {
		aadOnlyAuthentictionsEnabled = *admin.AzureADOnlyAuthentication
	}

	return []interface{}{
		map[string]interface{}{
			"login_username":              login,
			"object_id":                   sid,
			"tenant_id":                   tid,
			"azuread_authentication_only": aadOnlyAuthentictionsEnabled,
		},
	}
}

func flattenSqlServerRestorableDatabases(resp sql.RestorableDroppedDatabaseListResult) []string {
	if resp.Value == nil || len(*resp.Value) == 0 {
		return []string{}
	}
	res := make([]string, 0)
	for _, r := range *resp.Value {
		var id string
		if r.ID != nil {
			id = *r.ID
		}
		res = append(res, id)
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
