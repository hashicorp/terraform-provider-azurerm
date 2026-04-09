package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CopyActivityLogSettings struct {
	EnableReliableLogging *interface{} `json:"enableReliableLogging,omitempty"`
	LogLevel              *interface{} `json:"logLevel,omitempty"`
}
