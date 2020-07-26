package apimanagement

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApiManagementApi() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApiManagementApiRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ApiManagementApiName,
			},

			"api_management_name": azure.SchemaApiManagementDataSourceName(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"revision": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
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

			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"protocols": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"service_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"soap_pass_through": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"subscription_key_parameter_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"subscription_required": {
				Type:     schema.TypeBool,
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
		},
	}
}

func dataSourceApiManagementApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	name := d.Get("name").(string)
	revision := d.Get("revision").(string)

	apiId := fmt.Sprintf("%s;rev=%s", name, revision)
	resp, err := client.Get(ctx, resourceGroup, serviceName, apiId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("API %q Revision %q (API Management Service %q / Resource Group %q) does not exist!", name, revision, serviceName, resourceGroup)
		}

		return fmt.Errorf("retrieving API %q / Revision %q (API Management Service %q / Resource Group %q): %+v", name, revision, serviceName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("api_management_name", serviceName)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.APIContractProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("is_current", props.IsCurrent)
		d.Set("is_online", props.IsOnline)
		d.Set("path", props.Path)
		d.Set("revision", props.APIRevision)
		d.Set("service_url", props.ServiceURL)
		d.Set("soap_pass_through", string(props.APIType) == string(apimanagement.SoapPassThrough))
		d.Set("subscription_required", props.SubscriptionRequired)
		d.Set("version", props.APIVersion)
		d.Set("version_set_id", props.APIVersionSetID)

		if err := d.Set("protocols", flattenApiManagementApiDataSourceProtocols(props.Protocols)); err != nil {
			return fmt.Errorf("setting `protocols`: %s", err)
		}

		if err := d.Set("subscription_key_parameter_names", flattenApiManagementApiDataSourceSubscriptionKeyParamNames(props.SubscriptionKeyParameterNames)); err != nil {
			return fmt.Errorf("setting `subscription_key_parameter_names`: %+v", err)
		}
	}

	return nil
}

func flattenApiManagementApiDataSourceProtocols(input *[]apimanagement.Protocol) []string {
	if input == nil {
		return []string{}
	}

	results := make([]string, 0)
	for _, v := range *input {
		results = append(results, string(v))
	}

	return results
}
func flattenApiManagementApiDataSourceSubscriptionKeyParamNames(paramNames *apimanagement.SubscriptionKeyParameterNamesContract) []interface{} {
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
