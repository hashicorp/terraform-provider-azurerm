// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_script

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/scripts"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoDatabaseScriptCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.ScriptsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	databaseId, err := commonids.ParseKustoDatabaseID(d.Get("database_id").(string))
	if err != nil {
		return err
	}
	id := scripts.NewScriptID(databaseId.SubscriptionId, databaseId.ResourceGroupName, databaseId.KustoClusterName, databaseId.KustoDatabaseName, d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %q: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_kusto_script", id.ID())
			}
		}
	}

	clusterId := commonids.NewKustoClusterID(databaseId.SubscriptionId, databaseId.ResourceGroupName, databaseId.KustoClusterName)
	locks.ByID(clusterId.ID())
	defer locks.UnlockByID(clusterId.ID())

	forceUpdateTag := d.Get("force_an_update_when_value_changed").(string)
	if len(forceUpdateTag) == 0 {
		forceUpdateTag, _ = uuid.GenerateUUID()
	}

	parameters := scripts.Script{
		Properties: &scripts.ScriptProperties{
			ContinueOnErrors: pointer.To(d.Get("continue_on_errors_enabled").(bool)),
			ForceUpdateTag:   pointer.To(forceUpdateTag),
		},
	}

	if scriptURL, ok := d.GetOk("url"); ok {
		parameters.Properties.ScriptURL = pointer.To(scriptURL.(string))
	}

	if scriptURLSasToken, ok := d.GetOk("sas_token"); ok {
		parameters.Properties.ScriptURLSasToken = pointer.To(scriptURLSasToken.(string))
	}

	if scriptContent, ok := d.GetOk("script_content"); ok {
		parameters.Properties.ScriptContent = pointer.To(scriptContent.(string))
	}

	if scriptLevel, ok := d.GetOk("script_level"); ok {
		parameters.Properties.ScriptLevel = pointer.ToEnum[scripts.ScriptLevel](scriptLevel.(string))
	}

	if principalPermissionsAction, ok := d.GetOk("principal_permissions_action"); ok {
		parameters.Properties.PrincipalPermissionsAction = pointer.ToEnum[scripts.PrincipalPermissionsAction](principalPermissionsAction.(string))
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceKustoDatabaseScriptRead(d, meta)
}
