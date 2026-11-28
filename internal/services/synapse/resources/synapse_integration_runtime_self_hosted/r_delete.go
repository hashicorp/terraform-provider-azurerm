// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_integration_runtime_self_hosted

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseIntegrationRuntimeSelfHostedDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntimes.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroupName, id.WorkspaceName, id.IntegrationRuntimeName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}
