// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_dsc_nodeconfiguration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/dscnodeconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationDscNodeConfigurationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.DscNodeConfiguration
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dscnodeconfiguration.ParseNodeConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.NodeConfigurationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			configurationName := ""
			if props.Configuration != nil && props.Configuration.Name != nil {
				configurationName = *props.Configuration.Name
			}
			d.Set("configuration_name", configurationName)
		}
	}

	// cannot read back content_embedded as not part of body nor exposed through method

	return nil
}
