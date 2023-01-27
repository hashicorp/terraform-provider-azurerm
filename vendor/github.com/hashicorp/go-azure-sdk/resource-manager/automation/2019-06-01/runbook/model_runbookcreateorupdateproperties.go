package runbook

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunbookCreateOrUpdateProperties struct {
	Description        *string         `json:"description,omitempty"`
	Draft              *RunbookDraft   `json:"draft,omitempty"`
	LogActivityTrace   *int64          `json:"logActivityTrace,omitempty"`
	LogProgress        *bool           `json:"logProgress,omitempty"`
	LogVerbose         *bool           `json:"logVerbose,omitempty"`
	PublishContentLink *ContentLink    `json:"publishContentLink,omitempty"`
	RunbookType        RunbookTypeEnum `json:"runbookType"`
}
