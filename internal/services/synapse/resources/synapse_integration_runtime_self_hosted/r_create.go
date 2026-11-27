// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_integration_runtime_self_hosted

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/integrationruntimes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseIntegrationRuntimeSelfHostedCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	workspaceId, err := workspaces.ParseWorkspaceIDInsensitively(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	id := integrationruntimes.NewIntegrationRuntimeID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroupName, id.WorkspaceName, id.IntegrationRuntimeName, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_synapse_integration_runtime_self_hosted", id.ID())
			}
		}
	}

	integrationRuntime := synapse.IntegrationRuntimeResource{
		Name: pointer.To(id.IntegrationRuntimeName),
		Properties: synapse.SelfHostedIntegrationRuntime{
			Description: pointer.To(d.Get("description").(string)),
			Type:        synapse.TypeBasicIntegrationRuntimeTypeSelfHosted,
		},
	}

	future, err := client.Create(ctx, id.ResourceGroupName, id.WorkspaceName, id.IntegrationRuntimeName, integrationRuntime, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for %s: %+v", id, err)
	}

	return resourceSynapseIntegrationRuntimeSelfHostedRead(d, meta)
}
