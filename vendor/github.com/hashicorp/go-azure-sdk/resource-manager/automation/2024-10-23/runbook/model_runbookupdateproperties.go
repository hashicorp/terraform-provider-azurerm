package runbook

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunbookUpdateProperties struct {
	Description      *string `json:"description,omitempty"`
	LogActivityTrace *int64  `json:"logActivityTrace,omitempty"`
	LogProgress      *bool   `json:"logProgress,omitempty"`
	LogVerbose       *bool   `json:"logVerbose,omitempty"`
}
