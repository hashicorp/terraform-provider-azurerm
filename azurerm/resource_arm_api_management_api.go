package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/apimanagement/mgmt/2018-06-01-preview/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiManagementApiName,
			},

			"service_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiManagementServiceName,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ApiManagementApiPath,
			},

			"api_id": {
				Type:     schema.TypeString,
				Computed: true,
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

			"subscription_key_parameter_names": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"query": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"soap_pass_through": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
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
	d.Set("api_id", apiId)

	var properties *apimanagement.APICreateOrUpdateProperties

	_, hasImport := d.GetOk("import")

	// If import is used, we need to send properties to Azure API in two operations.
	// First we execute import and then updated the other props.
	if hasImport {
		properties = expandApiManagementImportProperties(d)
	} else {
		properties = expandApiManagementApiProperties(d)
	}

	apiParams := apimanagement.APICreateOrUpdateParameter{
		APICreateOrUpdateProperties: properties,
	}

	log.Printf("[DEBUG] Calling api with resource group %q, service name %q, api id %q", resGroup, serviceName, apiId)
	_, err := client.CreateOrUpdate(ctx, resGroup, serviceName, apiId, apiParams, "")
	if err != nil {
		return fmt.Errorf("Error creating/updating API Management API %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if hasImport {
		// Update with aditional properties not possible to send to Azure API during import
		updateParams := apimanagement.APICreateOrUpdateParameter{
			APICreateOrUpdateProperties: expandApiManagementApiProperties(d),
		}

		_, err := client.CreateOrUpdate(ctx, resGroup, serviceName, apiId, updateParams, "")

		if err != nil {
			return fmt.Errorf("Failed to update after import: %+v", err)
		}
	}

	read, err := client.Get(ctx, resGroup, serviceName, apiId)
	if err != nil {
		return fmt.Errorf("Error retrieving API Management API %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for API Management API %q (Resource Group %q)", name, resGroup)
	}

	d.SetId(*read.ID)

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

	d.Set("api_id", apiid)
	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("service_name", serviceName)

	if props := resp.APIContractProperties; props != nil {
		d.Set("display_name", props.DisplayName)
		d.Set("service_url", props.ServiceURL)
		d.Set("path", props.Path)
		d.Set("description", props.Description)
		d.Set("revision", props.APIRevision)
		d.Set("version", props.APIVersion)
		d.Set("version_set_id", props.APIVersionSetID)
		d.Set("is_current", props.IsCurrent)
		d.Set("is_online", props.IsOnline)
		d.Set("protocols", props.Protocols)
		d.Set("soap_pass_through", string(props.APIType) == string(apimanagement.SoapPassThrough))

		if err := d.Set("subscription_key_parameter_names", flattenApiManagementApiSubscriptionKeyParamNames(props.SubscriptionKeyParameterNames)); err != nil {
			return fmt.Errorf("Error setting `subscription_key_parameter_names`: %+v", err)
		}
	}

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
	soapPassThrough := d.Get("soap_pass_through").(bool)

	props := &apimanagement.APICreateOrUpdateProperties{
		Description:                   &description,
		DisplayName:                   &displayName,
		Path:                          &path,
		Protocols:                     expandApiManagementApiProtocols(d),
		ServiceURL:                    &serviceUrl,
		SubscriptionKeyParameterNames: expandApiManagementApiSubscriptionKeyParamNames(d),
	}

	if soapPassThrough {
		props.APIType = apimanagement.APIType(apimanagement.SoapPassThrough)
	}

	return props
}

func expandApiManagementApiSubscriptionKeyParamNames(d *schema.ResourceData) *apimanagement.SubscriptionKeyParameterNamesContract {
	var contract apimanagement.SubscriptionKeyParameterNamesContract

	if v, ok := d.GetOk("subscription_key_parameter_names.0.header"); ok {
		header := v.(string)
		contract.Header = &header
	}

	if v, ok := d.GetOk("subscription_key_parameter_names.0.query"); ok {
		query := v.(string)
		contract.Query = &query
	}

	if contract.Header == nil && contract.Query == nil {
		return nil
	}

	return &contract
}

func expandApiManagementApiProtocols(d *schema.ResourceData) *[]apimanagement.Protocol {
	protos := make([]apimanagement.Protocol, 0)

	if p, ok := d.GetOk("protocols"); ok {
		protocolsConfig := p.([]interface{})
		for _, v := range protocolsConfig {
			protos = append(protos, apimanagement.Protocol(v.(string)))
		}
	} else {
		// If not specified, set default to https
		protos = append(protos, apimanagement.ProtocolHTTPS)
	}

	return &protos
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

	if _, selectorUsed := d.GetOk("import.0.wsdl_selector"); selectorUsed {
		props.WsdlSelector = &apimanagement.APICreateOrUpdatePropertiesWsdlSelector{}

		if v, ok := d.GetOk("import.0.wsdl_selector.0.service_name"); ok {
			serviceName := v.(string)
			props.WsdlSelector.WsdlServiceName = &serviceName
		}

		if v, ok := d.GetOk("import.0.wsdl_selector.0.endpoint_name"); ok {
			endpointName := v.(string)
			props.WsdlSelector.WsdlEndpointName = &endpointName
		}
	}

	return props
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
