package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
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

			"display_name": {
				Type:     schema.TypeString,
				Required: true,
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

			"protocols": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(apimanagement.ProtocolHTTP),
						string(apimanagement.ProtocolHTTPS),
					}, false),
				},
				Required: true,
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
							}, false),
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
				Optional: true,
				Computed: true,
				Default:  nil,
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
	revision := int32(d.Get("revision").(int))

	apiId := fmt.Sprintf("%s;rev=%d", name, revision)
	d.Set("api_id", apiId)

	var properties *apimanagement.APICreateOrUpdateProperties

	_, hasImport := d.GetOk("import")

	// If import is used, we need to send properties to Azure API in two operations.
	// First we execute import and then updated the other props.
	if hasImport {
		properties = expandApiManagementImportProperties(d)
		log.Printf("[DEBUG] Importing API Management API %q of type %q", name, properties.ContentFormat)

		apiParams := apimanagement.APICreateOrUpdateParameter{
			APICreateOrUpdateProperties: properties,
		}

		log.Printf("[DEBUG] Calling api with resource group %q, service name %q, api id %q", resGroup, serviceName, apiId)
		_, err := client.CreateOrUpdate(ctx, resGroup, serviceName, apiId, apiParams, "")
		if err != nil {
			return fmt.Errorf("Error creating/updating API Management API %q (Resource Group %q): %+v", name, resGroup, err)
		}
	}

	updateParams := apimanagement.APICreateOrUpdateParameter{
		APICreateOrUpdateProperties: expandApiManagementApiProperties(d),
	}

	log.Printf("[DEBUG] Calling api with resource group %q, service name %q, api id %q", resGroup, serviceName, apiId)
	_, err := client.CreateOrUpdate(ctx, resGroup, serviceName, apiId, updateParams, "")

	if err != nil {
		return fmt.Errorf("Error creating/updating API Management API %q (Resource Group %q): %+v", name, resGroup, err)
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
			log.Printf("[DEBUG] API Management API %q (Service %q / Resource Group %q) does not exist - removing from state!", name, serviceName, resGroup)
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
	displayName := d.Get("display_name").(string)
	path := d.Get("path").(string)
	serviceUrl := d.Get("service_url").(string)
	description := d.Get("description").(string)

	props := &apimanagement.APICreateOrUpdateProperties{
		Description:                   &description,
		DisplayName:                   &displayName,
		Path:                          &path,
		Protocols:                     expandApiManagementApiProtocols(d),
		ServiceURL:                    &serviceUrl,
		SubscriptionKeyParameterNames: expandApiManagementApiSubscriptionKeyParamNames(d),
	}

	if v, ok := d.GetOk("soap_pass_through"); ok {
		soapPassThrough := v.(bool)

		if soapPassThrough {
			props.APIType = apimanagement.APIType(apimanagement.SoapPassThrough)
		} else {
			props.APIType = apimanagement.APIType(apimanagement.SoapToRest)
		}
	}

	return props
}

func expandApiManagementApiSubscriptionKeyParamNames(d *schema.ResourceData) *apimanagement.SubscriptionKeyParameterNamesContract {
	vs := d.Get("subscription_key_parameter_names").([]interface{})

	if len(vs) > 0 {
		v := vs[0].(map[string]interface{})
		contract := apimanagement.SubscriptionKeyParameterNamesContract{}

		query := v["query"].(string)
		header := v["header"].(string)

		if query != "" {
			contract.Query = utils.String(query)
		}

		if header != "" {
			contract.Header = utils.String(header)
		}

		if query != "" || header != "" {
			return &contract
		}
	}

	return nil
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

	importVs := d.Get("import").([]interface{})
	importV := importVs[0].(map[string]interface{})

	props.ContentFormat = apimanagement.ContentFormat(importV["content_format"].(string))

	cVal := importV["content_value"].(string)
	props.ContentValue = &cVal

	wsdlSelectorVs := importV["wsdl_selector"].([]interface{})

	if len(wsdlSelectorVs) > 0 {
		wsdlSelectorV := wsdlSelectorVs[0].(map[string]interface{})
		props.WsdlSelector = &apimanagement.APICreateOrUpdatePropertiesWsdlSelector{}

		wSvcName := wsdlSelectorV["service_name"].(string)
		wEndpName := wsdlSelectorV["endpoint_name"].(string)

		props.WsdlSelector.WsdlServiceName = &wSvcName
		props.WsdlSelector.WsdlEndpointName = &wEndpName
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
