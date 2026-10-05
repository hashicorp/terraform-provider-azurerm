package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HDInsightMapReduceActivityTypeProperties struct {
	Arguments             *[]interface{}                    `json:"arguments,omitempty"`
	ClassName             interface{}                       `json:"className"`
	Defines               *map[string]interface{}           `json:"defines,omitempty"`
	GetDebugInfo          *HDInsightActivityDebugInfoOption `json:"getDebugInfo,omitempty"`
	JarFilePath           interface{}                       `json:"jarFilePath"`
	JarLibs               *[]interface{}                    `json:"jarLibs,omitempty"`
	JarLinkedService      *LinkedServiceReference           `json:"jarLinkedService,omitempty"`
	StorageLinkedServices *[]LinkedServiceReference         `json:"storageLinkedServices,omitempty"`
}
