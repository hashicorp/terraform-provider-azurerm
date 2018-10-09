package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/apimanagement/mgmt/2018-06-01-preview/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementApiCreateUpdate,
		Read:   resourceArmApiManagementApiRead,
		Update: resourceArmApiManagementApiCreateUpdate,
		Delete: resourceArmApiManagementApiDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// ValidateFunc: azure.ValidateApiManagementApiName,
			},

			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"path": {
				Type:     schema.TypeString,
				Required: true,
			},

			"service_url": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"import": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"content_format": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(apimanagement.SwaggerJSON),
								string(apimanagement.SwaggerLinkJSON),
								string(apimanagement.WadlLinkJSON),
								string(apimanagement.WadlXML),
								string(apimanagement.Wsdl),
								string(apimanagement.WsdlLink),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"wsdl_selector": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"endpoint_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"protocols": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
					ValidateFunc: validation.StringInSlice([]string{
						string(apimanagement.ProtocolHTTP),
						string(apimanagement.ProtocolHTTPS),
					}, true),
				},
				Optional: true,
			},

			"oauth": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authorization_server_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"scope": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"subscription_key": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"query": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"soap_api_type": {
				Type:    schema.TypeString,
				Default: "",
				ValidateFunc: validation.StringInSlice([]string{
					string(apimanagement.HTTP),
					string(apimanagement.Soap),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				Optional:         true,
			},

			"revision": {
				Type:     schema.TypeInt,
				Default:  1,
				Optional: true,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"version_set": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"scheme": {
							Type: schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(apimanagement.VersioningSchemeQuery),
								string(apimanagement.VersioningSchemeHeader),
								string(apimanagement.VersioningSchemeSegment),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							Required:         true,
						},
						"query_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"header_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"version_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_current": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"is_online": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArmApiManagementApiCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementApiClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM API Management API creation.")

	resGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("service_name").(string)
	name := d.Get("name").(string)
	revision := int32(d.Get("revision").(int))

	apiId := fmt.Sprintf("%s;rev=%d", name, revision)

	var properties *apimanagement.APICreateOrUpdateProperties
	var updateProperties *apimanagement.APIContractUpdateProperties

	_, isImport := d.GetOk("import")

	if isImport {
		properties = expandApiManagementImportProperties(d)
		updateProperties = expandApiManagementApiUpdateProperties(d)
	} else {
		properties = expandApiManagementApiProperties(d)
	}

	apiParams := apimanagement.APICreateOrUpdateParameter{
		APICreateOrUpdateProperties: properties,
	}

	log.Printf("[DEBUG] Calling api with resource group %q, service name %q, api id %q", resGroup, serviceName, apiId)
	log.Printf("[DEBUG] Listing api params:")
	log.Printf("%+v\n", apiParams.APICreateOrUpdateProperties)

	apiContract, err := client.CreateOrUpdate(ctx, resGroup, serviceName, apiId, apiParams, "")
	if err != nil {
		return err
	}

	if isImport {
		updateParams := apimanagement.APIUpdateContract{
			APIContractUpdateProperties: updateProperties,
		}

		_, err := client.Update(ctx, resGroup, serviceName, apiId, updateParams, "")

		if err != nil {
			return fmt.Errorf("Failed to update after import: %+v", err)
		}
	}

	d.SetId(*apiContract.ID)

	return resourceArmApiManagementApiRead(d, meta)
}

func resourceArmApiManagementApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	apiManagementApiClient := meta.(*ArmClient).apiManagementApiClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiid := id.Path["apis"]

	ctx := client.StopContext
	resp, err := apiManagementApiClient.Get(ctx, resGroup, serviceName, apiid)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on API Management API %q on service %q (Resource Group %q): %+v", apiid, serviceName, resGroup, err)
	}

	log.Printf("%+v\n", resp)
	return nil
}

func resourceArmApiManagementApiDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementApiClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiid := id.Path["apis"]

	log.Printf("[DEBUG] Deleting api management api %s: %s", resGroup, apiid)

	resp, err := client.Delete(ctx, resGroup, serviceName, apiid, "*", utils.Bool(true))

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return err
	}

	return nil
}

func expandApiManagementApiProperties(d *schema.ResourceData) *apimanagement.APICreateOrUpdateProperties {
	revision := d.Get("revision").(string)
	displayName := d.Get("name").(string)
	path := d.Get("path").(string)
	serviceUrl := d.Get("service_url").(string)
	description := d.Get("description").(string)
	soapApiTypeConfig := d.Get("soap_api_type").(string)
	version := d.Get("version").(string)
	versionSetId := d.Get("version_set_id").(string)
	isCurrent := d.Get("is_current").(bool)
	isOnline := d.Get("is_online").(bool)

	soapApiType := apimanagement.APIType(soapApiTypeConfig)

	protos := make([]apimanagement.Protocol, 0)
	if p, ok := d.GetOk("protocols.0"); ok {
		protocolsConfig := p.([]interface{})
		for _, v := range protocolsConfig {
			protos = append(protos, apimanagement.Protocol(v.(string)))
		}
	}
	if len(protos) == 0 {
		protos = append(protos, apimanagement.ProtocolHTTPS)
	}

	// versionSetId := d.Get("api_version_set_id").(string)

	// var oAuth *apimanagement.OAuth2AuthenticationSettingsContract
	// if oauthConfig := d.Get("oauth").([]interface{}); oauthConfig != nil && len(oauthConfig) > 0 {
	// 	oAuth = expandApiManagementApiOAuth(oauthConfig)
	// }

	return &apimanagement.APICreateOrUpdateProperties{
		APIRevision: &revision,
		APIType:     soapApiType,
		APIVersion:  &version,
		// APIVersionSet: nil,
		APIVersionSetID: &versionSetId,
		// AuthenticationSettings: &apimanagement.AuthenticationSettingsContract{
		// 	OAuth2: oAuth,
		// },
		Description: &description,
		DisplayName: &displayName,
		IsCurrent:   &isCurrent,
		IsOnline:    &isOnline,
		Path:        &path,
		Protocols:   &protos,
		ServiceURL:  &serviceUrl,
		// SubscriptionKeyParameterNames: nil,
	}
}

func expandApiManagementApiUpdateProperties(d *schema.ResourceData) *apimanagement.APIContractUpdateProperties {
	name := d.Get("name").(string)
	path := d.Get("path").(string)
	serviceUrl := d.Get("service_url").(string)
	description := d.Get("description").(string)
	soapApiTypeConfig := d.Get("soap_api_type").(string)
	// revisionDescription := d.Get("revision_description").(string)

	protos := make([]apimanagement.Protocol, 0)

	if p, ok := d.GetOk("protocols.0"); ok {
		protocolsConfig := p.([]interface{})
		for _, v := range protocolsConfig {
			protos = append(protos, apimanagement.Protocol(v.(string)))
		}
	}

	if len(protos) == 0 {
		protos = append(protos, apimanagement.ProtocolHTTPS)
	}

	var soapApiType apimanagement.APIType

	switch s := strings.ToLower(soapApiTypeConfig); s {
	case "http":
		soapApiType = apimanagement.HTTP
	case "soap":
		soapApiType = apimanagement.Soap
	}

	// versionSetId := d.Get("api_version_set_id").(string)

	// var oAuth *apimanagement.OAuth2AuthenticationSettingsContract
	// if oauthConfig := d.Get("oauth").([]interface{}); oauthConfig != nil && len(oauthConfig) > 0 {
	// 	oAuth = expandApiManagementApiOAuth(oauthConfig)
	// }

	log.Printf("ServiceURL: %s", &serviceUrl)

	return &apimanagement.APIContractUpdateProperties{
		// APIRevision: nil,
		APIType: soapApiType,
		// APIVersion: nil,
		// APIVersionSetID: nil,
		// AuthenticationSettings: nil,
		Description: &description,
		DisplayName: &name,
		// IsCurrent: nil,
		// IsOnline: nil,
		Path:       &path,
		Protocols:  &protos,
		ServiceURL: &serviceUrl,
		// SubscriptionKeyParameterNames: nil,
	}
}

func expandApiManagementImportProperties(d *schema.ResourceData) *apimanagement.APICreateOrUpdateProperties {
	path := d.Get("path").(string)

	var contentFormat apimanagement.ContentFormat
	if v, ok := d.GetOk("import.0.content_format"); ok {
		contentFormat = apimanagement.ContentFormat(v.(string))
	}

	var contentValue string
	if v, ok := d.GetOk("import.0.content_value"); ok {
		contentValue = v.(string)
	}

	return &apimanagement.APICreateOrUpdateProperties{
		Path:          &path,
		ContentFormat: contentFormat,
		ContentValue:  &contentValue,
	}
}

func expandApiManagementApiOAuth(oauth []interface{}) *apimanagement.OAuth2AuthenticationSettingsContract {
	config := oauth[0].(map[string]interface{})

	authorization_server_id := config["authorization_server_id"].(string)
	scope := config["scope"].(string)

	return &apimanagement.OAuth2AuthenticationSettingsContract{
		AuthorizationServerID: &authorization_server_id,
		Scope: &scope,
	}
}

func flattenApiManagementApiContract(apiContract apimanagement.APIContract) error {
	return nil
}
