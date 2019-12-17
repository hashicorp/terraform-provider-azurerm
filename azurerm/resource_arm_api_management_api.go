package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiManagementApiName,
			},

			"api_management_name": azure.SchemaApiManagementName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ApiManagementApiPath,
			},

			"protocols": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(apimanagement.ProtocolHTTP),
						string(apimanagement.ProtocolHTTPS),
					}, false),
				},
			},

			"revision": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			// Optional
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
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
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},

									"endpoint_name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
								},
							},
						},
					},
				},
			},

			"service_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"subscription_key_parameter_names": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"query": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"soap_pass_through": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			// Computed
			"is_current": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_online": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"version_set_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceArmApiManagementApiCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error setting `version` without the required `version_set_id`")
	}

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apiId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing API %q (API Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
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
				APIType:       apiType,
				SoapAPIType:   soapApiType,
				ContentFormat: apimanagement.ContentFormat(contentFormat),
				ContentValue:  utils.String(contentValue),
				Path:          utils.String(path),
				APIVersion:    utils.String(version),
			},
		}
		wsdlSelectorVs := importV["wsdl_selector"].([]interface{})
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

		if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiId, apiParams, ""); err != nil {
			return fmt.Errorf("Error creating/updating API Management API %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	serviceUrl := d.Get("service_url").(string)

	protocolsRaw := d.Get("protocols").(*schema.Set).List()
	protocols := expandApiManagementApiProtocols(protocolsRaw)

	subscriptionKeyParameterNamesRaw := d.Get("subscription_key_parameter_names").([]interface{})
	subscriptionKeyParameterNames := expandApiManagementApiSubscriptionKeyParamNames(subscriptionKeyParameterNamesRaw)

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
		},
	}

	if versionSetId != "" {
		params.APICreateOrUpdateProperties.APIVersionSetID = utils.String(versionSetId)
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiId, params, ""); err != nil {
		return fmt.Errorf("Error creating/updating API %q / Revision %q (API Management Service %q / Resource Group %q): %+v", name, revision, serviceName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, serviceName, apiId)
	if err != nil {
		return fmt.Errorf("Error retrieving API %q / Revision %q (API Management Service %q / Resource Group %q): %+v", name, revision, serviceName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for API %q / Revision %q (API Management Service %q / Resource Group %q)", name, revision, serviceName, resourceGroup)
	}

	d.SetId(*read.ID)
	return resourceArmApiManagementApiRead(d, meta)
}

func resourceArmApiManagementApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiid := id.Path["apis"]

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

		return fmt.Errorf("Error retrieving API %q / Revision %q (API Management Service %q / Resource Group %q): %+v", name, revision, serviceName, resourceGroup, err)
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
		d.Set("version", props.APIVersion)
		d.Set("version_set_id", props.APIVersionSetID)

		if err := d.Set("protocols", flattenApiManagementApiProtocols(props.Protocols)); err != nil {
			return fmt.Errorf("Error setting `protocols`: %s", err)
		}

		if err := d.Set("subscription_key_parameter_names", flattenApiManagementApiSubscriptionKeyParamNames(props.SubscriptionKeyParameterNames)); err != nil {
			return fmt.Errorf("Error setting `subscription_key_parameter_names`: %+v", err)
		}
	}

	return nil
}

func resourceArmApiManagementApiDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiid := id.Path["apis"]

	name := apiid
	revision := ""
	if strings.Contains(apiid, ";") {
		name = strings.Split(apiid, ";")[0]
		revision = strings.Split(apiid, "=")[1]
	}

	deleteRevisions := utils.Bool(true)
	if resp, err := client.Delete(ctx, resourceGroup, serviceName, name, "", deleteRevisions); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting API %q / Revision %q (API Management Service %q / Resource Group %q): %s", name, revision, serviceName, resourceGroup, err)
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
