package healthcare

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	acrParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	fhirService "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/fhirservices"
	workspace "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/workspaces"
	storgaeAccountParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	storageAccountValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHealthcareApisFhirService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareApisFhirServiceCreateUpdate,
		Read:   resourceHealthcareApisFhirServiceRead,
		Update: resourceHealthcareApisFhirServiceCreateUpdate,
		Delete: resourceHealthcareApisFhirServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := fhirService.ParseFhirServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				//todo: check the validate function
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workspace.ValidateWorkspaceID,
			},

			"location": azure.SchemaLocation(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				//todo: confirm the default value later
				Default: string(fhirService.FhirServiceKindFhirNegativeRFour),
				ValidateFunc: validation.StringInSlice([]string{
					string(fhirService.FhirServiceKindFhirNegativeRFour),
					string(fhirService.FhirServiceKindFhirNegativeStuThree),
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

			"authentication_configuration": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authority": {
							Type:     pluginsdk.TypeString,
							Required: true,
							//todo: must follow https://login.microsoft.com/tenantid
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

			"identity": {
				Type: pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type":{
							Type: pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(fhirService.ManagedServiceIdentityTypeSystemAssigned),
							},false),
						},
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"acr_login_servers": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"cors_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"allowed_origins": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							MaxItems: 64,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
						"allowed_headers": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							MaxItems: 64,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
						"allowed_methods": {
							//check set or list via api call
							Type:     pluginsdk.TypeList,
							Optional: true,
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
						"allow_credentials": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"export_storage_account_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: storageAccountValidate.StorageAccountID,
			},
			"tags": commonschema.Tags(),
		},
	}

}

func resourceHealthcareApisFhirServiceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Healthcare Fhir Service creation.")

	//todo check other resource about this ID to Name practice
	workspace, err := workspace.ParseWorkspaceIDInsensitively(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	fhirServiceId := fhirService.NewFhirServiceID(workspace.SubscriptionId, workspace.ResourceGroupName, workspace.WorkspaceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, fhirServiceId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", fhirServiceId, err)
			}
		}

		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_healthcareapis_fhir_service", fhirServiceId.ID())
		}
	}

	parameters := fhirService.FhirService{
		Name:     utils.String(fhirServiceId.FhirServiceName),
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Identity: expandFhirServiceIdentity(d.Get("identity").([]interface{})),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &fhirService.FhirServiceProperties{
			AuthenticationConfiguration: expandFhirAuthentication(d.Get("authentication_configuration").([]interface{})),
			CorsConfiguration:           expandFhirCorsConfiguration(d.Get("cors_configuration").([]interface{})),
		},
	}

	kind := fhirService.FhirServiceKind(d.Get("kind").(string))
	parameters.Kind = &kind

	accessPolicyObjectIds, hasValues := d.GetOk("access_policy_object_ids")
	if hasValues {
		parameters.Properties.AccessPolicies = expandAccessPolicy(accessPolicyObjectIds.(*pluginsdk.Set).List())
	}

	if v := d.Get("export_storage_account_id").(string); v != "" {
		storageAccount, err := storgaeAccountParse.StorageAccountID(v)
		if err != nil {
			return err
		}
		parameters.Properties.ExportConfiguration = &fhirService.FhirServiceExportConfiguration{
			StorageAccountName: utils.String(storageAccount.Name),
		}
	}

	acrConfig, hasValues := d.GetOk("acr_login_servers")
	if hasValues {
		result, err := expandFhirAcrLoginServer(acrConfig.(*pluginsdk.Set).List(), meta, ctx)
		if err != nil {
			return fmt.Errorf("expanding Acr Login Server error: %+v", err)
		}
		parameters.Properties.AcrConfiguration = result
	}

	if err := client.CreateOrUpdateThenPoll(ctx, fhirServiceId, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", fhirServiceId, err)
	}

	d.SetId(fhirServiceId.ID())
	return resourceHealthcareApisFhirServiceRead(d, meta)
}

func resourceHealthcareApisFhirServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fhirService.ParseFhirServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.FhirServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if model.Kind != nil {
			d.Set("kind", model.Kind)
		}

		if props := model.Properties; props != nil {
			d.Set("access_policy_object_ids", flattenFhirAccessPolicy(props.AccessPolicies))
			d.Set("authentication_configuration", flattenFhirAuthentication(props.AuthenticationConfiguration))
			d.Set("acr_login_servers", flattenFhirAcrLoginServer(props.AcrConfiguration))
			d.Set("cors_configuration", flattenFhirCorsConfiguration(props.CorsConfiguration))
			if props.ExportConfiguration != nil && props.ExportConfiguration.StorageAccountName != nil {
				d.Set("export_storage_account_id", props.ExportConfiguration.StorageAccountName)
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceHealthcareApisFhirServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fhirService.ParseFhirServiceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(future.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return waitForHealthcareApiFhirServiceToBeDeleted(ctx, client, *id)
}

func waitForHealthcareApiFhirServiceToBeDeleted(ctx context.Context, client *fhirService.FhirServicesClient, id fhirService.FhirServiceId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context has no deadline")
	}

	log.Printf("[DEBUG] Waiting for %s to be deleted..", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: healthcareApiFhirServiceStateCodeRefreshFunc(ctx, client, id),
		Timeout: time.Until(deadline),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}
	return nil
}

func healthcareApiFhirServiceStateCodeRefreshFunc(ctx context.Context, client *fhirService.FhirServicesClient, id fhirService.FhirServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if res.HttpResponse != nil {
			log.Printf("Retrieving %s returned Status %d", id, res.HttpResponse.StatusCode)
		}

		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
			}
			return nil, "", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}

		return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
	}
}

func expandFhirServiceIdentity(input []interface{}) *fhirService.ServiceManagedIdentityIdentity{
	if len(input) == 0 {
		return &fhirService.ServiceManagedIdentityIdentity{
			Type: fhirService.ManagedServiceIdentityTypeNone,
		}
	}

		identity := input[0].(map[string]interface{})
		return &fhirService.ServiceManagedIdentityIdentity{
			Type: fhirService.ManagedServiceIdentityType(identity["type"].(string)),
		}
	}
}
func expandFhirAuthentication(input []interface{}) *fhirService.FhirServiceAuthenticationConfiguration {
	authConfig := input[0].(map[string]interface{})
	authority := authConfig["authority"].(string)
	audience := authConfig["audience"].(string)
	smartProxyEnabled := authConfig["smart_proxy_enabled"].(bool)

	auth := &fhirService.FhirServiceAuthenticationConfiguration{
		Authority:         utils.String(authority),
		Audience:          utils.String(audience),
		SmartProxyEnabled: utils.Bool(smartProxyEnabled),
	}

	return auth
}
func expandAccessPolicy(input []interface{}) *[]fhirService.FhirServiceAccessPolicyEntry {
	if len(input) == 0 {
		return nil
	}

	accessPolicySet := make([]fhirService.FhirServiceAccessPolicyEntry, 0)

	for _, objectId := range input {
		accessPolicyObjectId := fhirService.FhirServiceAccessPolicyEntry{
			ObjectId: objectId.(string),
		}
		accessPolicySet = append(accessPolicySet, accessPolicyObjectId)
	}

	return &accessPolicySet
}

func expandFhirCorsConfiguration(input []interface{}) *fhirService.FhirServiceCorsConfiguration {
	if len(input) == 0 {
		return &fhirService.FhirServiceCorsConfiguration{
			Origins:          &[]string{},
			Headers:          &[]string{},
			Methods:          &[]string{},
			AllowCredentials: utils.Bool(false),
		}
	}

	block := input[0].(map[string]interface{})

	allowedOrigins := *utils.ExpandStringSlice(block["allowed_origins"].(*pluginsdk.Set).List())
	allowedHeaders := *utils.ExpandStringSlice(block["allowed_headers"].(*pluginsdk.Set).List())
	allowedMethods := *utils.ExpandStringSlice(block["allowed_methods"].([]interface{}))
	allowCredentials := block["allow_credentials"].(bool)

	cors := &fhirService.FhirServiceCorsConfiguration{
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

func expandFhirAcrLoginServer(input []interface{}, meta interface{}, ctx context.Context) (*fhirService.FhirServiceAcrConfiguration, error) {
	acrLoginServers := make([]string, 0)
	acrClient := meta.(*clients.Client).Containers.RegistriesClient

	for _, item := range input {
		acrId, err := acrParse.RegistryID(item.(string))
		if err != nil {
			return nil, err
		}
		acrItem, err := acrClient.Get(ctx, acrId.ResourceGroup, acrId.Name)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", acrId, err)
		}
		if loginServer := acrItem.LoginServer; loginServer != nil {
			acrLoginServers = append(acrLoginServers, *loginServer)
		}
	}

	return &fhirService.FhirServiceAcrConfiguration{
		LoginServers: &acrLoginServers,
	}, nil
}

func flattenFhirAcrLoginServer(acrLoginServer *fhirService.FhirServiceAcrConfiguration) []string {
	result := make([]string, 0)
	if acrLoginServer == nil {
		return result
	}

	if loginServer := acrLoginServer.LoginServers; loginServer != nil {
		for _, serverId := range *loginServer {
			result = append(result, serverId)
		}
	}
	return result
}

func flattenFhirAccessPolicy(policies *[]fhirService.FhirServiceAccessPolicyEntry) []string {
	result := make([]string, 0)

	if policies == nil {
		return result
	}

	for _, policy := range *policies {
		if objectId := policy.ObjectId; objectId != "" {
			result = append(result, objectId)
		}
	}
	return result
}

func flattenFhirCorsConfiguration(corsConfig *fhirService.FhirServiceCorsConfiguration) []interface{} {
	if corsConfig == nil {
		return []interface{}{}
	}

	if corsConfig.Origins != nil && len(*corsConfig.Origins) == 0 &&
		corsConfig.Methods != nil && len(*corsConfig.Methods) == 0 &&
		corsConfig.Headers != nil && len(*corsConfig.Headers) == 0 &&
		corsConfig.AllowCredentials != nil && *corsConfig.AllowCredentials == false {
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
			"allow_credentials":  allowCredentials,
			"allowed_headers":    utils.FlattenStringSlice(corsConfig.Headers),
			"allowed_methods":    utils.FlattenStringSlice(corsConfig.Methods),
			"allowed_origins":    utils.FlattenStringSlice(corsConfig.Origins),
			"max_age_in_seconds": maxAge,
		},
	}
}

func flattenFhirAuthentication(authConfig *fhirService.FhirServiceAuthenticationConfiguration) []interface{} {
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
