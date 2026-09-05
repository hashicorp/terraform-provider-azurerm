// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      string(fhirservices.FhirServiceKindFhirNegativeRFour),
				ValidateFunc: validation.StringInSlice(fhirservices.PossibleValuesForFhirServiceKind(), false),
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

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	id := fhirservices.NewFhirServiceID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
		Kind:     pointer.ToEnum[fhirservices.FhirServiceKind](d.Get("kind").(string)),
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
		acrConfig.OciArtifacts = expandOciArtifacts(ociArtifactsRaw.([]interface{}))
	}
	loginServersRaw, hasValues := d.GetOk("container_registry_login_server_url")
	if hasValues {
		acrConfig.LoginServers = expandFhirAcrLoginServer(loginServersRaw.(*pluginsdk.Set).List())
	}
	parameters.Properties.AcrConfiguration = &acrConfig

	if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
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

			if err := tags.FlattenAndSet(d, m.Tags); err != nil {
				return err
			}
		}
	}
	return nil
}

func expandOciArtifacts(input []interface{}) *[]fhirservices.ServiceOciArtifactEntry {
	output := make([]fhirservices.ServiceOciArtifactEntry, 0)

	for _, artifactSet := range input {
		artifactRaw := artifactSet.(map[string]interface{})

		artifact := fhirservices.ServiceOciArtifactEntry{
			LoginServer: pointer.To(artifactRaw["login_server"].(string)),
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
		Kind:     pointer.ToEnum[fhirservices.FhirServiceKind](d.Get("kind").(string)),
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
		acrConfig.OciArtifacts = expandOciArtifacts(ociArtifactsRaw.([]interface{}))
	}
	loginServersRaw, hasValues := d.GetOk("container_registry_login_server_url")
	if hasValues {
		acrConfig.LoginServers = expandFhirAcrLoginServer(loginServersRaw.(*pluginsdk.Set).List())
	}
	parameters.Properties.AcrConfiguration = &acrConfig

	if err = client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
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

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to be deleted..", id)
	poller := custompollers.NewEventualConsistencyPoller(3, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, *id)
		return resp.HttpResponse, err
	}, custompollers.DefaultDeletionEventualConsistencyPollerOptions())
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}
	return nil
}

func expandFhirAuthentication(input []interface{}) *fhirservices.FhirServiceAuthenticationConfiguration {
	authConfig := input[0].(map[string]interface{})
	authority := authConfig["authority"].(string)
	audience := authConfig["audience"].(string)
	smartProxyEnabled := authConfig["smart_proxy_enabled"].(bool)

	return &fhirservices.FhirServiceAuthenticationConfiguration{
		Authority:         pointer.To(authority),
		Audience:          pointer.To(audience),
		SmartProxyEnabled: pointer.To(smartProxyEnabled),
	}
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

	allowedOrigins := *helpers.ExpandStringSlice(block["allowed_origins"].(*pluginsdk.Set).List())
	allowedHeaders := *helpers.ExpandStringSlice(block["allowed_headers"].(*pluginsdk.Set).List())
	allowedMethods := *helpers.ExpandStringSlice(block["allowed_methods"].(*pluginsdk.Set).List())

	cors := &fhirservices.FhirServiceCorsConfiguration{
		Origins:          &allowedOrigins,
		Headers:          &allowedHeaders,
		Methods:          &allowedMethods,
		AllowCredentials: pointer.To(block["credentials_allowed"].(bool)),
	}

	if v, ok := block["max_age_in_seconds"]; ok {
		cors.MaxAge = pointer.To(int64(v.(int)))
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

	return []interface{}{
		map[string]interface{}{
			"credentials_allowed": pointer.From(corsConfig.AllowCredentials),
			"allowed_headers":     helpers.FlattenStringSlice(corsConfig.Headers),
			"allowed_methods":     helpers.FlattenStringSlice(corsConfig.Methods),
			"allowed_origins":     helpers.FlattenStringSlice(corsConfig.Origins),
			"max_age_in_seconds":  maxAge,
		},
	}
}

func flattenFhirAuthentication(authConfig *fhirservices.FhirServiceAuthenticationConfiguration) []interface{} {
	if authConfig == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"audience":            pointer.From(authConfig.Audience),
			"authority":           pointer.From(authConfig.Authority),
			"smart_proxy_enabled": pointer.From(authConfig.SmartProxyEnabled),
		},
	}
}
