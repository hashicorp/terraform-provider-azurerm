package logic_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func actionExists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	return componentExists(ctx, clients, state, "Action", "actions")
}

func triggerExists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	return componentExists(ctx, clients, state, "Trigger", "triggers")
}

func componentExists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState, kind, propertyName string) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	workflowName := id.Path["workflows"]
	componentName := id.Path[propertyName]

	resp, err := clients.Logic.WorkflowClient.Get(ctx, id.ResourceGroup, workflowName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Logic App Workflow %s %s (resource group: %s): %v", kind, workflowName, id.ResourceGroup, err)
	}

	if resp.WorkflowProperties == nil {
		return utils.Bool(false), nil
	}

	if resp.WorkflowProperties.Definition == nil {
		return nil, fmt.Errorf("Logic App Workflow %s %s (resource group: %s) Definition is nil", kind, workflowName, id.ResourceGroup)
	}

	definition := resp.WorkflowProperties.Definition.(map[string]interface{})
	actions := definition[propertyName].(map[string]interface{})

	exists := false
	for k := range actions {
		if strings.EqualFold(k, componentName) {
			exists = true
			break
		}
	}

	return utils.Bool(exists), nil
}
