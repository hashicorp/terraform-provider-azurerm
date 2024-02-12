package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticParameters struct {
	AppInstance *string `json:"appInstance,omitempty"`
	Duration    *string `json:"duration,omitempty"`
	FilePath    *string `json:"filePath,omitempty"`
}
