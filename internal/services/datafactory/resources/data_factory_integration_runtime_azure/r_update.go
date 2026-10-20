// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_azure

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimedisableinteractivequery"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimeenableinteractivequery"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryIntegrationRuntimeAzureUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	disableInteractiveQueryClient := meta.(*clients.Client).DataFactory.IntegrationRuntimeDisableInteractiveQueryClient
	enableInteractiveQueryClient := meta.(*clients.Client).DataFactory.IntegrationRuntimeEnableInteractiveQueryClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntimes.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, integrationruntimes.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `Model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `Properties` was nil", id)
	}

	props, ok := existing.Model.Properties.(integrationruntimes.ManagedIntegrationRuntime)
	if !ok {
		return fmt.Errorf("retrieving %s: asserting `IntegrationRuntime` as `ManagedIntegrationRuntime`", id)
	}

	if props.TypeProperties.ComputeProperties == nil {
		props.TypeProperties.ComputeProperties = &integrationruntimes.IntegrationRuntimeComputeProperties{}
	}

	if props.TypeProperties.ComputeProperties.DataFlowProperties == nil {
		props.TypeProperties.ComputeProperties.DataFlowProperties = &integrationruntimes.IntegrationRuntimeDataFlowProperties{}
	}
	dfProps := props.TypeProperties.ComputeProperties.DataFlowProperties

	if d.HasChange("cleanup_enabled") {
		dfProps.Cleanup = pointer.To(d.Get("cleanup_enabled").(bool))
	}

	if d.HasChange("compute_type") {
		dfProps.ComputeType = pointer.ToEnum[integrationruntimes.DataFlowComputeType](d.Get("compute_type").(string))
	}

	if d.HasChange("core_count") {
		dfProps.CoreCount = pointer.To(int64(d.Get("core_count").(int)))
	}

	if d.HasChange("description") {
		props.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("time_to_live_min") {
		dfProps.TimeToLive = pointer.To(int64(d.Get("time_to_live_min").(int)))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model, integrationruntimes.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if d.HasChange("interactive_authoring_time_to_live_in_minutes") {
		poller := pollers.NewPoller(custompollers.NewDataFactoryIntegrationRuntimeStatusPoller(client, *id), time.Second*5, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for state change of %s: %+v", id, err)
		}

		if ttl := d.Get("interactive_authoring_time_to_live_in_minutes").(int); ttl > 0 {
			payload := integrationruntimeenableinteractivequery.EnableInteractiveQueryRequest{
				AutoTerminationMinutes: pointer.To(int64(ttl)),
			}

			if err := enableInteractiveQueryClient.IntegrationRuntimeEnableInteractiveQueryThenPoll(ctx, integrationruntimeenableinteractivequery.IntegrationRuntimeId(*id), payload); err != nil {
				return fmt.Errorf("enabling interactive authoring on %s: %+v", id, err)
			}
		} else {
			if err := disableInteractiveQueryClient.IntegrationRuntimeDisableInteractiveQueryThenPoll(ctx, integrationruntimedisableinteractivequery.IntegrationRuntimeId(*id)); err != nil {
				return fmt.Errorf("disabling interactive authoring on %s: %+v", id, err)
			}
		}
	}

	return resourceDataFactoryIntegrationRuntimeAzureRead(d, meta)
}
