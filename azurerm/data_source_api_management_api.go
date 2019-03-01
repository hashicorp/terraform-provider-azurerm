package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApiManagementApi() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApiManagementApiRead,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"protocols": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},

			"subscription_key_parameter_names": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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

			"soap_pass_through": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"revision": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
				Computed: true,
			},

			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"api_version_set_id": {
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

func dataSourceApiManagementApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementApiClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("service_name").(string)
	name := d.Get("name").(string)
	revision := int32(d.Get("revision").(int))

	apiId := fmt.Sprintf("%s;rev=%d", name, revision)

	resp, err := client.Get(ctx, resGroup, serviceName, apiId)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("API Management API %q (Service %q / Resource Group %q) was not found", name, serviceName, resGroup)
		}

		return fmt.Errorf("Error retrieving API Management API %q (Service %q / Resource Group %q): %+v", name, serviceName, resGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("service_name", serviceName)
	d.Set("resource_group_name", resGroup)

	if props := resp.APIContractProperties; props != nil {
		d.Set("service_url", props.ServiceURL)
		d.Set("path", props.Path)
		d.Set("description", props.Description)
		d.Set("revision", props.APIRevision)
		d.Set("api_version", props.APIVersion)
		d.Set("api_version_set_id", props.APIVersionSetID)
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
