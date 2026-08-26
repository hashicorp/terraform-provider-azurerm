package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HDInsightPigActivityTypeProperties struct {
	Arguments             *interface{}                      `json:"arguments,omitempty"`
	Defines               *map[string]interface{}           `json:"defines,omitempty"`
	GetDebugInfo          *HDInsightActivityDebugInfoOption `json:"getDebugInfo,omitempty"`
	ScriptLinkedService   *LinkedServiceReference           `json:"scriptLinkedService,omitempty"`
	ScriptPath            *interface{}                      `json:"scriptPath,omitempty"`
	StorageLinkedServices *[]LinkedServiceReference         `json:"storageLinkedServices,omitempty"`
}
