package automation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func importAutomationConnection(connectionType string) func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	return func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
		id, err := parse.ConnectionID(d.Id())
		if err != nil {
			return []*schema.ResourceData{}, err
		}

		client := meta.(*clients.Client).Automation.ConnectionClient
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.ConnectionName)
		if err != nil {
			return []*schema.ResourceData{}, fmt.Errorf("retrieving automation connection %q (Account %q / Resource Group %q): %+v", id.ConnectionName, id.AutomationAccountName, id.ResourceGroup, err)
		}

		if resp.ConnectionProperties == nil || resp.ConnectionProperties.ConnectionType == nil || resp.ConnectionProperties.ConnectionType.Name == nil {
			return []*schema.ResourceData{}, fmt.Errorf("retrieving automation connection %q (Account %q / Resource Group %q): `properties`, `properties.connectionType` or `properties.connectionType.name` was nil", id.ConnectionName, id.AutomationAccountName, id.ResourceGroup)
		}

		if *resp.ConnectionProperties.ConnectionType.Name != connectionType {
			return nil, fmt.Errorf(`automation connection "type" mismatch, expected "%s", got "%s"`, connectionType, *resp.ConnectionProperties.ConnectionType.Name)
		}
		return []*schema.ResourceData{d}, nil
	}
}
