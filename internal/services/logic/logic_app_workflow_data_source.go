// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

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

			"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

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

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceLogicAppWorkflowRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workflows.NewWorkflowID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("Logic App Workflow %s was not found", id)
		}

		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		identity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
		if err != nil {
			return err
		}
		d.Set("identity", identity)

		if props := model.Properties; props != nil {
			parameters := flattenLogicAppDataSourceWorkflowParameters(props.Parameters)
			if err := d.Set("parameters", parameters); err != nil {
				return fmt.Errorf("setting `parameters`: %+v", err)
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
				definitionRaw := *props.Definition
				if v, ok := definitionRaw.(map[string]interface{}); ok {
					d.Set("workflow_schema", v["$schema"].(string))
					d.Set("workflow_version", v["contentVersion"].(string))
				}
			}

			if props.IntegrationAccount != nil && props.IntegrationAccount.Id != nil {
				d.Set("logic_app_integration_account_id", props.IntegrationAccount.Id)
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func flattenLogicAppDataSourceWorkflowParameters(input *map[string]workflows.WorkflowParameter) map[string]interface{} {
	output := make(map[string]interface{})
	if input == nil {
		return output
	}

	for k, v := range *input {
		// we only support string parameters at this time
		if v.Value != nil {
			rawValue := *v.Value
			val, ok := rawValue.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping parameter %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
}
