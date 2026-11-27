// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_pipeline

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2026-04-01/pipelinegroups"
)

func expandPipelineExecutionPlacement(input []PipelinePlacementConstraintModel) *pipelinegroups.ExecutionPlacement {
	if len(input) == 0 {
		return nil
	}

	constraints := make([]pipelinegroups.PlacementConstraint, 0, len(input))
	for _, constraint := range input {
		constraints = append(constraints, pipelinegroups.PlacementConstraint{
			Capability: constraint.Capability,
			Operator:   pipelinegroups.CapabilityOperator(constraint.Operator),
			Values:     pointer.To(constraint.Values),
		})
	}

	return &pipelinegroups.ExecutionPlacement{
		Constraints: &constraints,
	}
}

func expandPipelineAzureMonitorWorkspaceLogExporters(input []PipelineAzureMonitorWorkspaceLogExporterModel) []pipelinegroups.Exporter {
	output := make([]pipelinegroups.Exporter, 0, len(input))

	for _, v := range input {
		output = append(output, pipelinegroups.Exporter{
			Name:                      v.Name,
			Type:                      pipelinegroups.ExporterTypeAzureMonitorWorkspaceLogs,
			AzureMonitorWorkspaceLogs: expandPipelineAzureMonitorWorkspaceLogExporter(v),
		})
	}

	return output
}

func expandPipelineAzureMonitorWorkspaceLogExporter(input PipelineAzureMonitorWorkspaceLogExporterModel) *pipelinegroups.AzureMonitorWorkspaceLogsExporter {
	output := pipelinegroups.AzureMonitorWorkspaceLogsExporter{}

	if len(input.Api) > 0 {
		api := input.Api[0]
		output.Api = pipelinegroups.AzureMonitorWorkspaceLogsApiConfig{
			DataCollectionEndpointURL: api.DataCollectionEndpointUrl,
			DataCollectionRule:        api.DataCollectionRuleImmutableId,
			Stream:                    api.Stream,
		}

		if len(api.Schema) > 0 {
			output.Api.Schema = expandPipelineSchemaMap(api.Schema[0])
		}
	}

	if input.PersistenceMaximumStorageUsageInGB != 0 || input.PersistenceRetentionPeriodInMinutes != 0 {
		output.Persistence = &pipelinegroups.ExporterPersistenceConfiguration{
			MaxStorageUsage: pointer.ToOrNil(input.PersistenceMaximumStorageUsageInGB),
			RetentionPeriod: pointer.ToOrNil(input.PersistenceRetentionPeriodInMinutes),
		}
	}

	return &output
}

func expandPipelineSchemaMap(input PipelineSchemaMapModel) pipelinegroups.SchemaMap {
	output := pipelinegroups.SchemaMap{
		RecordMap: expandPipelineFieldMaps(input.RecordMap),
	}

	if len(input.ResourceMap) > 0 {
		resourceMaps := make([]pipelinegroups.ResourceMap, 0, len(input.ResourceMap))
		for _, v := range input.ResourceMap {
			resourceMaps = append(resourceMaps, pipelinegroups.ResourceMap{
				From: v.From,
				To:   v.To,
			})
		}
		output.ResourceMap = &resourceMaps
	}

	if len(input.ScopeMap) > 0 {
		scopeMaps := make([]pipelinegroups.ScopeMap, 0, len(input.ScopeMap))
		for _, v := range input.ScopeMap {
			scopeMaps = append(scopeMaps, pipelinegroups.ScopeMap{
				From: v.From,
				To:   v.To,
			})
		}
		output.ScopeMap = &scopeMaps
	}

	return output
}

func expandPipelineFieldMaps(input []PipelineFieldMapModel) []pipelinegroups.RecordMap {
	output := make([]pipelinegroups.RecordMap, 0, len(input))
	for _, v := range input {
		output = append(output, pipelinegroups.RecordMap{
			From: v.From,
			To:   v.To,
		})
	}
	return output
}

func expandPipelineProcessors(input MonitorPipelineResourceModel) []pipelinegroups.Processor {
	output := make([]pipelinegroups.Processor, 0, len(input.BatchProcessor)+len(input.MicrosoftCommonSecurityLogProcessor)+len(input.MicrosoftSyslogProcessor)+len(input.TransformLanguageProcessor))

	for _, v := range input.BatchProcessor {
		output = append(output, pipelinegroups.Processor{
			Name: v.Name,
			Type: pipelinegroups.ProcessorTypeBatch,
			Batch: &pipelinegroups.BatchProcessor{
				BatchSize: pointer.To(v.BatchSize),
				Timeout:   pointer.To(v.TimeoutInMilliseconds),
			},
		})
	}

	for _, v := range input.MicrosoftCommonSecurityLogProcessor {
		output = append(output, pipelinegroups.Processor{
			Name: v.Name,
			Type: pipelinegroups.ProcessorTypeMicrosoftCommonSecurityLog,
		})
	}

	for _, v := range input.MicrosoftSyslogProcessor {
		output = append(output, pipelinegroups.Processor{
			Name: v.Name,
			Type: pipelinegroups.ProcessorTypeMicrosoftSyslog,
		})
	}

	for _, v := range input.TransformLanguageProcessor {
		output = append(output, pipelinegroups.Processor{
			Name: v.Name,
			Type: pipelinegroups.ProcessorTypeTransformLanguage,
			TransformLanguage: &pipelinegroups.TransformLanguageProcessor{
				TransformStatement: v.TransformStatement,
			},
		})
	}

	return output
}

func expandPipelineReceivers(input MonitorPipelineResourceModel) []pipelinegroups.Receiver {
	output := make([]pipelinegroups.Receiver, 0, len(input.OtlpReceiver)+len(input.SyslogReceiver))

	for _, v := range input.OtlpReceiver {
		output = append(output, pipelinegroups.Receiver{
			Name: v.Name,
			Type: pipelinegroups.ReceiverTypeOTLP,
			Otlp: &pipelinegroups.OtlpReceiver{
				Endpoint: v.Endpoint,
			},
			TlsConfiguration: pointer.ToOrNil(v.TlsConfigurationName),
		})
	}

	for _, v := range input.SyslogReceiver {
		syslog := &pipelinegroups.SyslogReceiver{
			Endpoint:           v.Endpoint,
			AllowSkipPriHeader: &v.AllowSkipPriorityHeader,
			TransportProtocol:  pointer.ToEnum[pipelinegroups.TransportProtocol](v.TransportProtocol),
		}
		if len(v.AllowedFormats) > 0 {
			syslog.AllowedFormats = pointer.ToEnumSlice[pipelinegroups.AllowedFormats](v.AllowedFormats)
		}

		output = append(output, pipelinegroups.Receiver{
			Name:             v.Name,
			Type:             pipelinegroups.ReceiverTypeSyslog,
			Syslog:           syslog,
			TlsConfiguration: pointer.ToOrNil(v.TlsConfigurationName),
		})
	}

	return output
}

func expandPipelineService(input []PipelineLogPipelineModel, persistentVolumeName string) pipelinegroups.Service {
	output := pipelinegroups.Service{
		Pipelines: []pipelinegroups.Pipeline{},
	}

	if persistentVolumeName != "" {
		output.Persistence = &pipelinegroups.PersistenceConfigurations{
			PersistentVolumeName: persistentVolumeName,
		}
	}

	pipelines := make([]pipelinegroups.Pipeline, 0, len(input))
	for _, p := range input {
		pipeline := pipelinegroups.Pipeline{
			Name:      p.Name,
			Type:      pipelinegroups.PipelineTypeLogs,
			Exporters: p.Exporters,
			Receivers: p.Receivers,
		}

		if len(p.Processors) > 0 {
			pipeline.Processors = &p.Processors
		}

		pipelines = append(pipelines, pipeline)
	}
	output.Pipelines = pipelines

	return output
}

func expandPipelineTlsConfigurations(input []PipelineTlsConfigurationModel) *[]pipelinegroups.TlsConfiguration {
	if len(input) == 0 {
		return pointer.To([]pipelinegroups.TlsConfiguration{})
	}

	output := make([]pipelinegroups.TlsConfiguration, 0, len(input))

	for _, v := range input {
		tlsConfiguration := pipelinegroups.TlsConfiguration{
			Name: v.Name,
			Mode: pointer.ToEnum[pipelinegroups.TlsMode](v.Mode),
		}

		if len(v.ClientCertificateAuthority) > 0 {
			tlsConfiguration.ClientCa = expandPipelineCertificateSource(v.ClientCertificateAuthority[0])
		}

		if len(v.TlsCertificate) > 0 {
			certificate := v.TlsCertificate[0]
			certificateWithKey := pipelinegroups.CertificateWithKey{}

			if len(certificate.Certificate) > 0 {
				certificateWithKey.Certificate = *expandPipelineCertificateSource(certificate.Certificate[0])
			}
			if len(certificate.PrivateKeySource) > 0 {
				certificateWithKey.PrivateKey = pipelinegroups.PrivateKeySource{
					Location:    certificate.PrivateKeySource[0].Location,
					SubLocation: certificate.PrivateKeySource[0].SubLocation,
					Type:        pipelinegroups.PrivateKeySourceTypeKubernetesSecret,
				}
			}

			tlsConfiguration.TlsCertificate = &certificateWithKey
		}

		output = append(output, tlsConfiguration)
	}

	return pointer.To(output)
}

func expandPipelineCertificateSource(input PipelineCertificateSourceModel) *pipelinegroups.CertificateSource {
	return &pipelinegroups.CertificateSource{
		Location:    input.Location,
		SubLocation: input.SubLocation,
		Type:        pipelinegroups.CertificateSourceType(input.Type),
	}
}
