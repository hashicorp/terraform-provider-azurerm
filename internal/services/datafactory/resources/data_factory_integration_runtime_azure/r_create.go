// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_azure

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimeenableinteractivequery"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryIntegrationRuntimeAzureCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	enableInteractiveQueryClient := meta.(*clients.Client).DataFactory.IntegrationRuntimeEnableInteractiveQueryClient
	managedVirtualNetworksClient := meta.(*clients.Client).DataFactory.ManagedVirtualNetworks

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := integrationruntimes.NewIntegrationRuntimeID(dataFactoryId.SubscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, integrationruntimes.DefaultGetOperationOptions())
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_factory_integration_runtime_azure", id.ID())
		}
	}

	managedIntegrationRuntime := integrationruntimes.ManagedIntegrationRuntime{
		Description: pointer.To(d.Get("description").(string)),
		Type:        integrationruntimes.IntegrationRuntimeTypeManaged,
		TypeProperties: integrationruntimes.ManagedIntegrationRuntimeTypeProperties{
			ComputeProperties: &integrationruntimes.IntegrationRuntimeComputeProperties{
				Location: pointer.To(location.Normalize(d.Get("location").(string))),
				DataFlowProperties: &integrationruntimes.IntegrationRuntimeDataFlowProperties{
					ComputeType: pointer.ToEnum[integrationruntimes.DataFlowComputeType](d.Get("compute_type").(string)),
					CoreCount:   pointer.To(int64(d.Get("core_count").(int))),
					TimeToLive:  pointer.To(int64(d.Get("time_to_live_min").(int))),
					Cleanup:     pointer.To(d.Get("cleanup_enabled").(bool)),
				},
			},
		},
	}

	if d.Get("virtual_network_enabled").(bool) {
		virtualNetworkName, err := helpers.GetManagedVirtualNetworkName(ctx, managedVirtualNetworksClient, id.SubscriptionId, id.ResourceGroupName, id.FactoryName)
		if err != nil {
			return err
		}
		if virtualNetworkName == nil {
			return fmt.Errorf("virtual network feature for azure integration runtime is only available after managed virtual network for this data factory is enabled")
		}
		managedIntegrationRuntime.ManagedVirtualNetwork = &integrationruntimes.ManagedVirtualNetworkReference{
			Type:          integrationruntimes.ManagedVirtualNetworkReferenceTypeManagedVirtualNetworkReference,
			ReferenceName: *virtualNetworkName,
		}
	}

	integrationRuntime := integrationruntimes.IntegrationRuntimeResource{
		Name:       &id.IntegrationRuntimeName,
		Properties: managedIntegrationRuntime,
	}

	if _, err := client.CreateOrUpdate(ctx, id, integrationRuntime, integrationruntimes.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if ttl := d.Get("interactive_authoring_time_to_live_in_minutes").(int); ttl != 0 {
		// Interactive Authoring/Query can only be modified once the integration runtime is online
		poller := pollers.NewPoller(custompollers.NewDataFactoryIntegrationRuntimeStatusPoller(client, id), time.Second*5, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for state change of %s: %+v", id, err)
		}

		payload := integrationruntimeenableinteractivequery.EnableInteractiveQueryRequest{
			AutoTerminationMinutes: pointer.To(int64(ttl)),
		}

		if err := enableInteractiveQueryClient.IntegrationRuntimeEnableInteractiveQueryThenPoll(ctx, integrationruntimeenableinteractivequery.IntegrationRuntimeId(id), payload); err != nil {
			return fmt.Errorf("enabling interactive authoring on %s: %+v", id, err)
		}
	}

	return resourceDataFactoryIntegrationRuntimeAzureRead(d, meta)
}
