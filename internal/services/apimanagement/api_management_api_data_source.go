package apimanagement

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceApiManagementApi() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceApiManagementApiRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ApiManagementApiName,
			},

			"api_management_name": schemaz.SchemaApiManagementDataSourceName(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"revision": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"is_current": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"is_online": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"path": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"protocols": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"service_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"soap_pass_through": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"subscription_key_parameter_names": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"header": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"query": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"subscription_required": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"version_set_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApiManagementApiRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	revision := d.Get("revision").(string)
	id := parse.NewApiID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	apiId := fmt.Sprintf("%s;rev=%s", id.Name, revision)
	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, apiId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("api %s Revision %q does not exist", id, revision)
		}

		return fmt.Errorf("retrieving API %s / Revision %q : %+v", id, revision, err)
	}

	d.SetId(id.ID())

	d.Set("api_management_name", id.ServiceName)
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.APIContractProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("is_current", props.IsCurrent)
		d.Set("is_online", props.IsOnline)
		d.Set("path", props.Path)
		d.Set("revision", props.APIRevision)
		d.Set("service_url", props.ServiceURL)
		d.Set("soap_pass_through", string(props.APIType) == string(apimanagement.SoapAPITypeSoapPassThrough))
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
