package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceLogicAppWorkflow() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceLogicAppWorkflowRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"logic_app_integration_account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			// TODO: should Parameters be split out into their own object to allow validation on the different sub-types?
			"parameters": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"workflow_schema": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"workflow_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"access_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"connector_endpoint_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},
			"connector_outbound_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},
			"workflow_endpoint_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},
			"workflow_outbound_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceLogicAppWorkflowRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Logic App Workflow %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.WorkflowProperties; props != nil {
		parameters := flattenLogicAppDataSourceWorkflowParameters(props.Parameters)
		if err := d.Set("parameters", parameters); err != nil {
			return fmt.Errorf("Error setting `parameters`: %+v", err)
		}

		d.Set("access_endpoint", props.AccessEndpoint)

		if props.EndpointsConfiguration == nil || props.EndpointsConfiguration.Connector == nil {
			d.Set("connector_endpoint_ip_addresses", []interface{}{})
			d.Set("connector_outbound_ip_addresses", []interface{}{})
		} else {
			d.Set("connector_endpoint_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Connector.AccessEndpointIPAddresses))
			d.Set("connector_outbound_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Connector.OutgoingIPAddresses))
		}

		if props.EndpointsConfiguration == nil || props.EndpointsConfiguration.Workflow == nil {
			d.Set("workflow_endpoint_ip_addresses", []interface{}{})
			d.Set("workflow_outbound_ip_addresses", []interface{}{})
		} else {
			d.Set("workflow_endpoint_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Workflow.AccessEndpointIPAddresses))
			d.Set("workflow_outbound_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Workflow.OutgoingIPAddresses))
		}

		if definition := props.Definition; definition != nil {
			if v, ok := definition.(map[string]interface{}); ok {
				d.Set("workflow_schema", v["$schema"].(string))
				d.Set("workflow_version", v["contentVersion"].(string))
			}
		}

		if props.IntegrationAccount != nil && props.IntegrationAccount.ID != nil {
			d.Set("logic_app_integration_account_id", props.IntegrationAccount.ID)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenLogicAppDataSourceWorkflowParameters(input map[string]*logic.WorkflowParameter) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		if v != nil {
			// we only support string parameters at this time
			val, ok := v.Value.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping parameter %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
}
