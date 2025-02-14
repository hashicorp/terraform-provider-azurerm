package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataLakeAnalyticsUSQLActivityTypeProperties struct {
	CompilationMode     *string                `json:"compilationMode,omitempty"`
	DegreeOfParallelism *int64                 `json:"degreeOfParallelism,omitempty"`
	Parameters          *map[string]string     `json:"parameters,omitempty"`
	Priority            *int64                 `json:"priority,omitempty"`
	RuntimeVersion      *string                `json:"runtimeVersion,omitempty"`
	ScriptLinkedService LinkedServiceReference `json:"scriptLinkedService"`
	ScriptPath          string                 `json:"scriptPath"`
}
