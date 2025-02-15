package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HDInsightHiveActivityTypeProperties struct {
	Arguments             *[]string                         `json:"arguments,omitempty"`
	Defines               *map[string]string                `json:"defines,omitempty"`
	GetDebugInfo          *HDInsightActivityDebugInfoOption `json:"getDebugInfo,omitempty"`
	QueryTimeout          *int64                            `json:"queryTimeout,omitempty"`
	ScriptLinkedService   *LinkedServiceReference           `json:"scriptLinkedService,omitempty"`
	ScriptPath            *string                           `json:"scriptPath,omitempty"`
	StorageLinkedServices *[]LinkedServiceReference         `json:"storageLinkedServices,omitempty"`
	Variables             *map[string]string                `json:"variables,omitempty"`
}
