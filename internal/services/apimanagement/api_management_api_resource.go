// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApi() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceApiManagementApiCreate,
		Read:   resourceApiManagementApiRead,
		Update: resourceApiManagementApiUpdate,
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

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiV0ToV1{},
		}),

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
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"openid_authentication"},
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
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"oauth2_authorization"},
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

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
				values := d.GetRawConfig().AsValueMap()
				if d.Get("version").(string) != "" && values["version_set_id"].IsNull() {
					return errors.New("setting `version` without the required `version_set_id`")
				}

				protocols := expandApiManagementApiProtocols(d.Get("protocols").(*pluginsdk.Set).List())
				if values["source_api_id"].IsNull() && (values["display_name"].IsNull() || protocols == nil || len(*protocols) == 0) {
					return errors.New("`display_name`, `protocols` are required when `source_api_id` is not set")
				}
				return nil
			}),
		),
	}

	return resource
}

func resourceApiManagementApiCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	revision := d.Get("revision").(string)
	path := d.Get("path").(string)
	apiId := fmt.Sprintf("%s;rev=%s", d.Get("name").(string), revision)
	version := d.Get("version").(string)
	versionSetId := d.Get("version_set_id").(string)
	displayName := d.Get("display_name").(string)
	protocolsRaw := d.Get("protocols").(*pluginsdk.Set).List()
	protocols := expandApiManagementApiProtocols(protocolsRaw)
	sourceApiId := d.Get("source_api_id").(string)

	id := api.NewApiID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), apiId)
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_api_management_api", id.ID())
	}

	apiType := api.ApiTypeHTTP
	if v, ok := d.GetOk("api_type"); ok {
		apiType = api.ApiType(v.(string))
	}
	soapApiType := soapApiTypeFromApiType(apiType)

	// If import is used, we need to send properties to Azure API in two operations.
	// First we execute import and then updated the other props.
	if importVs, ok := d.GetOk("import"); ok {
		if apiParams := expandApiManagementApiImport(importVs.([]interface{}), apiType, soapApiType,
			path, d.Get("service_url").(string), version, versionSetId); apiParams != nil {
			result, err := client.CreateOrUpdate(ctx, id, *apiParams, api.CreateOrUpdateOperationOptions{})
			if err != nil {
				return fmt.Errorf("creating with import of %s: %+v", id, err)
			}

			if pollerType := custompollers.NewAPIManagementAPIPoller(client, id, result.HttpResponse); pollerType != nil {
				poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("polling import %s: %+v", id, err)
				}
			}
		}
	}

	serviceUrl := d.Get("service_url").(string)
	subscriptionRequired := d.Get("subscription_required").(bool)

	subscriptionKeyParameterNames := expandApiManagementApiSubscriptionKeyParamNames(d.Get("subscription_key_parameter_names").([]interface{}))

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
			Path:                          path,
			Protocols:                     protocols,
			ServiceURL:                    pointer.To(serviceUrl),
			SubscriptionKeyParameterNames: subscriptionKeyParameterNames,
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

	if description, ok := d.GetOk("description"); ok {
		params.Properties.Description = pointer.To(description.(string))
	}

	if displayName != "" {
		params.Properties.DisplayName = pointer.To(displayName)
	}

	if version != "" {
		params.Properties.ApiVersion = pointer.To(version)
	}

	if versionSetId != "" {
		params.Properties.ApiVersionSetId = pointer.To(versionSetId)
	}

	if v, ok := d.GetOk("terms_of_service_url"); ok {
		params.Properties.TermsOfServiceURL = pointer.To(v.(string))
	}

	result, err := client.CreateOrUpdate(ctx, id, params, api.CreateOrUpdateOperationOptions{IfMatch: pointer.To("*")})
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if pollerType := custompollers.NewAPIManagementAPIPoller(client, id, result.HttpResponse); pollerType != nil {
		poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("polling creating/updating %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	return resourceApiManagementApiRead(d, meta)
}

func resourceApiManagementApiUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	path := d.Get("path").(string)
	version := d.Get("version").(string)
	versionSetId := d.Get("version_set_id").(string)
	displayName := d.Get("display_name").(string)
	protocolsRaw := d.Get("protocols").(*pluginsdk.Set).List()
	protocols := expandApiManagementApiProtocols(protocolsRaw)
	sourceApiId := d.Get("source_api_id").(string)
	serviceUrl := d.Get("service_url").(string)

	id, err := api.ParseApiID(d.Id())
	if err != nil {
		return err
	}

	apiType := api.ApiTypeHTTP
	if v, ok := d.GetOk("api_type"); ok {
		apiType = api.ApiType(v.(string))
	}
	soapApiType := soapApiTypeFromApiType(apiType)

	// If import is used, we need to send properties to Azure API in two operations.
	// First we execute import and then updated the other props.
	if d.HasChange("import") {
		if vs, hasImport := d.GetOk("import"); hasImport {
			d.Partial(true)
			if apiParams := expandApiManagementApiImport(vs.([]interface{}), apiType, soapApiType,
				path, serviceUrl, version, versionSetId); apiParams != nil {
				result, err := client.CreateOrUpdate(ctx, *id, *apiParams, api.CreateOrUpdateOperationOptions{})
				if err != nil {
					return fmt.Errorf("creating with import of %s: %+v", id, err)
				}

				if pollerType := custompollers.NewAPIManagementAPIPoller(client, *id, result.HttpResponse); pollerType != nil {
					poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
					if err := poller.PollUntilDone(ctx); err != nil {
						return fmt.Errorf("polling import %s: %+v", id, err)
					}
				}
			}
			d.Partial(false)
		}
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	existing := resp.Model.Properties
	if existing.Type != nil {
		soapApiType = soapApiTypeFromApiType(pointer.From(existing.Type))
	}
	prop := &api.ApiCreateOrUpdateProperties{
		Path:                          existing.Path,
		Protocols:                     existing.Protocols,
		ServiceURL:                    existing.ServiceURL,
		Description:                   existing.Description,
		ApiVersionDescription:         existing.ApiVersionDescription,
		ApiRevisionDescription:        existing.ApiRevisionDescription,
		SubscriptionRequired:          existing.SubscriptionRequired,
		SubscriptionKeyParameterNames: existing.SubscriptionKeyParameterNames,
		Contact:                       existing.Contact,
		License:                       existing.License,
		SourceApiId:                   existing.SourceApiId,
		DisplayName:                   existing.DisplayName,
		ApiVersion:                    existing.ApiVersion,
		ApiVersionSetId:               existing.ApiVersionSetId,
		TermsOfServiceURL:             existing.TermsOfServiceURL,
		Type:                          existing.Type,
		ApiType:                       pointer.To(soapApiType),
	}

	// For the setting of `AuthenticationSettingsContract`, the PUT payload restrictions are as follows:
	//   1. Cannot have both 'oAuth2' and 'openid' set
	//   2. Cannot use `OAuth2AuthenticationSettings` in combination with `OAuth2` nor `openid`
	//   3. Cannot use `OpenidAuthenticationSettings` in combination with `Openid` nor `OAuth2`
	// If specifying `oauth2_authorization`/`openid_authentication` when creating a resource and then updating the resource, the error #2/#3 mentioned above will occur.
	// This is because starting from the 2022-08-01 version, the Get API additionally returns a collection of `oauth2_authorization`/`openid_authentication` authentication settings, which property name is `OAuth2AuthenticationSettings`/`OpenidAuthenticationSetting`.
	// Given the API behavior, the update here should only read the specified property `oauth2_authorization`/`openid_authentication` to exclude `OAuth2AuthenticationSettings`/`OpenidAuthenticationSetting` to ensure the update works properly.
	if v := existing.AuthenticationSettings; v != nil {
		authenticationSettings := &api.AuthenticationSettingsContract{}
		if v.OAuth2 != nil {
			authenticationSettings.OAuth2 = v.OAuth2
			prop.AuthenticationSettings = authenticationSettings
		}

		if v.Openid != nil {
			authenticationSettings.Openid = v.Openid
			prop.AuthenticationSettings = authenticationSettings
		}
	}

	if d.HasChange("path") {
		prop.Path = path
	}

	if d.HasChange("protocols") {
		prop.Protocols = protocols
	}

	if d.HasChange("api_type") {
		prop.Type = pointer.To(apiType)
		prop.ApiType = pointer.To(soapApiType)
	}

	if d.HasChange("service_url") {
		prop.ServiceURL = pointer.To(serviceUrl)
	}

	if d.HasChange("description") {
		prop.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("revision_description") {
		prop.ApiRevisionDescription = pointer.To(d.Get("revision_description").(string))
	}

	if d.HasChange("version_description") {
		prop.ApiVersionDescription = pointer.To(d.Get("version_description").(string))
	}
	if d.HasChange("subscription_required") {
		prop.SubscriptionRequired = pointer.To(d.Get("subscription_required").(bool))
	}

	if d.HasChange("subscription_key_parameter_names") {
		subscriptionKeyParameterNamesRaw := d.Get("subscription_key_parameter_names").([]interface{})
		prop.SubscriptionKeyParameterNames = expandApiManagementApiSubscriptionKeyParamNames(subscriptionKeyParameterNamesRaw)
	}

	if d.HasChange("oauth2_authorization") {
		authenticationSettings := &api.AuthenticationSettingsContract{}
		oAuth2AuthorizationSettingsRaw := d.Get("oauth2_authorization").([]interface{})
		oAuth2AuthorizationSettings := expandApiManagementOAuth2AuthenticationSettingsContract(oAuth2AuthorizationSettingsRaw)
		authenticationSettings.OAuth2 = oAuth2AuthorizationSettings
		prop.AuthenticationSettings = authenticationSettings
	}

	if d.HasChange("openid_authentication") {
		authenticationSettings := &api.AuthenticationSettingsContract{}
		openIDAuthorizationSettingsRaw := d.Get("openid_authentication").([]interface{})
		openIDAuthorizationSettings := expandApiManagementOpenIDAuthenticationSettingsContract(openIDAuthorizationSettingsRaw)
		authenticationSettings.Openid = openIDAuthorizationSettings
		prop.AuthenticationSettings = authenticationSettings
	}

	if d.HasChange("contact") {
		prop.Contact = expandApiManagementApiContact(d.Get("contact").([]interface{}))
	}

	if d.HasChange("license") {
		prop.License = expandApiManagementApiLicense(d.Get("license").([]interface{}))
	}

	if d.HasChange("source_api_id") {
		prop.SourceApiId = pointer.To(sourceApiId)
	}

	if d.HasChange("display_name") {
		prop.DisplayName = pointer.To(displayName)
	}

	if d.HasChange("version") {
		prop.ApiVersion = pointer.To(version)
	}

	if d.HasChange("version_set_id") {
		prop.ApiVersionSetId = pointer.To(versionSetId)
	}

	if d.HasChange("terms_of_service_url") {
		prop.TermsOfServiceURL = pointer.To(d.Get("terms_of_service_url").(string))
	}

	params := api.ApiCreateOrUpdateParameter{
		Properties: prop,
	}

	result, err := client.CreateOrUpdate(ctx, *id, params, api.CreateOrUpdateOperationOptions{IfMatch: pointer.To("*")})
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if pollerType := custompollers.NewAPIManagementAPIPoller(client, *id, result.HttpResponse); pollerType != nil {
		poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("polling creating/updating %s: %+v", id, err)
		}
	}

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

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("api_management_name", id.ServiceName)
	d.Set("name", getApiName(id.ApiId))
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
			d.Set("service_url", pointer.From(props.ServiceURL))
			d.Set("revision", pointer.From(props.ApiRevision))
			d.Set("subscription_required", pointer.From(props.SubscriptionRequired))
			d.Set("version", pointer.From(props.ApiVersion))
			d.Set("version_set_id", pointer.From(props.ApiVersionSetId))
			d.Set("revision_description", pointer.From(props.ApiRevisionDescription))
			d.Set("version_description", pointer.From(props.ApiVersionDescription))
			d.Set("terms_of_service_url", pointer.From(props.TermsOfServiceURL))

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

	if resp, err := client.Delete(ctx, *id, api.DefaultDeleteOperationOptions()); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}

func soapApiTypeFromApiType(apiType api.ApiType) api.SoapApiType {
	return map[api.ApiType]api.SoapApiType{
		api.ApiTypeGraphql:   api.SoapApiTypeGraphql,
		api.ApiTypeHTTP:      api.SoapApiTypeHTTP,
		api.ApiTypeSoap:      api.SoapApiTypeSoap,
		api.ApiTypeWebsocket: api.SoapApiTypeWebsocket,
	}[apiType]
}

func expandApiManagementApiImport(importVs []interface{}, apiType api.ApiType, soapApiType api.SoapApiType, path, serviceUrl, version, versionSetId string) *api.ApiCreateOrUpdateParameter {
	if len(importVs) == 0 || importVs[0] == nil {
		return nil
	}

	importV := importVs[0].(map[string]interface{})
	if len(importV) == 0 {
		return nil
	}

	contentFormat := importV["content_format"].(string)
	contentValue := importV["content_value"].(string)

	apiParams := api.ApiCreateOrUpdateParameter{
		Properties: &api.ApiCreateOrUpdateProperties{
			Type:    pointer.To(apiType),
			ApiType: pointer.To(soapApiType),
			Format:  pointer.To(api.ContentFormat(contentFormat)),
			Value:   pointer.To(contentValue),
			Path:    path,
		},
	}

	wsdlSelectorVs := importV["wsdl_selector"].([]interface{})
	if len(wsdlSelectorVs) > 0 && wsdlSelectorVs[0] != nil {
		if wsdlSelectorV := wsdlSelectorVs[0].(map[string]interface{}); len(wsdlSelectorV) > 0 {
			wSvcName := wsdlSelectorV["service_name"].(string)
			wEndpName := wsdlSelectorV["endpoint_name"].(string)

			apiParams.Properties.WsdlSelector = &api.ApiCreateOrUpdatePropertiesWsdlSelector{
				WsdlServiceName:  pointer.To(wSvcName),
				WsdlEndpointName: pointer.To(wEndpName),
			}
		}
	}
	if serviceUrl != "" {
		apiParams.Properties.ServiceURL = pointer.To(serviceUrl)
	}

	if version != "" {
		apiParams.Properties.ApiVersion = pointer.To(version)
	}

	if versionSetId != "" {
		apiParams.Properties.ApiVersionSetId = pointer.To(versionSetId)
	}

	return &apiParams
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
