// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_integration_runtime_self_hosted

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/integrationruntimes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseIntegrationRuntimeSelfHostedRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.IntegrationRuntimesClient
	authKeysClient := meta.(*clients.Client).Synapse.IntegrationRuntimeAuthKeysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntimes.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.WorkspaceName, id.IntegrationRuntimeName, "")
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.IntegrationRuntimeName)
	d.Set("synapse_workspace_id", workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	selfHostedIntegrationRuntime, convertSuccess := resp.Properties.AsSelfHostedIntegrationRuntime()

	if !convertSuccess {
		return fmt.Errorf("converting integration runtime to Self-Hosted integration runtime (%q)", id)
	}

	d.Set("description", selfHostedIntegrationRuntime.Description)

	respKey, err := authKeysClient.List(ctx, id.ResourceGroupName, id.WorkspaceName, id.IntegrationRuntimeName)
	if err != nil {
		if response.WasNotFound(respKey.Response.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving auth keys (%q): %+v", id, err)
	}

	d.Set("authorization_key_primary", respKey.AuthKey1)
	d.Set("authorization_key_secondary", respKey.AuthKey2)

	return nil
}
