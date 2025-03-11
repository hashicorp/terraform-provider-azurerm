package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExecutePowerQueryActivityTypeProperties struct {
	Compute                  *ExecuteDataFlowActivityTypePropertiesCompute `json:"compute,omitempty"`
	ContinuationSettings     *ContinuationSettingsReference                `json:"continuationSettings,omitempty"`
	ContinueOnError          *bool                                         `json:"continueOnError,omitempty"`
	DataFlow                 DataFlowReference                             `json:"dataFlow"`
	IntegrationRuntime       *IntegrationRuntimeReference                  `json:"integrationRuntime,omitempty"`
	Queries                  *[]PowerQuerySinkMapping                      `json:"queries,omitempty"`
	RunConcurrently          *bool                                         `json:"runConcurrently,omitempty"`
	Sinks                    *map[string]PowerQuerySink                    `json:"sinks,omitempty"`
	SourceStagingConcurrency *int64                                        `json:"sourceStagingConcurrency,omitempty"`
	Staging                  *DataFlowStagingInfo                          `json:"staging,omitempty"`
	TraceLevel               *interface{}                                  `json:"traceLevel,omitempty"`
}
