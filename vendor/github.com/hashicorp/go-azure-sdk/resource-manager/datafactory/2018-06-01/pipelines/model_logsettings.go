package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogSettings struct {
	CopyActivityLogSettings *CopyActivityLogSettings `json:"copyActivityLogSettings,omitempty"`
	EnableCopyActivityLog   *bool                    `json:"enableCopyActivityLog,omitempty"`
	LogLocationSettings     LogLocationSettings      `json:"logLocationSettings"`
}
