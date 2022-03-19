package healthcare

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2021-11-01/healthcareapis"
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
			_, err := parse.FhirServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
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
				Default: string(healthcareapis.KindFhirR4),
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

			"authentication_configuration": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authority": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(healthcareapis.ManagedServiceIdentityTypeSystemAssigned),
							}, false),
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

			// can't use the registry ID due to the ID cannot be obtained when setting the property in state file
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
							Type:     pluginsdk.TypeSet,
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

			"export_storage_account_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
			return tf.ImportAsExistsError("azurerm_healthcareapis_fhir_service", fhirServiceId.ID())
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
			AuthenticationConfiguration: expandFhirAuthentication(d.Get("authentication_configuration").([]interface{})),
			CorsConfiguration:           expandFhirCorsConfiguration(d.Get("cors_configuration").([]interface{})),
		},
	}

	accessPolicyObjectIds, hasValues := d.GetOk("access_policy_object_ids")
	if hasValues {
		parameters.FhirServiceProperties.AccessPolicies = expandAccessPolicy(accessPolicyObjectIds.(*pluginsdk.Set).List())
	}

	if v := d.Get("export_storage_account_name").(string); v != "" {
		//storageAccount, err := storgaeAccountParse.StorageAccountID(v)
		//if err != nil {
		//	return err
		//}
		parameters.FhirServiceProperties.ExportConfiguration = &healthcareapis.FhirServiceExportConfiguration{
			StorageAccountName: utils.String(v),
		}
	}

	acrConfig, hasValues := d.GetOk("acr_login_servers")
	if hasValues {
		result := expandFhirAcrLoginServer(acrConfig.(*pluginsdk.Set).List())
		parameters.FhirServiceProperties.AcrConfiguration = result
	}

	future, err := client.CreateOrUpdate(ctx, fhirServiceId.ResourceGroup, fhirServiceId.WorkspaceName, fhirServiceId.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", fhirServiceId, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", fhirServiceId, err)
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
		d.Set("authentication_configuration", flattenFhirAuthentication(props.AuthenticationConfiguration))
		d.Set("cors_configuration", flattenFhirCorsConfiguration(props.CorsConfiguration))
		d.Set("acr_login_servers", flattenFhirAcrLoginServer(props.AcrConfiguration))
		if props.ExportConfiguration != nil && props.ExportConfiguration.StorageAccountName != nil {
			d.Set("export_storage_account_name", props.ExportConfiguration.StorageAccountName)
		}

		if err := tags.FlattenAndSet(d, resp.Tags); err != nil {
			return err
		}
	}

	return nil
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
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}
	return nil
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
	allowCredentials := block["allow_credentials"].(bool)

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

//func expandFhirAcrLoginServer(input []interface{}, meta interface{}, ctx context.Context) (*fhirService.FhirServiceAcrConfiguration, error) {
//	acrLoginServers := make([]string, 0)
//	acrClient := meta.(*clients.Client).Containers.RegistriesClient
//
//	for _, item := range input {
//		acrId, err := acrParse.RegistryID(item.(string))
//		if err != nil {
//			return nil, err
//		}
//		acrItem, err := acrClient.Get(ctx, acrId.ResourceGroup, acrId.Name)
//		if err != nil {
//			return nil, fmt.Errorf("retrieving %s: %+v", acrId, err)
//		}
//		if loginServer := acrItem.LoginServer; loginServer != nil {
//			acrLoginServers = append(acrLoginServers, *loginServer)
//		}
//	}
//
//	return &fhirService.FhirServiceAcrConfiguration{
//		LoginServers: &acrLoginServers,
//	}, nil
//}

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
		for _, serverId := range *loginServer {
			result = append(result, serverId)
		}
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
