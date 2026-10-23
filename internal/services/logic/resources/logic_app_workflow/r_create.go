// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_workflow

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogicAppWorkflowCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workflows.NewWorkflowID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Logic App Workflow %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_logic_app_workflow", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))

	workflowSchema := d.Get("workflow_schema").(string)
	workflowVersion := d.Get("workflow_version").(string)
	workflowParameters, err := expandLogicAppWorkflowWorkflowParameters(d.Get("workflow_parameters").(map[string]any))
	if err != nil {
		return fmt.Errorf("expanding `workflow_parameters`: %+v", err)
	}

	parameters, err := expandLogicAppWorkflowParameters(d.Get("parameters").(map[string]any), workflowParameters)
	if err != nil {
		return err
	}
	t := d.Get("tags").(map[string]any)

	isEnabled := workflows.WorkflowStateEnabled
	if v := d.Get("enabled").(bool); !v {
		isEnabled = workflows.WorkflowStateDisabled
	}

	identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	var definition any = map[string]any{
		"$schema":        workflowSchema,
		"contentVersion": workflowVersion,
		"actions":        make(map[string]any),
		"triggers":       make(map[string]any),
		"parameters":     workflowParameters,
	}

	properties := workflows.Workflow{
		Identity: identity,
		Location: pointer.To(location),
		Properties: &workflows.WorkflowProperties{
			Definition: &definition,
			Parameters: parameters,
			State:      &isEnabled,
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("access_control"); ok {
		properties.Properties.AccessControl = expandLogicAppWorkflowAccessControl(v.([]any))
	}

	if iseID, ok := d.GetOk("integration_service_environment_id"); ok {
		properties.Properties.IntegrationServiceEnvironment = &workflows.ResourceReference{
			Id: pointer.To(iseID.(string)),
		}
	}

	if v, ok := d.GetOk("logic_app_integration_account_id"); ok {
		properties.Properties.IntegrationAccount = &workflows.ResourceReference{
			Id: pointer.To(v.(string)),
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceLogicAppWorkflowRead(d, meta)
}
