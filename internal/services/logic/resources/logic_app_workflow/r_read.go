// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_workflow

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogicAppWorkflowRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workflows.ParseWorkflowID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Logic App Workflow %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %s: %+v", id, err)
	}

	d.Set("name", id.WorkflowName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if loc := model.Location; loc != nil {
			d.Set("location", location.Normalize(*loc))
		}

		identity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
		if err != nil {
			return err
		}
		d.Set("identity", identity)

		if props := model.Properties; props != nil {
			d.Set("access_endpoint", props.AccessEndpoint)

			if err := d.Set("access_control", flattenLogicAppWorkflowFlowAccessControl(props.AccessControl)); err != nil {
				return fmt.Errorf("setting `access_control`: %+v", err)
			}

			if props.State != nil && *props.State != "" {
				d.Set("enabled", *props.State == workflows.WorkflowStateEnabled)
			}

			if props.EndpointsConfiguration == nil || props.EndpointsConfiguration.Connector == nil {
				d.Set("connector_endpoint_ip_addresses", []any{})
				d.Set("connector_outbound_ip_addresses", []any{})
			} else {
				d.Set("connector_endpoint_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Connector.AccessEndpointIPAddresses))
				d.Set("connector_outbound_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Connector.OutgoingIPAddresses))
			}

			if props.EndpointsConfiguration == nil || props.EndpointsConfiguration.Workflow == nil {
				d.Set("workflow_endpoint_ip_addresses", []any{})
				d.Set("workflow_outbound_ip_addresses", []any{})
			} else {
				d.Set("workflow_endpoint_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Workflow.AccessEndpointIPAddresses))
				d.Set("workflow_outbound_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Workflow.OutgoingIPAddresses))
			}
			if definition := props.Definition; definition != nil {
				definitionRaw := *props.Definition
				if v, ok := definitionRaw.(map[string]any); ok {
					if v["$schema"] != nil {
						d.Set("workflow_schema", v["$schema"].(string))
					}
					if v["contentVersion"] != nil {
						d.Set("workflow_version", v["contentVersion"].(string))
					}
					if p, ok := v["parameters"]; ok {
						workflowParameters, err := flattenLogicAppWorkflowWorkflowParameters(p.(map[string]any))
						if err != nil {
							return fmt.Errorf("flattening `workflow_parameters`: %+v", err)
						}
						if err := d.Set("workflow_parameters", workflowParameters); err != nil {
							return fmt.Errorf("setting `workflow_parameters`: %+v", err)
						}

						// The props.Parameters (the value of the param) is accompany with the "parameters" (the definition of the param) inside the props.Definition.
						// We will need to make use of the definition of the parameters in order to properly flatten the value of the parameters being set (for kinds of types).
						parameters, err := flattenLogicAppWorkflowParameters(d, props.Parameters, p.(map[string]any))
						if err != nil {
							return fmt.Errorf("flattening `parameters`: %v", err)
						}
						if err := d.Set("parameters", parameters); err != nil {
							return fmt.Errorf("setting `parameters`: %+v", err)
						}
					}
				}
			}

			integrationServiceEnvironmentId := ""
			if props.IntegrationServiceEnvironment != nil && props.IntegrationServiceEnvironment.Id != nil {
				integrationServiceEnvironmentId = *props.IntegrationServiceEnvironment.Id
			}
			d.Set("integration_service_environment_id", integrationServiceEnvironmentId)

			if props.IntegrationAccount != nil && props.IntegrationAccount.Id != nil {
				d.Set("logic_app_integration_account_id", props.IntegrationAccount.Id)
			}

			integrationAccountId := ""
			if props.IntegrationAccount != nil && props.IntegrationAccount.Id != nil {
				integrationAccountId = *props.IntegrationAccount.Id
			}
			d.Set("logic_app_integration_account_id", integrationAccountId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
