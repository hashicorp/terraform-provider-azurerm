package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtensionDataSource struct {
	ExtensionName     string                             `json:"extensionName"`
	ExtensionSettings *interface{}                       `json:"extensionSettings,omitempty"`
	InputDataSources  *[]string                          `json:"inputDataSources,omitempty"`
	Name              *string                            `json:"name,omitempty"`
	Streams           *[]KnownExtensionDataSourceStreams `json:"streams,omitempty"`
}
