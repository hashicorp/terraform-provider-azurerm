// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_integration_runtime_azure

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/integrationruntimes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseIntegrationRuntimeAzureRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.IntegrationRuntimesClient
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

	managedIntegrationRuntime, convertSuccess := resp.Properties.AsManagedIntegrationRuntime()
	if !convertSuccess {
		return fmt.Errorf("converting managed integration runtime to Azure integration runtime (%q)", id)
	}

	d.Set("description", managedIntegrationRuntime.Description)
	if computeProps := managedIntegrationRuntime.ComputeProperties; computeProps != nil {
		d.Set("location", computeProps.Location)
		if dataFlowProps := computeProps.DataFlowProperties; dataFlowProps != nil {
			d.Set("compute_type", dataFlowProps.ComputeType)
			d.Set("core_count", dataFlowProps.CoreCount)
			d.Set("time_to_live_min", dataFlowProps.TimeToLive)
		}
	}

	return nil
}
