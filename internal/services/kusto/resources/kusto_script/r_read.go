// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_script

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/scripts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoDatabaseScriptRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.ScriptsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scripts.ParseScriptID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.ScriptName)
	d.Set("database_id", commonids.NewKustoDatabaseID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.DatabaseName).ID())

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil {
			d.Set("continue_on_errors_enabled", props.ContinueOnErrors)
			d.Set("force_an_update_when_value_changed", props.ForceUpdateTag)
			d.Set("url", props.ScriptURL)
			if props.ScriptLevel != nil {
				d.Set("script_level", pointer.FromEnum(props.ScriptLevel))
			}
			if props.PrincipalPermissionsAction != nil {
				d.Set("principal_permissions_action", pointer.FromEnum(props.PrincipalPermissionsAction))
			}
		}
	}
	return nil
}
