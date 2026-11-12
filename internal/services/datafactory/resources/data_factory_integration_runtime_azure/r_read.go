// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_azure

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryIntegrationRuntimeAzureRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntimes.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName)

	resp, err := client.Get(ctx, *id, integrationruntimes.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.IntegrationRuntimeName)
	d.Set("data_factory_id", dataFactoryId.ID())

	if model := resp.Model; model != nil {
		runTime, ok := model.Properties.(integrationruntimes.ManagedIntegrationRuntime)
		if !ok {
			return fmt.Errorf("asserting `IntegrationRuntime` as `ManagedIntegrationRuntime` for %s", *id)
		}

		d.Set("description", pointer.From(runTime.Description))
		d.Set("virtual_network_enabled", runTime.ManagedVirtualNetwork != nil && runTime.ManagedVirtualNetwork.ReferenceName != "")

		if computeProps := runTime.TypeProperties.ComputeProperties; computeProps != nil {
			d.Set("location", location.NormalizeNilable(computeProps.Location))

			if dataFlowProps := computeProps.DataFlowProperties; dataFlowProps != nil {
				d.Set("compute_type", pointer.FromEnum(dataFlowProps.ComputeType))
				d.Set("core_count", dataFlowProps.CoreCount)
				d.Set("time_to_live_min", dataFlowProps.TimeToLive)
				d.Set("cleanup_enabled", dataFlowProps.Cleanup)
			}
		}

		// @sreallymatt: This is currently non-functional, the API doesn't return the InteractiveQuery properties
		// See: https://github.com/Azure/azure-rest-api-specs/issues/39594
		// Can be uncommented and the additional `GetStatus` request below removed once this issue is resolved.
		// if iqProps := runTime.TypeProperties.InteractiveQuery; iqProps != nil {
		//	d.Set("interactive_authoring_time_to_live_in_minutes", iqProps.AutoTerminationMinutes)
		// }
	}

	status, err := client.GetStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving status for %s: %+v", id, err)
	}

	if model := status.Model; model != nil {
		rt, ok := model.Properties.(integrationruntimes.ManagedIntegrationRuntimeStatus)
		if !ok {
			return fmt.Errorf("asserting `IntegrationRuntimeStatus` as `ManagedIntegrationRuntimeStatus` for %s", *id)
		}

		if iq := rt.TypeProperties.InteractiveQuery; iq != nil {
			d.Set("interactive_authoring_time_to_live_in_minutes", iq.AutoTerminationMinutes)
			// The API doesn't reset the TTL value when interactive authoring is disabled. If disabled, reset TTL.
			if pointer.From(iq.Status) == integrationruntimes.InteractiveCapabilityStatusDisabled {
				d.Set("interactive_authoring_time_to_live_in_minutes", nil)
			}
		}
	}

	return nil
}
