// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationConnectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Connection
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connection.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("read request on %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.ConnectionType != nil {
				d.Set("type", props.ConnectionType.Name)
			}

			if props.FieldDefinitionValues != nil {
				if err := d.Set("values", props.FieldDefinitionValues); err != nil {
					return fmt.Errorf("setting `values`: %+v", err)
				}
			}
			d.Set("description", props.Description)
		}
	}

	return nil
}
