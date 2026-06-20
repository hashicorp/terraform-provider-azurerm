package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogStorageSettings struct {
	EnableReliableLogging *interface{}           `json:"enableReliableLogging,omitempty"`
	LinkedServiceName     LinkedServiceReference `json:"linkedServiceName"`
	LogLevel              *interface{}           `json:"logLevel,omitempty"`
	Path                  *interface{}           `json:"path,omitempty"`
}
