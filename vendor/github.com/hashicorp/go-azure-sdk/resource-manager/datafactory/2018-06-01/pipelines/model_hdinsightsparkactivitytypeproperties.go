package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HDInsightSparkActivityTypeProperties struct {
	Arguments             *[]interface{}                    `json:"arguments,omitempty"`
	ClassName             *string                           `json:"className,omitempty"`
	EntryFilePath         interface{}                       `json:"entryFilePath"`
	GetDebugInfo          *HDInsightActivityDebugInfoOption `json:"getDebugInfo,omitempty"`
	ProxyUser             *interface{}                      `json:"proxyUser,omitempty"`
	RootPath              interface{}                       `json:"rootPath"`
	SparkConfig           *map[string]interface{}           `json:"sparkConfig,omitempty"`
	SparkJobLinkedService *LinkedServiceReference           `json:"sparkJobLinkedService,omitempty"`
}
