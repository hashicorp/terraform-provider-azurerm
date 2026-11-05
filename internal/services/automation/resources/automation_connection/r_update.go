// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationConnectionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Connection
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connection.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	parameters := connection.ConnectionUpdateParameters{
		Name:       pointer.To(id.ConnectionName),
		Properties: &connection.ConnectionUpdateProperties{},
	}

	if d.HasChange("description") {
		parameters.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("values") {
		values := helpers.ExpandStringInterfaceMap(d.Get("values").(map[string]any))
		parameters.Properties.FieldDefinitionValues = &values
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAutomationConnectionRead(d, meta)
}
