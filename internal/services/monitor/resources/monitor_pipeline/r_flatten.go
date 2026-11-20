// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_pipeline

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2026-04-01/pipelinegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *pipelinegroups.PipelineGroupId, model *pipelinegroups.PipelineGroup) error {
	state := MonitorPipelineResourceModel{
		Name:              id.PipelineGroupName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		if model.ExtendedLocation != nil && model.ExtendedLocation.Name != "" {
			customLocationId, err := customlocations.ParseCustomLocationIDInsensitively(model.ExtendedLocation.Name)
			if err != nil {
				return err
			}
			state.CustomLocationId = customLocationId.ID()
		}

		if props := model.Properties; props != nil {
			state.ExecutionPlacementConstraint = flattenPipelineExecutionPlacement(props.ExecutionPlacement)
			state.AzureMonitorWorkspaceLogExporter = flattenPipelineAzureMonitorWorkspaceLogExporters(props.Exporters)
			flattenPipelineProcessors(props.Processors, &state)
			flattenPipelineReceivers(props.Receivers, &state)
			state.ReplicaCount = pointer.From(props.Replicas)
			flattenPipelineService(props.Service, &state)
			state.TlsConfiguration = flattenPipelineTlsConfigurations(props.TlsConfigurations)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenPipelineExecutionPlacement(input *pipelinegroups.ExecutionPlacement) []PipelinePlacementConstraintModel {
	if input == nil || input.Constraints == nil {
		return []PipelinePlacementConstraintModel{}
	}

	output := make([]PipelinePlacementConstraintModel, 0, len(*input.Constraints))
	for _, constraint := range *input.Constraints {
		output = append(output, PipelinePlacementConstraintModel{
			Capability: constraint.Capability,
			Operator:   string(constraint.Operator),
			Values:     pointer.From(constraint.Values),
		})
	}

	return output
}

func flattenPipelineAzureMonitorWorkspaceLogExporters(input []pipelinegroups.Exporter) []PipelineAzureMonitorWorkspaceLogExporterModel {
	output := make([]PipelineAzureMonitorWorkspaceLogExporterModel, 0, len(input))

	for _, v := range input {
		if v.Type != pipelinegroups.ExporterTypeAzureMonitorWorkspaceLogs || v.AzureMonitorWorkspaceLogs == nil {
			continue
		}

		exporter := flattenPipelineAzureMonitorWorkspaceLogExporter(*v.AzureMonitorWorkspaceLogs)
		exporter.Name = v.Name
		output = append(output, exporter)
	}

	return output
}

func flattenPipelineAzureMonitorWorkspaceLogExporter(input pipelinegroups.AzureMonitorWorkspaceLogsExporter) PipelineAzureMonitorWorkspaceLogExporterModel {
	output := PipelineAzureMonitorWorkspaceLogExporterModel{
		Api: []PipelineAzureMonitorWorkspaceLogsApiModel{
			{
				DataCollectionEndpointUrl:     input.Api.DataCollectionEndpointURL,
				DataCollectionRuleImmutableId: input.Api.DataCollectionRule,
				Stream:                        input.Api.Stream,
				Schema:                        flattenPipelineSchemaMap(input.Api.Schema),
			},
		},
	}

	if input.Persistence != nil {
		output.PersistenceMaximumStorageUsageInGB = pointer.From(input.Persistence.MaxStorageUsage)
		output.PersistenceRetentionPeriodInMinutes = pointer.From(input.Persistence.RetentionPeriod)
	}

	return output
}

func flattenPipelineSchemaMap(input pipelinegroups.SchemaMap) []PipelineSchemaMapModel {
	output := PipelineSchemaMapModel{
		RecordMap: flattenPipelineFieldMaps(input.RecordMap),
	}

	if input.ResourceMap != nil {
		resourceMaps := make([]PipelineFieldMapModel, 0, len(*input.ResourceMap))
		for _, v := range *input.ResourceMap {
			resourceMaps = append(resourceMaps, PipelineFieldMapModel{
				From: v.From,
				To:   v.To,
			})
		}
		output.ResourceMap = resourceMaps
	}

	if input.ScopeMap != nil {
		scopeMaps := make([]PipelineFieldMapModel, 0, len(*input.ScopeMap))
		for _, v := range *input.ScopeMap {
			scopeMaps = append(scopeMaps, PipelineFieldMapModel{
				From: v.From,
				To:   v.To,
			})
		}
		output.ScopeMap = scopeMaps
	}

	return []PipelineSchemaMapModel{
		output,
	}
}

func flattenPipelineFieldMaps(input []pipelinegroups.RecordMap) []PipelineFieldMapModel {
	output := make([]PipelineFieldMapModel, 0, len(input))
	for _, v := range input {
		output = append(output, PipelineFieldMapModel{
			From: v.From,
			To:   v.To,
		})
	}
	return output
}

func flattenPipelineProcessors(input []pipelinegroups.Processor, output *MonitorPipelineResourceModel) {
	for _, v := range input {
		switch v.Type {
		case pipelinegroups.ProcessorTypeBatch:
			processor := PipelineBatchProcessorModel{
				Name:                  v.Name,
				BatchSize:             defaultPipelineBatchSize,
				TimeoutInMilliseconds: defaultPipelineTimeoutInMilliseconds,
			}
			if v.Batch != nil {
				processor.BatchSize = pointer.From(v.Batch.BatchSize)
				processor.TimeoutInMilliseconds = pointer.From(v.Batch.Timeout)
			}
			output.BatchProcessor = append(output.BatchProcessor, processor)
		case pipelinegroups.ProcessorTypeMicrosoftCommonSecurityLog:
			output.MicrosoftCommonSecurityLogProcessor = append(output.MicrosoftCommonSecurityLogProcessor, PipelineMicrosoftCommonSecurityLogProcessorModel{
				Name: v.Name,
			})
		case pipelinegroups.ProcessorTypeMicrosoftSyslog:
			output.MicrosoftSyslogProcessor = append(output.MicrosoftSyslogProcessor, PipelineMicrosoftSyslogProcessorModel{
				Name: v.Name,
			})
		case pipelinegroups.ProcessorTypeTransformLanguage:
			processor := PipelineTransformLanguageProcessorModel{
				Name: v.Name,
			}
			if v.TransformLanguage != nil {
				processor.TransformStatement = v.TransformLanguage.TransformStatement
			}
			output.TransformLanguageProcessor = append(output.TransformLanguageProcessor, processor)
		}
	}
}

func flattenPipelineReceivers(input []pipelinegroups.Receiver, output *MonitorPipelineResourceModel) {
	for _, v := range input {
		switch v.Type {
		case pipelinegroups.ReceiverTypeOTLP:
			receiver := PipelineOtlpReceiverModel{
				Name:                 v.Name,
				TlsConfigurationName: pointer.From(v.TlsConfiguration),
			}
			if v.Otlp != nil {
				receiver.Endpoint = v.Otlp.Endpoint
			}
			output.OtlpReceiver = append(output.OtlpReceiver, receiver)
		case pipelinegroups.ReceiverTypeSyslog:
			receiver := PipelineSyslogReceiverModel{
				Name:                 v.Name,
				TlsConfigurationName: pointer.From(v.TlsConfiguration),
			}
			if v.Syslog != nil {
				receiver.Endpoint = v.Syslog.Endpoint
				receiver.AllowSkipPriorityHeader = pointer.From(v.Syslog.AllowSkipPriHeader)
				receiver.AllowedFormats = pointer.FromEnumSlice(v.Syslog.AllowedFormats)
				receiver.TransportProtocol = pointer.FromEnum(v.Syslog.TransportProtocol)
			}
			output.SyslogReceiver = append(output.SyslogReceiver, receiver)
		}
	}
}

func flattenPipelineService(input pipelinegroups.Service, output *MonitorPipelineResourceModel) {
	output.PersistentVolumeName = pointer.From(input.Persistence).PersistentVolumeName
	output.LogPipeline = make([]PipelineLogPipelineModel, 0, len(input.Pipelines))
	for _, p := range input.Pipelines {
		if p.Type != pipelinegroups.PipelineTypeLogs {
			continue
		}
		output.LogPipeline = append(output.LogPipeline, PipelineLogPipelineModel{
			Name:       p.Name,
			Exporters:  p.Exporters,
			Receivers:  p.Receivers,
			Processors: pointer.From(p.Processors),
		})
	}
}

func flattenPipelineTlsConfigurations(input *[]pipelinegroups.TlsConfiguration) []PipelineTlsConfigurationModel {
	if input == nil {
		return []PipelineTlsConfigurationModel{}
	}

	output := make([]PipelineTlsConfigurationModel, 0, len(*input))

	for _, v := range *input {
		tlsConfiguration := PipelineTlsConfigurationModel{
			Name: v.Name,
			Mode: pointer.FromEnum(v.Mode),
		}

		tlsConfiguration.ClientCertificateAuthority = flattenPipelineCertificateSource(v.ClientCa)

		if v.TlsCertificate != nil {
			tlsConfiguration.TlsCertificate = []PipelineTlsCertificateModel{
				{
					Certificate: flattenPipelineCertificateSource(&v.TlsCertificate.Certificate),
					PrivateKeySource: []PipelinePrivateKeySourceModel{
						{
							Location:    v.TlsCertificate.PrivateKey.Location,
							SubLocation: v.TlsCertificate.PrivateKey.SubLocation,
						},
					},
				},
			}
		}

		output = append(output, tlsConfiguration)
	}

	return output
}

func flattenPipelineCertificateSource(input *pipelinegroups.CertificateSource) []PipelineCertificateSourceModel {
	if input == nil {
		return []PipelineCertificateSourceModel{}
	}

	return []PipelineCertificateSourceModel{
		{
			Location:    input.Location,
			SubLocation: input.SubLocation,
			Type:        string(input.Type),
		},
	}
}
