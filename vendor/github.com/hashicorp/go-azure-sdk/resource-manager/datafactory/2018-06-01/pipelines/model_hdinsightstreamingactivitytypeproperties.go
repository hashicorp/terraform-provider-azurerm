package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HDInsightStreamingActivityTypeProperties struct {
	Arguments             *[]interface{}                    `json:"arguments,omitempty"`
	Combiner              *interface{}                      `json:"combiner,omitempty"`
	CommandEnvironment    *[]interface{}                    `json:"commandEnvironment,omitempty"`
	Defines               *map[string]interface{}           `json:"defines,omitempty"`
	FileLinkedService     *LinkedServiceReference           `json:"fileLinkedService,omitempty"`
	FilePaths             []interface{}                     `json:"filePaths"`
	GetDebugInfo          *HDInsightActivityDebugInfoOption `json:"getDebugInfo,omitempty"`
	Input                 interface{}                       `json:"input"`
	Mapper                interface{}                       `json:"mapper"`
	Output                interface{}                       `json:"output"`
	Reducer               interface{}                       `json:"reducer"`
	StorageLinkedServices *[]LinkedServiceReference         `json:"storageLinkedServices,omitempty"`
}
