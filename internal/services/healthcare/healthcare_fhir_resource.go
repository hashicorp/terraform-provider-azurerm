package healthcare

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2021-11-01/healthcareapis"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FhirServiceID(id)
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
				ValidateFunc: validate.WorkspaceID,
			},

			"location": commonschema.Location(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(healthcareapis.KindFhirR4),
				ValidateFunc: validation.StringInSlice([]string{
					string(healthcareapis.KindFhirR4),
					string(healthcareapis.KindFhirStu3),
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

			"identity": commonschema.SystemAssignedIdentityOptional(),

			// can't use the registry ID due to the ID cannot be obtained when setting the property in state file
			"container_registry_login_server_url": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
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

			"tags": commonschema.Tags(),
		},
	}

}

func resourceHealthcareApisFhirServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Healthcare Fhir Service creation.")

	workspace, err := parse.WorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	fhirServiceId := parse.NewFhirServiceID(workspace.SubscriptionId, workspace.ResourceGroup, workspace.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, fhirServiceId.ResourceGroup, fhirServiceId.WorkspaceName, fhirServiceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", fhirServiceId, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_healthcare_fhir_service", fhirServiceId.ID())
		}
	}

	identity, err := expandFhirManagedIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := healthcareapis.FhirService{
		Identity: identity,
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     healthcareapis.FhirServiceKind(d.Get("kind").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		FhirServiceProperties: &healthcareapis.FhirServiceProperties{
			AuthenticationConfiguration: expandFhirAuthentication(d.Get("authentication").([]interface{})),
			CorsConfiguration:           expandFhirCorsConfiguration(d.Get("cors").([]interface{})),
		},
	}

	accessPolicyObjectIds, hasValues := d.GetOk("access_policy_object_ids")
	if hasValues {
		parameters.FhirServiceProperties.AccessPolicies = expandAccessPolicy(accessPolicyObjectIds.(*pluginsdk.Set).List())
	}

	storageAcc, hasValues := d.GetOk("configuration_export_storage_account_name")
	if hasValues {
		parameters.FhirServiceProperties.ExportConfiguration = &healthcareapis.FhirServiceExportConfiguration{
			StorageAccountName: utils.String(storageAcc.(string)),
		}
	}

	acrConfig, hasValues := d.GetOk("container_registry_login_server_url")
	if hasValues {
		result := expandFhirAcrLoginServer(acrConfig.(*pluginsdk.Set).List())
		parameters.FhirServiceProperties.AcrConfiguration = result
	}

	future, err := client.CreateOrUpdate(ctx, fhirServiceId.ResourceGroup, fhirServiceId.WorkspaceName, fhirServiceId.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", fhirServiceId, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", fhirServiceId, err)
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 12,
		Delay:                     60 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"Creating", "Updating", "Verifying"},
		Target:                    []string{"Succeeded"},
		Refresh:                   fhirServiceCreateStateRefreshFunc(ctx, client, fhirServiceId),
		Timeout:                   d.Timeout(pluginsdk.TimeoutUpdate),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Fhir Service %s to settle down: %+v", fhirServiceId, err)
	}

	d.SetId(fhirServiceId.ID())
	return resourceHealthcareApisFhirServiceRead(d, meta)
}

func resourceHealthcareApisFhirServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FhirServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	workSpaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("workspace_id", workSpaceId.ID())

	if resp.Location != nil {
		d.Set("location", location.NormalizeNilable(resp.Location))
	}

	if err := d.Set("identity", flattenFhirManagedIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	d.Set("kind", resp.Kind)

	if props := resp.FhirServiceProperties; props != nil {
		d.Set("access_policy_object_ids", flattenFhirAccessPolicy(props.AccessPolicies))
		d.Set("authentication", flattenFhirAuthentication(props.AuthenticationConfiguration))
		d.Set("cors", flattenFhirCorsConfiguration(props.CorsConfiguration))
		d.Set("container_registry_login_server_url", flattenFhirAcrLoginServer(props.AcrConfiguration))
		if props.ExportConfiguration != nil && props.ExportConfiguration.StorageAccountName != nil {
			d.Set("configuration_export_storage_account_name", props.ExportConfiguration.StorageAccountName)
		}

		if err := tags.FlattenAndSet(d, resp.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceHealthcareApisFhirServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspace, err := parse.WorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	fhirServiceId := parse.NewFhirServiceID(workspace.SubscriptionId, workspace.ResourceGroup, workspace.Name, d.Get("name").(string))

	identity, err := expandFhirManagedIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := healthcareapis.FhirService{
		Identity: identity,
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     healthcareapis.FhirServiceKind(d.Get("kind").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		FhirServiceProperties: &healthcareapis.FhirServiceProperties{
			AuthenticationConfiguration: expandFhirAuthentication(d.Get("authentication").([]interface{})),
			CorsConfiguration:           expandFhirCorsConfiguration(d.Get("cors").([]interface{})),
			AccessPolicies:              expandAccessPolicy(d.Get("access_policy_object_ids").(*pluginsdk.Set).List()),
		},
	}

	storageAcc, hasValues := d.GetOk("configuration_export_storage_account_name")
	if hasValues {
		parameters.FhirServiceProperties.ExportConfiguration = &healthcareapis.FhirServiceExportConfiguration{
			StorageAccountName: utils.String(storageAcc.(string)),
		}
	}

	acrConfig, hasValues := d.GetOk("container_registry_login_server_url")
	if hasValues {
		result := expandFhirAcrLoginServer(acrConfig.(*pluginsdk.Set).List())
		parameters.FhirServiceProperties.AcrConfiguration = result
	}

	future, err := client.CreateOrUpdate(ctx, fhirServiceId.ResourceGroup, fhirServiceId.WorkspaceName, fhirServiceId.Name, parameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", fhirServiceId, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", fhirServiceId, err)
	}

	d.SetId(fhirServiceId.ID())
	return resourceHealthcareApisFhirServiceRead(d, meta)
}

func resourceHealthcareApisFhirServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FhirServiceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, id.WorkspaceName)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
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

func fhirServiceStateStatusCodeRefreshFunc(ctx context.Context, client *healthcareapis.FhirServicesClient, id parse.FhirServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, "Deleted", nil
			}
			return nil, "Error", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}

		return res, "Pending", nil
	}
}

func expandFhirManagedIdentity(input []interface{}) (*healthcareapis.ServiceManagedIdentityIdentity, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &healthcareapis.ServiceManagedIdentityIdentity{
		Type: healthcareapis.ServiceManagedIdentityType(string(expanded.Type)),
	}, nil
}

func flattenFhirManagedIdentity(input *healthcareapis.ServiceManagedIdentityIdentity) []interface{} {
	var transition *identity.SystemAssigned

	if input != nil {
		transition = &identity.SystemAssigned{
			Type: identity.Type(string(input.Type)),
		}
		if input.PrincipalID != nil {
			principalID := *input.PrincipalID
			transition.PrincipalId = principalID.String()
		}
		if input.TenantID != nil {
			tenantID := *input.TenantID
			transition.TenantId = tenantID.String()
		}
	}

	return identity.FlattenSystemAssigned(transition)
}

func expandFhirAuthentication(input []interface{}) *healthcareapis.FhirServiceAuthenticationConfiguration {
	authConfig := input[0].(map[string]interface{})
	authority := authConfig["authority"].(string)
	audience := authConfig["audience"].(string)
	smartProxyEnabled := authConfig["smart_proxy_enabled"].(bool)

	auth := &healthcareapis.FhirServiceAuthenticationConfiguration{
		Authority:         utils.String(authority),
		Audience:          utils.String(audience),
		SmartProxyEnabled: utils.Bool(smartProxyEnabled),
	}

	return auth
}
func expandAccessPolicy(input []interface{}) *[]healthcareapis.FhirServiceAccessPolicyEntry {
	if len(input) == 0 {
		return nil
	}

	accessPolicySet := make([]healthcareapis.FhirServiceAccessPolicyEntry, 0)

	for _, objectId := range input {
		accessPolicyObjectId := healthcareapis.FhirServiceAccessPolicyEntry{
			ObjectID: utils.String(objectId.(string)),
		}
		accessPolicySet = append(accessPolicySet, accessPolicyObjectId)
	}

	return &accessPolicySet
}

func expandFhirCorsConfiguration(input []interface{}) *healthcareapis.FhirServiceCorsConfiguration {
	if len(input) == 0 {
		return &healthcareapis.FhirServiceCorsConfiguration{
			Origins:          &[]string{},
			Headers:          &[]string{},
			Methods:          &[]string{},
			AllowCredentials: utils.Bool(false),
		}
	}

	block := input[0].(map[string]interface{})

	allowedOrigins := *utils.ExpandStringSlice(block["allowed_origins"].(*pluginsdk.Set).List())
	allowedHeaders := *utils.ExpandStringSlice(block["allowed_headers"].(*pluginsdk.Set).List())
	allowedMethods := *utils.ExpandStringSlice(block["allowed_methods"].(*pluginsdk.Set).List())
	allowCredentials := block["credentials_allowed"].(bool)

	cors := &healthcareapis.FhirServiceCorsConfiguration{
		Origins:          &allowedOrigins,
		Headers:          &allowedHeaders,
		Methods:          &allowedMethods,
		AllowCredentials: &allowCredentials,
	}

	if v, ok := block["max_age_in_seconds"]; ok {
		maxAgeInSeconds := int32(v.(int))
		cors.MaxAge = &maxAgeInSeconds
	}

	return cors
}

func expandFhirAcrLoginServer(input []interface{}) *healthcareapis.FhirServiceAcrConfiguration {
	acrLoginServers := make([]string, 0)

	if len(input) == 0 {
		return &healthcareapis.FhirServiceAcrConfiguration{
			LoginServers: &acrLoginServers,
		}
	}

	for _, item := range input {
		acrLoginServers = append(acrLoginServers, item.(string))
	}
	return &healthcareapis.FhirServiceAcrConfiguration{
		LoginServers: &acrLoginServers,
	}
}

func flattenFhirAcrLoginServer(acrLoginServer *healthcareapis.FhirServiceAcrConfiguration) []string {
	result := make([]string, 0)
	if acrLoginServer == nil {
		return result
	}

	if loginServer := acrLoginServer.LoginServers; loginServer != nil {
		result = append(result, *loginServer...)
	}
	return result
}

func flattenFhirAccessPolicy(policies *[]healthcareapis.FhirServiceAccessPolicyEntry) []string {
	result := make([]string, 0)

	if policies == nil {
		return result
	}

	for _, policy := range *policies {
		if objectId := policy.ObjectID; objectId != nil {
			result = append(result, *objectId)
		}
	}
	return result
}

func flattenFhirCorsConfiguration(corsConfig *healthcareapis.FhirServiceCorsConfiguration) []interface{} {
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

func flattenFhirAuthentication(authConfig *healthcareapis.FhirServiceAuthenticationConfiguration) []interface{} {
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

func fhirServiceCreateStateRefreshFunc(ctx context.Context, client *healthcareapis.FhirServicesClient, fhirServiceId parse.FhirServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, fhirServiceId.ResourceGroup, fhirServiceId.WorkspaceName, fhirServiceId.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil, "", fmt.Errorf("unable to retrieve iot connector %q: %+v", fhirServiceId, err)
			}
			return nil, "Error", fmt.Errorf("polling for the status of %s: %+v", fhirServiceId, err)
		}

		return resp, string(resp.ProvisioningState), nil
	}
}
