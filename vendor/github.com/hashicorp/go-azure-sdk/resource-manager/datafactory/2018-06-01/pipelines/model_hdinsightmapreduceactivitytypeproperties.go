package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HDInsightMapReduceActivityTypeProperties struct {
	Arguments             *[]string                         `json:"arguments,omitempty"`
	ClassName             string                            `json:"className"`
	Defines               *map[string]string                `json:"defines,omitempty"`
	GetDebugInfo          *HDInsightActivityDebugInfoOption `json:"getDebugInfo,omitempty"`
	JarFilePath           string                            `json:"jarFilePath"`
	JarLibs               *[]string                         `json:"jarLibs,omitempty"`
	JarLinkedService      *LinkedServiceReference           `json:"jarLinkedService,omitempty"`
	StorageLinkedServices *[]LinkedServiceReference         `json:"storageLinkedServices,omitempty"`
}
