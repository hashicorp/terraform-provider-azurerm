// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_pipeline

const (
	defaultPipelineBatchSize             = 8192
	defaultPipelineTimeoutInMilliseconds = 200
)

func (r Resource) ModelObject() any {
	return &MonitorPipelineResourceModel{}
}

type MonitorPipelineResourceModel struct {
	Name                                string                                             `tfschema:"name"`
	ResourceGroupName                   string                                             `tfschema:"resource_group_name"`
	Location                            string                                             `tfschema:"location"`
	CustomLocationId                    string                                             `tfschema:"custom_location_id"`
	LogPipeline                         []PipelineLogPipelineModel                         `tfschema:"log_pipeline"`
	AzureMonitorWorkspaceLogExporter    []PipelineAzureMonitorWorkspaceLogExporterModel    `tfschema:"azure_monitor_workspace_log_exporter"`
	BatchProcessor                      []PipelineBatchProcessorModel                      `tfschema:"batch_processor"`
	ExecutionPlacementConstraint        []PipelinePlacementConstraintModel                 `tfschema:"execution_placement_constraint"`
	MicrosoftCommonSecurityLogProcessor []PipelineMicrosoftCommonSecurityLogProcessorModel `tfschema:"microsoft_common_security_log_processor"`
	MicrosoftSyslogProcessor            []PipelineMicrosoftSyslogProcessorModel            `tfschema:"microsoft_syslog_processor"`
	OtlpReceiver                        []PipelineOtlpReceiverModel                        `tfschema:"otlp_receiver"`
	PersistentVolumeName                string                                             `tfschema:"persistent_volume_name"`
	ReplicaCount                        int64                                              `tfschema:"replica_count"`
	SyslogReceiver                      []PipelineSyslogReceiverModel                      `tfschema:"syslog_receiver"`
	TlsConfiguration                    []PipelineTlsConfigurationModel                    `tfschema:"tls_configuration"`
	TransformLanguageProcessor          []PipelineTransformLanguageProcessorModel          `tfschema:"transform_language_processor"`
	Tags                                map[string]string                                  `tfschema:"tags"`
}

type PipelineLogPipelineModel struct {
	Name       string   `tfschema:"name"`
	Exporters  []string `tfschema:"exporters"`
	Receivers  []string `tfschema:"receivers"`
	Processors []string `tfschema:"processors"`
}

type PipelineAzureMonitorWorkspaceLogExporterModel struct {
	Name                                string                                      `tfschema:"name"`
	Api                                 []PipelineAzureMonitorWorkspaceLogsApiModel `tfschema:"api"`
	PersistenceMaximumStorageUsageInGB  int64                                       `tfschema:"persistence_maximum_storage_usage_in_gb"`
	PersistenceRetentionPeriodInMinutes int64                                       `tfschema:"persistence_retention_period_in_minutes"`
}

type PipelineBatchProcessorModel struct {
	Name                  string `tfschema:"name"`
	BatchSize             int64  `tfschema:"batch_size"`
	TimeoutInMilliseconds int64  `tfschema:"timeout_in_milliseconds"`
}

type PipelinePlacementConstraintModel struct {
	Capability string   `tfschema:"capability"`
	Operator   string   `tfschema:"operator"`
	Values     []string `tfschema:"values"`
}

type PipelineMicrosoftCommonSecurityLogProcessorModel struct {
	Name string `tfschema:"name"`
}

type PipelineMicrosoftSyslogProcessorModel struct {
	Name string `tfschema:"name"`
}

type PipelineOtlpReceiverModel struct {
	Name                 string `tfschema:"name"`
	Endpoint             string `tfschema:"endpoint"`
	TlsConfigurationName string `tfschema:"tls_configuration_name"`
}

type PipelineSyslogReceiverModel struct {
	Name                    string   `tfschema:"name"`
	Endpoint                string   `tfschema:"endpoint"`
	AllowSkipPriorityHeader bool     `tfschema:"allow_skip_priority_header"`
	AllowedFormats          []string `tfschema:"allowed_formats"`
	TlsConfigurationName    string   `tfschema:"tls_configuration_name"`
	TransportProtocol       string   `tfschema:"transport_protocol"`
}

type PipelineTlsConfigurationModel struct {
	Name                       string                           `tfschema:"name"`
	ClientCertificateAuthority []PipelineCertificateSourceModel `tfschema:"client_certificate_authority"`
	Mode                       string                           `tfschema:"mode"`
	TlsCertificate             []PipelineTlsCertificateModel    `tfschema:"tls_certificate"`
}

type PipelineTransformLanguageProcessorModel struct {
	Name               string `tfschema:"name"`
	TransformStatement string `tfschema:"transform_statement"`
}

type PipelineAzureMonitorWorkspaceLogsApiModel struct {
	DataCollectionEndpointUrl     string                   `tfschema:"data_collection_endpoint_url"`
	DataCollectionRuleImmutableId string                   `tfschema:"data_collection_rule_immutable_id"`
	Schema                        []PipelineSchemaMapModel `tfschema:"schema"`
	Stream                        string                   `tfschema:"stream"`
}

type PipelineCertificateSourceModel struct {
	Location    string `tfschema:"location"`
	SubLocation string `tfschema:"sub_location"`
	Type        string `tfschema:"type"`
}

type PipelineTlsCertificateModel struct {
	Certificate      []PipelineCertificateSourceModel `tfschema:"certificate"`
	PrivateKeySource []PipelinePrivateKeySourceModel  `tfschema:"private_key_source"`
}

type PipelineSchemaMapModel struct {
	RecordMap   []PipelineFieldMapModel `tfschema:"record_map"`
	ResourceMap []PipelineFieldMapModel `tfschema:"resource_map"`
	ScopeMap    []PipelineFieldMapModel `tfschema:"scope_map"`
}

type PipelinePrivateKeySourceModel struct {
	Location    string `tfschema:"location"`
	SubLocation string `tfschema:"sub_location"`
}

type PipelineFieldMapModel struct {
	From string `tfschema:"from"`
	To   string `tfschema:"to"`
}
