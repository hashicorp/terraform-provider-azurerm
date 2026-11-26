package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExecuteDataFlowActivityTypeProperties struct {
	Compute                  *ExecuteDataFlowActivityTypePropertiesCompute `json:"compute,omitempty"`
	ContinuationSettings     *ContinuationSettingsReference                `json:"continuationSettings,omitempty"`
	ContinueOnError          *interface{}                                  `json:"continueOnError,omitempty"`
	DataFlow                 DataFlowReference                             `json:"dataFlow"`
	IntegrationRuntime       *IntegrationRuntimeReference                  `json:"integrationRuntime,omitempty"`
	RunConcurrently          *interface{}                                  `json:"runConcurrently,omitempty"`
	SourceStagingConcurrency *interface{}                                  `json:"sourceStagingConcurrency,omitempty"`
	Staging                  *DataFlowStagingInfo                          `json:"staging,omitempty"`
	TraceLevel               *interface{}                                  `json:"traceLevel,omitempty"`
}
