// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

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
	apiId := fmt.Sprintf("%s;rev=%s", d.Get("name").(string), revision)

	id := api.NewApiID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), apiId)
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s does not exist", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("api_management_name", id.ServiceName)
	name := getApiName(id.ApiId)
	d.Set("name", name)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("display_name", pointer.From(props.DisplayName))
			d.Set("is_current", pointer.From(props.IsCurrent))
			d.Set("is_online", pointer.From(props.IsOnline))
			d.Set("path", props.Path)
			d.Set("revision", pointer.From(props.ApiRevision))
			d.Set("service_url", pointer.From(props.ServiceURL))
			d.Set("soap_pass_through", pointer.From(props.Type) == api.ApiTypeSoap)
			d.Set("subscription_required", pointer.From(props.SubscriptionRequired))
			d.Set("version", pointer.From(props.ApiVersion))
			d.Set("version_set_id", pointer.From(props.ApiVersionSetId))
			if err := d.Set("protocols", flattenApiManagementApiDataSourceProtocols(props.Protocols)); err != nil {
				return fmt.Errorf("setting `protocols`: %s", err)
			}

			if err := d.Set("subscription_key_parameter_names", flattenApiManagementApiDataSourceSubscriptionKeyParamNames(props.SubscriptionKeyParameterNames)); err != nil {
				return fmt.Errorf("setting `subscription_key_parameter_names`: %+v", err)
			}
		}
	}
	return nil
}

func flattenApiManagementApiDataSourceProtocols(input *[]api.Protocol) []string {
	if input == nil {
		return []string{}
	}

	results := make([]string, 0)
	for _, v := range *input {
		results = append(results, string(v))
	}

	return results
}

func flattenApiManagementApiDataSourceSubscriptionKeyParamNames(paramNames *api.SubscriptionKeyParameterNamesContract) []interface{} {
	if paramNames == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["header"] = pointer.From(paramNames.Header)
	result["query"] = pointer.From(paramNames.Query)

	return []interface{}{result}
}
