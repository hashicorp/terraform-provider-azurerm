package apimanagement

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementApi() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiCreateUpdate,
		Read:   resourceApiManagementApiRead,
		Update: resourceApiManagementApiCreateUpdate,
		Delete: resourceApiManagementApiDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementApiName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"path": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ApiManagementApiPath,
			},

			"protocols": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(apimanagement.ProtocolHTTP),
						string(apimanagement.ProtocolHTTPS),
					}, false),
				},
			},

			"revision": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Optional
			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"import": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"content_value": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"content_format": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(apimanagement.Openapi),
								string(apimanagement.Openapijson),
								string(apimanagement.OpenapijsonLink),
								string(apimanagement.OpenapiLink),
								string(apimanagement.SwaggerJSON),
								string(apimanagement.SwaggerLinkJSON),
								string(apimanagement.WadlLinkJSON),
								string(apimanagement.WadlXML),
								string(apimanagement.Wsdl),
								string(apimanagement.WsdlLink),
							}, false),
						},

						"wsdl_selector": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"service_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"endpoint_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},

			"service_url": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"subscription_key_parameter_names": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"header": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"query": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"subscription_required": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"soap_pass_through": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"oauth2_authorization": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authorization_server_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.ApiManagementChildName,
						},
						"scope": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							// There is currently no validation, as any length and characters can be used in the field
						},
					},
				},
			},

			"openid_authentication": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"openid_provider_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.ApiManagementChildName,
						},
						"bearer_token_sending_methods": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(apimanagement.BearerTokenSendingMethodsAuthorizationHeader),
									string(apimanagement.BearerTokenSendingMethodsQuery),
								}, false),
							},
						},
					},
				},
			},

			// Computed
			"is_current": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"is_online": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
				Optional: true,
			},

			"version_set_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementApiCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	name := d.Get("name").(string)
	revision := d.Get("revision").(string)
	path := d.Get("path").(string)
	apiId := fmt.Sprintf("%s;rev=%s", name, revision)
	version := d.Get("version").(string)
	versionSetId := d.Get("version_set_id").(string)

	if version != "" && versionSetId == "" {
		return fmt.Errorf("setting `version` without the required `version_set_id`")
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apiId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing API %q (API Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_api", *existing.ID)
		}
	}

	var apiType apimanagement.APIType
	var soapApiType apimanagement.SoapAPIType

	soapPassThrough := d.Get("soap_pass_through").(bool)
	if soapPassThrough {
		apiType = apimanagement.Soap
		soapApiType = apimanagement.SoapPassThrough
	} else {
		apiType = apimanagement.HTTP
		soapApiType = apimanagement.SoapToRest
	}

	// If import is used, we need to send properties to Azure API in two operations.
	// First we execute import and then updated the other props.
	if vs, hasImport := d.GetOk("import"); hasImport {
		importVs := vs.([]interface{})
		importV := importVs[0].(map[string]interface{})
		contentFormat := importV["content_format"].(string)
		contentValue := importV["content_value"].(string)

		log.Printf("[DEBUG] Importing API Management API %q of type %q", name, contentFormat)
		apiParams := apimanagement.APICreateOrUpdateParameter{
			APICreateOrUpdateProperties: &apimanagement.APICreateOrUpdateProperties{
				APIType:     apiType,
				SoapAPIType: soapApiType,
				Format:      apimanagement.ContentFormat(contentFormat),
				Value:       utils.String(contentValue),
				Path:        utils.String(path),
				APIVersion:  utils.String(version),
			},
		}
		wsdlSelectorVs := importV["wsdl_selector"].([]interface{})

		// `wsdl_selector` is necessary under format `wsdl`
		if len(wsdlSelectorVs) == 0 && contentFormat == string(apimanagement.Wsdl) {
			return fmt.Errorf("`wsdl_selector` is required when content format is `wsdl` in API Management API %q", name)
		}

		if len(wsdlSelectorVs) > 0 {
			wsdlSelectorV := wsdlSelectorVs[0].(map[string]interface{})
			wSvcName := wsdlSelectorV["service_name"].(string)
			wEndpName := wsdlSelectorV["endpoint_name"].(string)

			apiParams.APICreateOrUpdateProperties.WsdlSelector = &apimanagement.APICreateOrUpdatePropertiesWsdlSelector{
				WsdlServiceName:  utils.String(wSvcName),
				WsdlEndpointName: utils.String(wEndpName),
			}
		}

		if versionSetId != "" {
			apiParams.APICreateOrUpdateProperties.APIVersionSetID = utils.String(versionSetId)
		}

		future, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiId, apiParams, "")
		if err != nil {
			return fmt.Errorf("creating/updating API Management API %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on creating/updating API Management API %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	serviceUrl := d.Get("service_url").(string)
	subscriptionRequired := d.Get("subscription_required").(bool)

	protocolsRaw := d.Get("protocols").(*pluginsdk.Set).List()
	protocols := expandApiManagementApiProtocols(protocolsRaw)

	subscriptionKeyParameterNamesRaw := d.Get("subscription_key_parameter_names").([]interface{})
	subscriptionKeyParameterNames := expandApiManagementApiSubscriptionKeyParamNames(subscriptionKeyParameterNamesRaw)

	authenticationSettings := &apimanagement.AuthenticationSettingsContract{}

	oAuth2AuthorizationSettingsRaw := d.Get("oauth2_authorization").([]interface{})
	oAuth2AuthorizationSettings := expandApiManagementOAuth2AuthenticationSettingsContract(oAuth2AuthorizationSettingsRaw)
	authenticationSettings.OAuth2 = oAuth2AuthorizationSettings

	openIDAuthorizationSettingsRaw := d.Get("openid_authentication").([]interface{})
	openIDAuthorizationSettings := expandApiManagementOpenIDAuthenticationSettingsContract(openIDAuthorizationSettingsRaw)
	authenticationSettings.Openid = openIDAuthorizationSettings

	params := apimanagement.APICreateOrUpdateParameter{
		APICreateOrUpdateProperties: &apimanagement.APICreateOrUpdateProperties{
			APIType:                       apiType,
			SoapAPIType:                   soapApiType,
			Description:                   utils.String(description),
			DisplayName:                   utils.String(displayName),
			Path:                          utils.String(path),
			Protocols:                     protocols,
			ServiceURL:                    utils.String(serviceUrl),
			SubscriptionKeyParameterNames: subscriptionKeyParameterNames,
			APIVersion:                    utils.String(version),
			SubscriptionRequired:          &subscriptionRequired,
			AuthenticationSettings:        authenticationSettings,
		},
	}

	if versionSetId != "" {
		params.APICreateOrUpdateProperties.APIVersionSetID = utils.String(versionSetId)
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiId, params, "")
	if err != nil {
		return fmt.Errorf("creating/updating API Management API %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating API Management API %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, serviceName, apiId)
	if err != nil {
		return fmt.Errorf("retrieving API %q / Revision %q (API Management Service %q / Resource Group %q): %+v", name, revision, serviceName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for API %q / Revision %q (API Management Service %q / Resource Group %q)", name, revision, serviceName, resourceGroup)
	}

	d.SetId(*read.ID)
	return resourceApiManagementApiRead(d, meta)
}

func resourceApiManagementApiRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	apiid := id.Name

	name := apiid
	revision := ""
	if strings.Contains(apiid, ";") {
		name = strings.Split(apiid, ";")[0]
		revision = strings.Split(apiid, "=")[1]
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiid)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] API %q Revision %q (API Management Service %q / Resource Group %q) does not exist - removing from state!", name, revision, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving API %q / Revision %q (API Management Service %q / Resource Group %q): %+v", name, revision, serviceName, resourceGroup, err)
	}

	d.Set("api_management_name", serviceName)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.APIContractProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("is_current", props.IsCurrent)
		d.Set("is_online", props.IsOnline)
		d.Set("path", props.Path)
		d.Set("service_url", props.ServiceURL)
		d.Set("revision", props.APIRevision)
		d.Set("soap_pass_through", string(props.APIType) == string(apimanagement.SoapPassThrough))
		d.Set("subscription_required", props.SubscriptionRequired)
		d.Set("version", props.APIVersion)
		d.Set("version_set_id", props.APIVersionSetID)

		if err := d.Set("protocols", flattenApiManagementApiProtocols(props.Protocols)); err != nil {
			return fmt.Errorf("setting `protocols`: %s", err)
		}

		if err := d.Set("subscription_key_parameter_names", flattenApiManagementApiSubscriptionKeyParamNames(props.SubscriptionKeyParameterNames)); err != nil {
			return fmt.Errorf("setting `subscription_key_parameter_names`: %+v", err)
		}

		if err := d.Set("oauth2_authorization", flattenApiManagementOAuth2Authorization(props.AuthenticationSettings.OAuth2)); err != nil {
			return fmt.Errorf("setting `oauth2_authorization`: %+v", err)
		}

		if err := d.Set("openid_authentication", flattenApiManagementOpenIDAuthentication(props.AuthenticationSettings.Openid)); err != nil {
			return fmt.Errorf("setting `openid_authentication`: %+v", err)
		}
	}

	return nil
}

func resourceApiManagementApiDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	apiid := id.Name

	name := apiid
	revision := ""
	if strings.Contains(apiid, ";") {
		name = strings.Split(apiid, ";")[0]
		revision = strings.Split(apiid, "=")[1]
	}

	deleteRevisions := utils.Bool(true)
	if resp, err := client.Delete(ctx, resourceGroup, serviceName, name, "", deleteRevisions); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting API %q / Revision %q (API Management Service %q / Resource Group %q): %s", name, revision, serviceName, resourceGroup, err)
		}
	}

	return nil
}

func expandApiManagementApiProtocols(input []interface{}) *[]apimanagement.Protocol {
	results := make([]apimanagement.Protocol, 0)

	for _, v := range input {
		results = append(results, apimanagement.Protocol(v.(string)))
	}

	return &results
}

func flattenApiManagementApiProtocols(input *[]apimanagement.Protocol) []string {
	if input == nil {
		return []string{}
	}

	results := make([]string, 0)
	for _, v := range *input {
		results = append(results, string(v))
	}

	return results
}

func expandApiManagementApiSubscriptionKeyParamNames(input []interface{}) *apimanagement.SubscriptionKeyParameterNamesContract {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	query := v["query"].(string)
	header := v["header"].(string)
	contract := apimanagement.SubscriptionKeyParameterNamesContract{
		Query:  utils.String(query),
		Header: utils.String(header),
	}
	return &contract
}

func flattenApiManagementApiSubscriptionKeyParamNames(paramNames *apimanagement.SubscriptionKeyParameterNamesContract) []interface{} {
	if paramNames == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if paramNames.Header != nil {
		result["header"] = *paramNames.Header
	}

	if paramNames.Query != nil {
		result["query"] = *paramNames.Query
	}

	return []interface{}{result}
}

func expandApiManagementOAuth2AuthenticationSettingsContract(input []interface{}) *apimanagement.OAuth2AuthenticationSettingsContract {
	if len(input) == 0 {
		return nil
	}

	oAuth2AuthorizationV := input[0].(map[string]interface{})
	return &apimanagement.OAuth2AuthenticationSettingsContract{
		AuthorizationServerID: utils.String(oAuth2AuthorizationV["authorization_server_name"].(string)),
		Scope:                 utils.String(oAuth2AuthorizationV["scope"].(string)),
	}
}

func flattenApiManagementOAuth2Authorization(input *apimanagement.OAuth2AuthenticationSettingsContract) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	authServerId := ""
	if input.AuthorizationServerID != nil {
		authServerId = *input.AuthorizationServerID
	}
	result["authorization_server_name"] = authServerId
	if input.Scope != nil {
		result["scope"] = *input.Scope
	}

	return []interface{}{result}
}

func expandApiManagementOpenIDAuthenticationSettingsContract(input []interface{}) *apimanagement.OpenIDAuthenticationSettingsContract {
	if len(input) == 0 {
		return nil
	}

	openIDAuthorizationV := input[0].(map[string]interface{})
	return &apimanagement.OpenIDAuthenticationSettingsContract{
		OpenidProviderID:          utils.String(openIDAuthorizationV["openid_provider_name"].(string)),
		BearerTokenSendingMethods: expandApiManagementOpenIDAuthenticationSettingsBearerTokenSendingMethods(openIDAuthorizationV["bearer_token_sending_methods"].(*pluginsdk.Set).List()),
	}
}

func expandApiManagementOpenIDAuthenticationSettingsBearerTokenSendingMethods(input []interface{}) *[]apimanagement.BearerTokenSendingMethods {
	if input == nil {
		return nil
	}
	results := make([]apimanagement.BearerTokenSendingMethods, 0)

	for _, v := range input {
		results = append(results, apimanagement.BearerTokenSendingMethods(v.(string)))
	}

	return &results
}

func flattenApiManagementOpenIDAuthentication(input *apimanagement.OpenIDAuthenticationSettingsContract) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	openIdProviderId := ""
	if input.OpenidProviderID != nil {
		openIdProviderId = *input.OpenidProviderID
	}
	result["openid_provider_name"] = openIdProviderId

	bearerTokenSendingMethods := make([]interface{}, 0)
	if s := input.BearerTokenSendingMethods; s != nil {
		for _, v := range *s {
			bearerTokenSendingMethods = append(bearerTokenSendingMethods, string(v))
		}
	}
	result["bearer_token_sending_methods"] = pluginsdk.NewSet(pluginsdk.HashString, bearerTokenSendingMethods)

	return []interface{}{result}
}
