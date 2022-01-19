package apimanagement

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementApi() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiCreateUpdate,
		Read:   resourceApiManagementApiRead,
		Update: resourceApiManagementApiCreateUpdate,
		Delete: resourceApiManagementApiDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApiID(id)
			return err
		}),

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
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"path": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ApiManagementApiPath,
			},

			"protocols": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
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

			"revision_description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
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
								string(apimanagement.ContentFormatOpenapi),
								string(apimanagement.ContentFormatOpenapijson),
								string(apimanagement.ContentFormatOpenapijsonLink),
								string(apimanagement.ContentFormatOpenapiLink),
								string(apimanagement.ContentFormatSwaggerJSON),
								string(apimanagement.ContentFormatSwaggerLinkJSON),
								string(apimanagement.ContentFormatWadlLinkJSON),
								string(apimanagement.ContentFormatWadlXML),
								string(apimanagement.ContentFormatWsdl),
								string(apimanagement.ContentFormatWsdlLink),
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

			"source_api_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.ApiID,
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

			"version_description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewApiID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	revision := d.Get("revision").(string)
	path := d.Get("path").(string)
	apiId := fmt.Sprintf("%s;rev=%s", id.Name, revision)
	version := d.Get("version").(string)
	versionSetId := d.Get("version_set_id").(string)
	displayName := d.Get("display_name").(string)
	protocolsRaw := d.Get("protocols").(*pluginsdk.Set).List()
	protocols := expandApiManagementApiProtocols(protocolsRaw)
	sourceApiId := d.Get("source_api_id").(string)

	if version != "" && versionSetId == "" {
		return fmt.Errorf("setting `version` without the required `version_set_id`")
	}

	if sourceApiId == "" && (displayName == "" || protocols == nil || len(*protocols) == 0) {
		return fmt.Errorf("`display_name`, `protocols` are required when `source_api_id` is not set")
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, apiId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_api", id.ID())
		}
	}

	var apiType apimanagement.APIType
	var soapApiType apimanagement.SoapAPIType

	soapPassThrough := d.Get("soap_pass_through").(bool)
	if soapPassThrough {
		apiType = apimanagement.APITypeSoap
		soapApiType = apimanagement.SoapAPITypeSoapPassThrough
	} else {
		apiType = apimanagement.APITypeHTTP
		soapApiType = apimanagement.SoapAPITypeSoapToRest
	}

	// If import is used, we need to send properties to Azure API in two operations.
	// First we execute import and then updated the other props.
	if vs, hasImport := d.GetOk("import"); hasImport {
		importVs := vs.([]interface{})
		importV := importVs[0].(map[string]interface{})
		contentFormat := importV["content_format"].(string)
		contentValue := importV["content_value"].(string)

		log.Printf("[DEBUG] Importing API Management API %q of type %q", id.Name, contentFormat)
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
		if len(wsdlSelectorVs) == 0 && contentFormat == string(apimanagement.ContentFormatWsdl) {
			return fmt.Errorf("`wsdl_selector` is required when content format is `wsdl` in API Management API %q", id.Name)
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

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, apiId, apiParams, "")
		if err != nil {
			return fmt.Errorf("creating/updating %s: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on creating/updating %s: %+v", id, err)
		}
	}

	description := d.Get("description").(string)
	serviceUrl := d.Get("service_url").(string)
	subscriptionRequired := d.Get("subscription_required").(bool)

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
			Path:                          utils.String(path),
			Protocols:                     protocols,
			ServiceURL:                    utils.String(serviceUrl),
			SubscriptionKeyParameterNames: subscriptionKeyParameterNames,
			APIVersion:                    utils.String(version),
			SubscriptionRequired:          &subscriptionRequired,
			AuthenticationSettings:        authenticationSettings,
			APIRevisionDescription:        utils.String(d.Get("revision_description").(string)),
			APIVersionDescription:         utils.String(d.Get("version_description").(string)),
		},
	}

	if sourceApiId != "" {
		params.APICreateOrUpdateProperties.SourceAPIID = &sourceApiId
	}

	if displayName != "" {
		params.APICreateOrUpdateProperties.DisplayName = &displayName
	}

	if versionSetId != "" {
		params.APICreateOrUpdateProperties.APIVersionSetID = utils.String(versionSetId)
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, apiId, params, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
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

	name := id.Name
	revision := ""
	if strings.Contains(id.Name, ";") {
		name = strings.Split(id.Name, ";")[0]
		revision = strings.Split(id.Name, "=")[1]
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] API %q Revision %q (API Management Service %q / Resource Group %q) does not exist - removing from state!", name, revision, id.ServiceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving API %q / Revision %q (API Management Service %q / Resource Group %q): %+v", name, revision, id.ServiceName, id.ResourceGroup, err)
	}

	d.Set("api_management_name", id.ServiceName)
	d.Set("name", name)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.APIContractProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("is_current", props.IsCurrent)
		d.Set("is_online", props.IsOnline)
		d.Set("path", props.Path)
		d.Set("service_url", props.ServiceURL)
		d.Set("revision", props.APIRevision)
		d.Set("soap_pass_through", string(props.APIType) == string(apimanagement.SoapAPITypeSoapPassThrough))
		d.Set("subscription_required", props.SubscriptionRequired)
		d.Set("version", props.APIVersion)
		d.Set("version_set_id", props.APIVersionSetID)
		d.Set("revision_description", props.APIRevisionDescription)
		d.Set("version_description", props.APIVersionDescription)

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

	name := id.Name
	revision := ""
	if strings.Contains(id.Name, ";") {
		name = strings.Split(id.Name, ";")[0]
		revision = strings.Split(id.Name, "=")[1]
	}

	deleteRevisions := utils.Bool(true)
	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, name, "", deleteRevisions); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting API %q / Revision %q (API Management Service %q / Resource Group %q): %s", name, revision, id.ServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandApiManagementApiProtocols(input []interface{}) *[]apimanagement.Protocol {
	if len(input) == 0 {
		return nil
	}
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
