// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApi() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceApiManagementApiCreateUpdate,
		Read:   resourceApiManagementApiRead,
		Update: resourceApiManagementApiCreateUpdate,
		Delete: resourceApiManagementApiDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := api.ParseApiID(id)
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

			"resource_group_name": commonschema.ResourceGroupName(),

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
						string(api.ProtocolHTTP),
						string(api.ProtocolHTTPS),
						string(api.ProtocolWs),
						string(api.ProtocolWss),
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
			"api_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(api.ApiTypeGraphql),
					string(api.ApiTypeHTTP),
					string(api.ApiTypeSoap),
					string(api.ApiTypeWebsocket),
				}, false),
			},

			"contact": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"email": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.EmailAddress,
						},
						"name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"url": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},
					},
				},
			},

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
								string(api.ContentFormatOpenapi),
								string(api.ContentFormatOpenapiPositivejson),
								string(api.ContentFormatOpenapiPositivejsonNegativelink),
								string(api.ContentFormatOpenapiNegativelink),
								string(api.ContentFormatSwaggerNegativejson),
								string(api.ContentFormatSwaggerNegativelinkNegativejson),
								string(api.ContentFormatWadlNegativelinkNegativejson),
								string(api.ContentFormatWadlNegativexml),
								string(api.ContentFormatWsdl),
								string(api.ContentFormatWsdlNegativelink),
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

			"license": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"url": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
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

			"terms_of_service_url": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
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
									string(api.BearerTokenSendingMethodsAuthorizationHeader),
									string(api.BearerTokenSendingMethodsQuery),
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

	if !features.FourPointOhBeta() {
		resource.Schema["api_type"].ConflictsWith = []string{"soap_pass_through"}

		resource.Schema["soap_pass_through"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			Deprecated:    "`soap_pass_through` will be removed in favour of the property `api_type` in version 4.0 of the AzureRM Provider",
			ConflictsWith: []string{"api_type"},
		}
	}

	return resource
}

func resourceApiManagementApiCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := api.NewApiID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	revision := d.Get("revision").(string)
	path := d.Get("path").(string)
	apiId := fmt.Sprintf("%s;rev=%s", d.Get("name").(string), revision)
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

	newId := api.NewApiID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), apiId)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, newId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of an existing %s: %+v", newId, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_api", newId.ID())
		}
	}

	apiType := api.ApiTypeHTTP
	if v, ok := d.GetOk("api_type"); ok {
		apiType = api.ApiType(v.(string))
	}
	if !features.FourPointOhBeta() {
		if d.Get("soap_pass_through").(bool) {
			apiType = api.ApiTypeSoap
		}
	}

	soapApiType := map[api.ApiType]api.SoapApiType{
		api.ApiTypeGraphql:   api.SoapApiTypeGraphql,
		api.ApiTypeHTTP:      api.SoapApiTypeHTTP,
		api.ApiTypeSoap:      api.SoapApiTypeSoap,
		api.ApiTypeWebsocket: api.SoapApiTypeWebsocket,
	}[apiType]

	// If import is used, we need to send properties to Azure API in two operations.
	// First we execute import and then updated the other props.
	if vs, hasImport := d.GetOk("import"); hasImport {
		importVs := vs.([]interface{})
		importV := importVs[0].(map[string]interface{})
		contentFormat := importV["content_format"].(string)
		contentValue := importV["content_value"].(string)

		log.Printf("[DEBUG] Importing API Management API %q of type %q", id.ApiId, contentFormat)
		apiParams := api.ApiCreateOrUpdateParameter{
			Properties: &api.ApiCreateOrUpdateProperties{
				Type:       pointer.To(apiType),
				ApiType:    pointer.To(soapApiType),
				Format:     pointer.To(api.ContentFormat(contentFormat)),
				Value:      pointer.To(contentValue),
				Path:       path,
				ApiVersion: pointer.To(version),
			},
		}
		wsdlSelectorVs := importV["wsdl_selector"].([]interface{})

		if len(wsdlSelectorVs) > 0 {
			wsdlSelectorV := wsdlSelectorVs[0].(map[string]interface{})
			wSvcName := wsdlSelectorV["service_name"].(string)
			wEndpName := wsdlSelectorV["endpoint_name"].(string)

			apiParams.Properties.WsdlSelector = &api.ApiCreateOrUpdatePropertiesWsdlSelector{
				WsdlServiceName:  pointer.To(wSvcName),
				WsdlEndpointName: pointer.To(wEndpName),
			}
		}

		if versionSetId != "" {
			apiParams.Properties.ApiVersionSetId = pointer.To(versionSetId)
		}
		if err := client.CreateOrUpdateThenPoll(ctx, id, apiParams, api.CreateOrUpdateOperationOptions{}); err != nil {
			return fmt.Errorf("creating/updating %s: %+v", id, err)
		}

	}

	description := d.Get("description").(string)
	serviceUrl := d.Get("service_url").(string)
	subscriptionRequired := d.Get("subscription_required").(bool)

	subscriptionKeyParameterNamesRaw := d.Get("subscription_key_parameter_names").([]interface{})
	subscriptionKeyParameterNames := expandApiManagementApiSubscriptionKeyParamNames(subscriptionKeyParameterNamesRaw)

	authenticationSettings := &api.AuthenticationSettingsContract{}

	oAuth2AuthorizationSettingsRaw := d.Get("oauth2_authorization").([]interface{})
	oAuth2AuthorizationSettings := expandApiManagementOAuth2AuthenticationSettingsContract(oAuth2AuthorizationSettingsRaw)
	authenticationSettings.OAuth2 = oAuth2AuthorizationSettings

	openIDAuthorizationSettingsRaw := d.Get("openid_authentication").([]interface{})
	openIDAuthorizationSettings := expandApiManagementOpenIDAuthenticationSettingsContract(openIDAuthorizationSettingsRaw)
	authenticationSettings.Openid = openIDAuthorizationSettings

	contactInfoRaw := d.Get("contact").([]interface{})
	contactInfo := expandApiManagementApiContact(contactInfoRaw)

	licenseInfoRaw := d.Get("license").([]interface{})
	licenseInfo := expandApiManagementApiLicense(licenseInfoRaw)

	params := api.ApiCreateOrUpdateParameter{
		Properties: &api.ApiCreateOrUpdateProperties{
			Type:                          pointer.To(apiType),
			ApiType:                       pointer.To(soapApiType),
			Description:                   pointer.To(description),
			Path:                          path,
			Protocols:                     protocols,
			ServiceUrl:                    pointer.To(serviceUrl),
			SubscriptionKeyParameterNames: subscriptionKeyParameterNames,
			ApiVersion:                    pointer.To(version),
			SubscriptionRequired:          &subscriptionRequired,
			AuthenticationSettings:        authenticationSettings,
			ApiRevisionDescription:        pointer.To(d.Get("revision_description").(string)),
			ApiVersionDescription:         pointer.To(d.Get("version_description").(string)),
			Contact:                       contactInfo,
			License:                       licenseInfo,
		},
	}

	if sourceApiId != "" {
		params.Properties.SourceApiId = pointer.To(sourceApiId)
	}
	if displayName != "" {
		params.Properties.DisplayName = pointer.To(displayName)
	}
	if versionSetId != "" {
		params.Properties.ApiVersionSetId = pointer.To(versionSetId)
	}

	if v, ok := d.GetOk("terms_of_service_url"); ok {
		params.Properties.TermsOfServiceUrl = pointer.To(v.(string))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, newId, params, api.CreateOrUpdateOperationOptions{IfMatch: pointer.To("*")}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceApiManagementApiRead(d, meta)
}

func resourceApiManagementApiRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := api.ParseApiID(d.Id())
	if err != nil {
		return err
	}

	name := getApiName(id.ApiId)
	newId := api.NewApiID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, name)
	resp, err := client.Get(ctx, newId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", newId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", newId, err)
	}

	d.Set("api_management_name", id.ServiceName)
	d.Set("name", name)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			apiType := string(pointer.From(props.Type))
			if len(apiType) == 0 {
				apiType = string(api.ApiTypeHTTP)
			}
			d.Set("api_type", apiType)
			d.Set("description", pointer.From(props.Description))
			d.Set("display_name", pointer.From(props.DisplayName))
			d.Set("is_current", pointer.From(props.IsCurrent))
			d.Set("is_online", pointer.From(props.IsOnline))
			d.Set("path", props.Path)
			d.Set("service_url", pointer.From(props.ServiceUrl))
			d.Set("revision", pointer.From(props.ApiRevision))
			if !features.FourPointOhBeta() {
				d.Set("soap_pass_through", apiType == string(api.ApiTypeSoap))
			}
			d.Set("subscription_required", pointer.From(props.SubscriptionRequired))
			d.Set("version", pointer.From(props.ApiVersion))
			d.Set("version_set_id", pointer.From(props.ApiVersionSetId))
			d.Set("revision_description", pointer.From(props.ApiRevisionDescription))
			d.Set("version_description", pointer.From(props.ApiVersionDescription))
			d.Set("terms_of_service_url", pointer.From(props.TermsOfServiceUrl))

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

			if err := d.Set("contact", flattenApiManagementApiContact(props.Contact)); err != nil {
				return fmt.Errorf("setting `contact`: %+v", err)
			}

			if err := d.Set("license", flattenApiManagementApiLicense(props.License)); err != nil {
				return fmt.Errorf("setting `license`: %+v", err)
			}
		}
	}
	return nil
}

func resourceApiManagementApiDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := api.ParseApiID(d.Id())
	if err != nil {
		return err
	}

	name := getApiName(id.ApiId)

	newId := api.NewApiID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, name)
	if resp, err := client.Delete(ctx, newId, api.DeleteOperationOptions{DeleteRevisions: pointer.To(true)}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", newId, err)
		}
	}

	return nil
}

func expandApiManagementApiProtocols(input []interface{}) *[]api.Protocol {
	if len(input) == 0 {
		return nil
	}
	results := make([]api.Protocol, 0)

	for _, v := range input {
		results = append(results, api.Protocol(v.(string)))
	}

	return &results
}

func flattenApiManagementApiProtocols(input *[]api.Protocol) []string {
	if input == nil {
		return []string{}
	}

	results := make([]string, 0)
	for _, v := range *input {
		results = append(results, string(v))
	}

	return results
}

func expandApiManagementApiSubscriptionKeyParamNames(input []interface{}) *api.SubscriptionKeyParameterNamesContract {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	query := v["query"].(string)
	header := v["header"].(string)
	contract := api.SubscriptionKeyParameterNamesContract{
		Query:  pointer.To(query),
		Header: pointer.To(header),
	}
	return &contract
}

func flattenApiManagementApiSubscriptionKeyParamNames(paramNames *api.SubscriptionKeyParameterNamesContract) []interface{} {
	if paramNames == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["header"] = pointer.From(paramNames.Header)
	result["query"] = pointer.From(paramNames.Query)

	return []interface{}{result}
}

func expandApiManagementOAuth2AuthenticationSettingsContract(input []interface{}) *api.OAuth2AuthenticationSettingsContract {
	if len(input) == 0 {
		return nil
	}

	oAuth2AuthorizationV := input[0].(map[string]interface{})
	return &api.OAuth2AuthenticationSettingsContract{
		AuthorizationServerId: pointer.To(oAuth2AuthorizationV["authorization_server_name"].(string)),
		Scope:                 pointer.To(oAuth2AuthorizationV["scope"].(string)),
	}
}

func flattenApiManagementOAuth2Authorization(input *api.OAuth2AuthenticationSettingsContract) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["authorization_server_name"] = pointer.From(input.AuthorizationServerId)
	result["scope"] = pointer.From(input.Scope)

	return []interface{}{result}
}

func expandApiManagementOpenIDAuthenticationSettingsContract(input []interface{}) *api.OpenIdAuthenticationSettingsContract {
	if len(input) == 0 {
		return nil
	}

	openIDAuthorizationV := input[0].(map[string]interface{})
	return &api.OpenIdAuthenticationSettingsContract{
		OpenidProviderId:          pointer.To(openIDAuthorizationV["openid_provider_name"].(string)),
		BearerTokenSendingMethods: expandApiManagementOpenIDAuthenticationSettingsBearerTokenSendingMethods(openIDAuthorizationV["bearer_token_sending_methods"].(*pluginsdk.Set).List()),
	}
}

func expandApiManagementOpenIDAuthenticationSettingsBearerTokenSendingMethods(input []interface{}) *[]api.BearerTokenSendingMethods {
	if input == nil {
		return nil
	}
	results := make([]api.BearerTokenSendingMethods, 0)

	for _, v := range input {
		results = append(results, api.BearerTokenSendingMethods(v.(string)))
	}

	return &results
}

func flattenApiManagementOpenIDAuthentication(input *api.OpenIdAuthenticationSettingsContract) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["openid_provider_name"] = pointer.From(input.OpenidProviderId)

	bearerTokenSendingMethods := make([]interface{}, 0)
	if s := input.BearerTokenSendingMethods; s != nil {
		for _, v := range *s {
			bearerTokenSendingMethods = append(bearerTokenSendingMethods, string(v))
		}
	}
	result["bearer_token_sending_methods"] = pluginsdk.NewSet(pluginsdk.HashString, bearerTokenSendingMethods)

	return []interface{}{result}
}

func expandApiManagementApiContact(input []interface{}) *api.ApiContactInformation {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &api.ApiContactInformation{
		Email: pointer.To(v["email"].(string)),
		Name:  pointer.To(v["name"].(string)),
		Url:   pointer.To(v["url"].(string)),
	}
}

func flattenApiManagementApiContact(contact *api.ApiContactInformation) []interface{} {
	if contact == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["email"] = pointer.From(contact.Email)
	result["name"] = pointer.From(contact.Name)
	result["url"] = pointer.From(contact.Url)

	return []interface{}{result}
}

func expandApiManagementApiLicense(input []interface{}) *api.ApiLicenseInformation {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &api.ApiLicenseInformation{
		Name: pointer.To(v["name"].(string)),
		Url:  pointer.To(v["url"].(string)),
	}
}

func flattenApiManagementApiLicense(license *api.ApiLicenseInformation) []interface{} {
	if license == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["name"] = pointer.From(license.Name)
	result["url"] = pointer.From(license.Url)

	return []interface{}{result}
}

func getApiName(apiId string) string {
	name := apiId
	if strings.Contains(apiId, ";") {
		name = strings.Split(apiId, ";")[0]
	}

	return name
}
