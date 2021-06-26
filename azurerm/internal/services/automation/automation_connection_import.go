package automation

import (
	"context"
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func importAutomationConnection(connectionType string) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.ConnectionID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		client := meta.(*clients.Client).Automation.ConnectionClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving automation connection %q (Account %q / Resource Group %q): %+v", id.Name, id.AutomationAccountName, id.ResourceGroup, err)
		}

		if resp.ConnectionProperties == nil || resp.ConnectionProperties.ConnectionType == nil || resp.ConnectionProperties.ConnectionType.Name == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving automation connection %q (Account %q / Resource Group %q): `properties`, `properties.connectionType` or `properties.connectionType.name` was nil", id.Name, id.AutomationAccountName, id.ResourceGroup)
		}

		if *resp.ConnectionProperties.ConnectionType.Name != connectionType {
			return nil, fmt.Errorf(`automation connection "type" mismatch, expected "%s", got "%s"`, connectionType, *resp.ConnectionProperties.ConnectionType.Name)
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}
