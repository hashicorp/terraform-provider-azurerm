package logic_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func actionExists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	return componentExists(ctx, clients, state, "Action", "actions")
}

func triggerExists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	return componentExists(ctx, clients, state, "Trigger", "triggers")
}

func componentExists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState, kind, propertyName string) (*bool, error) {
	var resourceGroup string
	var workflowName string
	var name string

	if kind == "Action" {
		id, err := parse.ActionID(state.ID)
		if err != nil {
			return nil, err
		}
		resourceGroup = id.ResourceGroup
		workflowName = id.WorkflowName
		name = id.Name
	} else {
		id, err := parse.TriggerID(state.ID)
		if err != nil {
			return nil, err
		}
		resourceGroup = id.ResourceGroup
		workflowName = id.WorkflowName
		name = id.Name
	}

	resp, err := clients.Logic.WorkflowClient.Get(ctx, resourceGroup, workflowName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Logic App Workflow %s %s (resource group: %s): %v", kind, workflowName, resourceGroup, err)
	}

	if resp.WorkflowProperties == nil {
		return utils.Bool(false), nil
	}

	if resp.WorkflowProperties.Definition == nil {
		return nil, fmt.Errorf("Logic App Workflow %s %s (resource group: %s) Definition is nil", kind, workflowName, resourceGroup)
	}

	definition := resp.WorkflowProperties.Definition.(map[string]interface{})
	actions := definition[propertyName].(map[string]interface{})

	exists := false
	for k := range actions {
		if strings.EqualFold(k, name) {
			exists = true
			break
		}
	}

	return utils.Bool(exists), nil
}
