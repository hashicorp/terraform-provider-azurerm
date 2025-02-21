package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HDInsightStreamingActivityTypeProperties struct {
	Arguments             *[]string                         `json:"arguments,omitempty"`
	Combiner              *string                           `json:"combiner,omitempty"`
	CommandEnvironment    *[]string                         `json:"commandEnvironment,omitempty"`
	Defines               *map[string]string                `json:"defines,omitempty"`
	FileLinkedService     *LinkedServiceReference           `json:"fileLinkedService,omitempty"`
	FilePaths             []string                          `json:"filePaths"`
	GetDebugInfo          *HDInsightActivityDebugInfoOption `json:"getDebugInfo,omitempty"`
	Input                 string                            `json:"input"`
	Mapper                string                            `json:"mapper"`
	Output                string                            `json:"output"`
	Reducer               string                            `json:"reducer"`
	StorageLinkedServices *[]LinkedServiceReference         `json:"storageLinkedServices,omitempty"`
}
