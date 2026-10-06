package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReportableException struct {
	ActionableMessage *string `json:"actionableMessage,omitempty"`
	FilePath          *string `json:"filePath,omitempty"`
	HResult           *int64  `json:"hResult,omitempty"`
	LineNumber        *string `json:"lineNumber,omitempty"`
	Message           *string `json:"message,omitempty"`
	StackTrace        *string `json:"stackTrace,omitempty"`
}
