// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2026-04-01/pipelinegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.PipelineGroupsClient

			var model MonitorPipelineResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := pipelinegroups.NewPipelineGroupID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := pipelinegroups.PipelineGroup{
				ExtendedLocation: &pipelinegroups.AzureResourceManagerCommonTypesExtendedLocation{
					Name: model.CustomLocationId,
					Type: pipelinegroups.ExtendedLocationTypeCustomLocation,
				},
				Location: location.Normalize(model.Location),
				Properties: &pipelinegroups.PipelineGroupProperties{
					ExecutionPlacement: expandPipelineExecutionPlacement(model.ExecutionPlacementConstraint),
					Exporters:          expandPipelineAzureMonitorWorkspaceLogExporters(model.AzureMonitorWorkspaceLogExporter),
					Processors:         expandPipelineProcessors(model),
					Receivers:          expandPipelineReceivers(model),
					Replicas:           pointer.ToOrNil(model.ReplicaCount),
					Service:            expandPipelineService(model.LogPipeline, model.PersistentVolumeName),
				},
				Tags: pointer.To(model.Tags),
			}

			if len(model.TlsConfiguration) > 0 {
				payload.Properties.TlsConfigurations = expandPipelineTlsConfigurations(model.TlsConfiguration)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.ResourceData.SetId(id.ID())
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}
