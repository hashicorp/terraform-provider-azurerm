// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2026-04-01/pipelinegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.PipelineGroupsClient

			id, err := pipelinegroups.ParsePipelineGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model MonitorPipelineResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", *id)
			}

			payload := *resp.Model

			if metadata.ResourceData.HasChange("execution_placement_constraint") {
				payload.Properties.ExecutionPlacement = expandPipelineExecutionPlacement(model.ExecutionPlacementConstraint)
			}

			if metadata.ResourceData.HasChange("azure_monitor_workspace_log_exporter") {
				payload.Properties.Exporters = expandPipelineAzureMonitorWorkspaceLogExporters(model.AzureMonitorWorkspaceLogExporter)
			}

			batchProcessorChanged := metadata.ResourceData.HasChange("batch_processor")
			microsoftCommonSecurityLogProcessorChanged := metadata.ResourceData.HasChange("microsoft_common_security_log_processor")
			microsoftSyslogProcessorChanged := metadata.ResourceData.HasChange("microsoft_syslog_processor")
			transformLanguageProcessorChanged := metadata.ResourceData.HasChange("transform_language_processor")
			if batchProcessorChanged || microsoftCommonSecurityLogProcessorChanged || microsoftSyslogProcessorChanged || transformLanguageProcessorChanged {
				changedTypes := map[pipelinegroups.ProcessorType]bool{
					pipelinegroups.ProcessorTypeBatch:                      batchProcessorChanged,
					pipelinegroups.ProcessorTypeMicrosoftCommonSecurityLog: microsoftCommonSecurityLogProcessorChanged,
					pipelinegroups.ProcessorTypeMicrosoftSyslog:            microsoftSyslogProcessorChanged,
					pipelinegroups.ProcessorTypeTransformLanguage:          transformLanguageProcessorChanged,
				}
				payload.Properties.Processors = mergePipelineProcessors(payload.Properties.Processors, expandPipelineProcessors(model), changedTypes)
			}

			otlpReceiverChanged := metadata.ResourceData.HasChange("otlp_receiver")
			syslogReceiverChanged := metadata.ResourceData.HasChange("syslog_receiver")
			if otlpReceiverChanged || syslogReceiverChanged {
				changedTypes := map[pipelinegroups.ReceiverType]bool{
					pipelinegroups.ReceiverTypeOTLP:   otlpReceiverChanged,
					pipelinegroups.ReceiverTypeSyslog: syslogReceiverChanged,
				}
				payload.Properties.Receivers = mergePipelineReceivers(payload.Properties.Receivers, expandPipelineReceivers(model), changedTypes)
			}

			if metadata.ResourceData.HasChange("replica_count") {
				payload.Properties.Replicas = pointer.ToOrNil(model.ReplicaCount)
			}

			if metadata.ResourceData.HasChanges("log_pipeline", "persistent_volume_name") {
				payload.Properties.Service = expandPipelineService(model.LogPipeline, model.PersistentVolumeName)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("tls_configuration") {
				payload.Properties.TlsConfigurations = expandPipelineTlsConfigurations(model.TlsConfiguration)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func mergePipelineProcessors(existing, configured []pipelinegroups.Processor, changedTypes map[pipelinegroups.ProcessorType]bool) []pipelinegroups.Processor {
	output := make([]pipelinegroups.Processor, 0, len(existing)+len(configured))
	for _, processor := range existing {
		if !changedTypes[processor.Type] {
			output = append(output, processor)
		}
	}
	for _, processor := range configured {
		if changedTypes[processor.Type] {
			output = append(output, processor)
		}
	}
	return output
}

func mergePipelineReceivers(existing, configured []pipelinegroups.Receiver, changedTypes map[pipelinegroups.ReceiverType]bool) []pipelinegroups.Receiver {
	output := make([]pipelinegroups.Receiver, 0, len(existing)+len(configured))
	for _, receiver := range existing {
		if !changedTypes[receiver.Type] {
			output = append(output, receiver)
		}
	}
	for _, receiver := range configured {
		if changedTypes[receiver.Type] {
			output = append(output, receiver)
		}
	}
	return output
}
