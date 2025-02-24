package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HDInsightSparkActivityTypeProperties struct {
	Arguments             *[]string                         `json:"arguments,omitempty"`
	ClassName             *string                           `json:"className,omitempty"`
	EntryFilePath         string                            `json:"entryFilePath"`
	GetDebugInfo          *HDInsightActivityDebugInfoOption `json:"getDebugInfo,omitempty"`
	ProxyUser             *string                           `json:"proxyUser,omitempty"`
	RootPath              string                            `json:"rootPath"`
	SparkConfig           *map[string]string                `json:"sparkConfig,omitempty"`
	SparkJobLinkedService *LinkedServiceReference           `json:"sparkJobLinkedService,omitempty"`
}
