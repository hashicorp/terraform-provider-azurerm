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
			"name": { // DisplayName
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
				Computed: true, // Azure API sets protocols to https by default
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
				Computed: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"is_current": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_online": {
				Type:     schema.TypeBool,
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

	//Currently we don't support revisions, so we use 1 as default
	apiId := fmt.Sprintf("%s;rev=%d", name, 1)

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

	name := apiid
	if strings.Contains(apiid, ";") {
		name = strings.Split(apiid, ";")[0]
	}

	ctx := client.StopContext
	resp, err := apiManagementApiClient.Get(ctx, resGroup, serviceName, apiid)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on API Management API %q on service %q (Resource Group %q): %+v", apiid, serviceName, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("service_name", serviceName)

	if props := resp.APIContractProperties; props != nil {
		d.Set("path", props.Path)
		d.Set("description", props.Description)
		d.Set("soap_api_type", props.APIType)
		d.Set("revision", props.APIRevision)
		d.Set("version", props.APIVersion)
		d.Set("version_set_id", props.APIVersionSetID)
		d.Set("is_current", props.IsCurrent)
		d.Set("is_online", props.IsOnline)
		d.Set("protocols", props.Protocols)
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
	displayName := d.Get("name").(string)
	path := d.Get("path").(string)
	serviceUrl := d.Get("service_url").(string)
	description := d.Get("description").(string)
	soapApiTypeConfig := d.Get("soap_api_type").(string)
	soapApiType := apimanagement.APIType(soapApiTypeConfig)

	// versionSetId := d.Get("api_version_set_id").(string)

	// var oAuth *apimanagement.OAuth2AuthenticationSettingsContract
	// if oauthConfig := d.Get("oauth").([]interface{}); oauthConfig != nil && len(oauthConfig) > 0 {
	// 	oAuth = expandApiManagementApiOAuth(oauthConfig)
	// }

	return &apimanagement.APICreateOrUpdateProperties{
		APIType: soapApiType,
		// AuthenticationSettings: &apimanagement.AuthenticationSettingsContract{
		// 	OAuth2: oAuth,
		// },
		Description: &description,
		DisplayName: &displayName,
		Path:        &path,
		Protocols:   expandApiManagementApiProtocols(d),
		ServiceURL:  &serviceUrl,
		// SubscriptionKeyParameterNames: nil,
	}
}

func expandApiManagementApiProtocols(d *schema.ResourceData) *[]apimanagement.Protocol {
	protos := make([]apimanagement.Protocol, 0)

	if p, ok := d.GetOk("protocols"); ok {
		protocolsConfig := p.([]interface{})
		for _, v := range protocolsConfig {
			protos = append(protos, apimanagement.Protocol(v.(string)))
		}
	}

	return &protos
}

func expandApiManagementApiUpdateProperties(d *schema.ResourceData) *apimanagement.APIContractUpdateProperties {
	name := d.Get("name").(string)
	path := d.Get("path").(string)
	serviceUrl := d.Get("service_url").(string)
	description := d.Get("description").(string)
	soapApiTypeConfig := d.Get("soap_api_type").(string)

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
		APIType: soapApiType,
		// AuthenticationSettings: nil,
		Description: &description,
		DisplayName: &name,
		Path:        &path,
		Protocols:   expandApiManagementApiProtocols(d),
		ServiceURL:  &serviceUrl,
		// SubscriptionKeyParameterNames: nil,
	}
}

func expandApiManagementImportProperties(d *schema.ResourceData) *apimanagement.APICreateOrUpdateProperties {
	path := d.Get("path").(string)

	props := &apimanagement.APICreateOrUpdateProperties{
		Path: &path,
	}

	if v, ok := d.GetOk("import.0.content_format"); ok {
		props.ContentFormat = apimanagement.ContentFormat(v.(string))
	}

	if v, ok := d.GetOk("import.0.content_value"); ok {
		content_val := v.(string)
		props.ContentValue = &content_val
	}

	return props
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
