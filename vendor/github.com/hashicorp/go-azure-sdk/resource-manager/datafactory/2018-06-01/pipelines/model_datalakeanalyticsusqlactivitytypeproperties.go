package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataLakeAnalyticsUSQLActivityTypeProperties struct {
	CompilationMode     *interface{}            `json:"compilationMode,omitempty"`
	DegreeOfParallelism *int64                  `json:"degreeOfParallelism,omitempty"`
	Parameters          *map[string]interface{} `json:"parameters,omitempty"`
	Priority            *int64                  `json:"priority,omitempty"`
	RuntimeVersion      *interface{}            `json:"runtimeVersion,omitempty"`
	ScriptLinkedService LinkedServiceReference  `json:"scriptLinkedService"`
	ScriptPath          interface{}             `json:"scriptPath"`
}
