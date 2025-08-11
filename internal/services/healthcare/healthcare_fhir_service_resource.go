// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHealthcareApisFhirService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareApisFhirServiceCreate,
		Read:   resourceHealthcareApisFhirServiceRead,
		Update: resourceHealthcareApisFhirServiceUpdate,
		Delete: resourceHealthcareApisFhirServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.HealthCareFhirV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := fhirservices.ParseFhirServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FhirServiceName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"location": commonschema.Location(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(fhirservices.FhirServiceKindFhirNegativeRFour),
				ValidateFunc: validation.StringInSlice([]string{
					string(fhirservices.FhirServiceKindFhirNegativeRFour),
					string(fhirservices.FhirServiceKindFhirNegativeStuThree),
				}, false),
			},

			"access_policy_object_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},

			"authentication": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authority": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"audience": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"smart_proxy_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

			// can't use the registry ID due to the ID cannot be obtained when setting the property in state file
			"container_registry_login_server_url": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"oci_artifact": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"login_server": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"image_name": {
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
							Optional:     true,
						},

						"digest": {
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
							Optional:     true,
						},
					},
				},
			},

			"cors": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"allowed_origins": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MaxItems: 64,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"allowed_headers": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MaxItems: 64,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"allowed_methods": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MaxItems: 64,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"DELETE",
									"GET",
									"HEAD",
									"MERGE",
									"POST",
									"OPTIONS",
									"PUT",
									"PATCH",
								}, false),
							},
						},

						"max_age_in_seconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 2000000000),
						},

						"credentials_allowed": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"configuration_export_storage_account_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceHealthcareApisFhirServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Healthcare Fhir Service creation.")

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	id := fhirservices.NewFhirServiceID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_healthcare_fhir_service", id.ID())
		}
	}

	i, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := fhirservices.FhirService{
		Identity: i,
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     pointer.To(fhirservices.FhirServiceKind(d.Get("kind").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &fhirservices.FhirServiceProperties{
			AuthenticationConfiguration: expandFhirAuthentication(d.Get("authentication").([]interface{})),
			CorsConfiguration:           expandFhirCorsConfiguration(d.Get("cors").([]interface{})),
		},
	}

	accessPolicyObjectIds, hasValues := d.GetOk("access_policy_object_ids")
	if hasValues {
		parameters.Properties.AccessPolicies = expandAccessPolicy(accessPolicyObjectIds.(*pluginsdk.Set).List())
	}

	storageAcc, hasValues := d.GetOk("configuration_export_storage_account_name")
	if hasValues {
		parameters.Properties.ExportConfiguration = &fhirservices.FhirServiceExportConfiguration{
			StorageAccountName: pointer.To(storageAcc.(string)),
		}
	}

	acrConfig := fhirservices.FhirServiceAcrConfiguration{}
	ociArtifactsRaw, hasValues := d.GetOk("oci_artifact")
	if hasValues {
		ociArtifacts := expandOciArtifacts(ociArtifactsRaw.([]interface{}))
		acrConfig.OciArtifacts = ociArtifacts
	}
	loginServersRaw, hasValues := d.GetOk("container_registry_login_server_url")
	if hasValues {
		loginServers := expandFhirAcrLoginServer(loginServersRaw.(*pluginsdk.Set).List())
		acrConfig.LoginServers = loginServers
	}
	parameters.Properties.AcrConfiguration = &acrConfig

	err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 12,
		Delay:                     60 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"Creating", "Updating", "Verifying"},
		Target:                    []string{"Succeeded"},
		Refresh:                   fhirServiceCreateStateRefreshFunc(ctx, client, id),
		Timeout:                   d.Timeout(pluginsdk.TimeoutUpdate),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Fhir Service %s to settle down: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceHealthcareApisFhirServiceRead(d, meta)
}

func resourceHealthcareApisFhirServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fhirservices.ParseFhirServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.FhirServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	workSpaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("workspace_id", workSpaceId.ID())

	if m := resp.Model; m != nil {
		d.Set("location", location.NormalizeNilable(m.Location))

		i, err := identity.FlattenLegacySystemAndUserAssignedMap(m.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", i); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
		d.Set("kind", string(pointer.From(m.Kind)))

		if props := m.Properties; props != nil {
			d.Set("access_policy_object_ids", flattenFhirAccessPolicy(props.AccessPolicies))
			d.Set("authentication", flattenFhirAuthentication(props.AuthenticationConfiguration))
			d.Set("cors", flattenFhirCorsConfiguration(props.CorsConfiguration))
			d.Set("container_registry_login_server_url", flattenFhirAcrLoginServer(props.AcrConfiguration))
			if acrConfig := props.AcrConfiguration; acrConfig != nil {
				if artifacts := acrConfig.OciArtifacts; artifacts != nil {
					d.Set("oci_artifact", flattenOciArtifacts(artifacts))
				}
			}
			if props.ExportConfiguration != nil && props.ExportConfiguration.StorageAccountName != nil {
				d.Set("configuration_export_storage_account_name", props.ExportConfiguration.StorageAccountName)
			}
			if props.PublicNetworkAccess != nil {
				d.Set("public_network_access_enabled", pointer.From(props.PublicNetworkAccess) == fhirservices.PublicNetworkAccessEnabled)
			}

			return tags.FlattenAndSet(d, m.Tags)
		}
	}
	return nil
}

func expandOciArtifacts(input []interface{}) *[]fhirservices.ServiceOciArtifactEntry {
	output := make([]fhirservices.ServiceOciArtifactEntry, 0)

	for _, artifactSet := range input {
		artifactRaw := artifactSet.(map[string]interface{})

		loginServer := artifactRaw["login_server"].(string)
		artifact := fhirservices.ServiceOciArtifactEntry{
			LoginServer: &loginServer,
			ImageName:   nil,
			Digest:      nil,
		}
		if image := artifactRaw["image_name"].(string); image != "" {
			artifact.ImageName = &image
		}
		if digest := artifactRaw["digest"].(string); digest != "" {
			artifact.Digest = &digest
		}

		output = append(output, artifact)
	}

	return &output
}

func resourceHealthcareApisFhirServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspace, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	id := fhirservices.NewFhirServiceID(workspace.SubscriptionId, workspace.ResourceGroupName, workspace.WorkspaceName, d.Get("name").(string))

	i, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := fhirservices.FhirService{
		Identity: i,
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     pointer.To(fhirservices.FhirServiceKind(d.Get("kind").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &fhirservices.FhirServiceProperties{
			AuthenticationConfiguration: expandFhirAuthentication(d.Get("authentication").([]interface{})),
			CorsConfiguration:           expandFhirCorsConfiguration(d.Get("cors").([]interface{})),
			AccessPolicies:              expandAccessPolicy(d.Get("access_policy_object_ids").(*pluginsdk.Set).List()),
		},
	}

	storageAcc, hasValues := d.GetOk("configuration_export_storage_account_name")
	if hasValues {
		parameters.Properties.ExportConfiguration = &fhirservices.FhirServiceExportConfiguration{
			StorageAccountName: pointer.To(storageAcc.(string)),
		}
	}

	acrConfig := fhirservices.FhirServiceAcrConfiguration{}
	ociArtifactsRaw, hasValues := d.GetOk("oci_artifact")
	if hasValues {
		ociArtifacts := expandOciArtifacts(ociArtifactsRaw.([]interface{}))
		acrConfig.OciArtifacts = ociArtifacts
	}
	loginServersRaw, hasValues := d.GetOk("container_registry_login_server_url")
	if hasValues {
		loginServers := expandFhirAcrLoginServer(loginServersRaw.(*pluginsdk.Set).List())
		acrConfig.LoginServers = loginServers
	}
	parameters.Properties.AcrConfiguration = &acrConfig

	err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceHealthcareApisFhirServiceRead(d, meta)
}

func resourceHealthcareApisFhirServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fhirservices.ParseFhirServiceID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to be deleted..", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Pending"},
		Target:                    []string{"Deleted"},
		Refresh:                   fhirServiceStateStatusCodeRefreshFunc(ctx, client, *id),
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
		ContinuousTargetOccurence: 3,
		PollInterval:              10 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}
	return nil
}

func fhirServiceStateStatusCodeRefreshFunc(ctx context.Context, client *fhirservices.FhirServicesClient, id fhirservices.FhirServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return res, "Deleted", nil
			}
			return nil, "Error", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}

		return res, "Pending", nil
	}
}

func expandFhirAuthentication(input []interface{}) *fhirservices.FhirServiceAuthenticationConfiguration {
	authConfig := input[0].(map[string]interface{})
	authority := authConfig["authority"].(string)
	audience := authConfig["audience"].(string)
	smartProxyEnabled := authConfig["smart_proxy_enabled"].(bool)

	auth := &fhirservices.FhirServiceAuthenticationConfiguration{
		Authority:         pointer.To(authority),
		Audience:          pointer.To(audience),
		SmartProxyEnabled: pointer.To(smartProxyEnabled),
	}

	return auth
}

func expandAccessPolicy(input []interface{}) *[]fhirservices.FhirServiceAccessPolicyEntry {
	if len(input) == 0 {
		return nil
	}

	accessPolicySet := make([]fhirservices.FhirServiceAccessPolicyEntry, 0)

	for _, objectId := range input {
		accessPolicyObjectId := fhirservices.FhirServiceAccessPolicyEntry{
			ObjectId: objectId.(string),
		}
		accessPolicySet = append(accessPolicySet, accessPolicyObjectId)
	}

	return &accessPolicySet
}

func expandFhirCorsConfiguration(input []interface{}) *fhirservices.FhirServiceCorsConfiguration {
	if len(input) == 0 {
		return &fhirservices.FhirServiceCorsConfiguration{
			Origins:          &[]string{},
			Headers:          &[]string{},
			Methods:          &[]string{},
			AllowCredentials: pointer.To(false),
		}
	}

	block := input[0].(map[string]interface{})

	allowedOrigins := *utils.ExpandStringSlice(block["allowed_origins"].(*pluginsdk.Set).List())
	allowedHeaders := *utils.ExpandStringSlice(block["allowed_headers"].(*pluginsdk.Set).List())
	allowedMethods := *utils.ExpandStringSlice(block["allowed_methods"].(*pluginsdk.Set).List())
	allowCredentials := block["credentials_allowed"].(bool)

	cors := &fhirservices.FhirServiceCorsConfiguration{
		Origins:          &allowedOrigins,
		Headers:          &allowedHeaders,
		Methods:          &allowedMethods,
		AllowCredentials: &allowCredentials,
	}

	if v, ok := block["max_age_in_seconds"]; ok {
		maxAgeInSeconds := int64(v.(int))
		cors.MaxAge = &maxAgeInSeconds
	}

	return cors
}

func expandFhirAcrLoginServer(input []interface{}) *[]string {
	acrLoginServers := make([]string, 0)

	if len(input) == 0 {
		return &acrLoginServers
	}

	for _, item := range input {
		acrLoginServers = append(acrLoginServers, item.(string))
	}
	return &acrLoginServers
}

func flattenFhirAcrLoginServer(acrConfig *fhirservices.FhirServiceAcrConfiguration) []string {
	result := make([]string, 0)
	if acrConfig == nil {
		return result
	}

	if loginServer := acrConfig.LoginServers; loginServer != nil {
		result = append(result, *loginServer...)
	}
	return result
}

func flattenFhirAccessPolicy(policies *[]fhirservices.FhirServiceAccessPolicyEntry) []string {
	result := make([]string, 0)

	if policies == nil {
		return result
	}

	for _, policy := range *policies {
		result = append(result, policy.ObjectId)
	}
	return result
}

func flattenOciArtifacts(artifacts *[]fhirservices.ServiceOciArtifactEntry) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if artifacts == nil {
		return result
	}
	for _, artifact := range *artifacts {
		artifactRaw := make(map[string]interface{})

		if loginServer := artifact.LoginServer; loginServer != nil {
			artifactRaw["login_server"] = *loginServer
		}
		if imageName := artifact.ImageName; imageName != nil {
			artifactRaw["image_name"] = *imageName
		}
		if digest := artifact.Digest; digest != nil {
			artifactRaw["digest"] = *digest
		}
		result = append(result, artifactRaw)
	}

	return result
}

func flattenFhirCorsConfiguration(corsConfig *fhirservices.FhirServiceCorsConfiguration) []interface{} {
	if corsConfig == nil {
		return []interface{}{}
	}

	if corsConfig.Origins != nil && len(*corsConfig.Origins) == 0 &&
		corsConfig.Methods != nil && len(*corsConfig.Methods) == 0 &&
		corsConfig.Headers != nil && len(*corsConfig.Headers) == 0 &&
		corsConfig.AllowCredentials != nil && !*corsConfig.AllowCredentials {
		return []interface{}{}
	}

	var maxAge int
	if corsConfig.MaxAge != nil {
		maxAge = int(*corsConfig.MaxAge)
	}

	allowCredentials := false
	if corsConfig.AllowCredentials != nil {
		allowCredentials = *corsConfig.AllowCredentials
	}

	return []interface{}{
		map[string]interface{}{
			"credentials_allowed": allowCredentials,
			"allowed_headers":     utils.FlattenStringSlice(corsConfig.Headers),
			"allowed_methods":     utils.FlattenStringSlice(corsConfig.Methods),
			"allowed_origins":     utils.FlattenStringSlice(corsConfig.Origins),
			"max_age_in_seconds":  maxAge,
		},
	}
}

func flattenFhirAuthentication(authConfig *fhirservices.FhirServiceAuthenticationConfiguration) []interface{} {
	if authConfig == nil {
		return []interface{}{}
	}

	authority := ""
	if authConfig.Authority != nil {
		authority = *authConfig.Authority
	}

	audience := ""
	if authConfig.Audience != nil {
		audience = *authConfig.Audience
	}

	smartProxyEnabled := false
	if authConfig.SmartProxyEnabled != nil {
		smartProxyEnabled = *authConfig.SmartProxyEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"audience":            audience,
			"authority":           authority,
			"smart_proxy_enabled": smartProxyEnabled,
		},
	}
}

func fhirServiceCreateStateRefreshFunc(ctx context.Context, client *fhirservices.FhirServicesClient, fhirServiceId fhirservices.FhirServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, fhirServiceId)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return nil, "", fmt.Errorf("unable to retrieve iot connector %q: %+v", fhirServiceId, err)
			}
			return nil, "Error", fmt.Errorf("polling for the status of %s: %+v", fhirServiceId, err)
		}

		if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.ProvisioningState == nil {
			return resp, "Error", fmt.Errorf("model or properties or ProvisioningState is nil")
		}

		return resp, string(*resp.Model.Properties.ProvisioningState), nil
	}
}
