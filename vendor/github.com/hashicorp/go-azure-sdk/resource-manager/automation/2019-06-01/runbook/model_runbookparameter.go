package runbook

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunbookParameter struct {
	DefaultValue *string `json:"defaultValue,omitempty"`
	IsMandatory  *bool   `json:"isMandatory,omitempty"`
	Position     *int64  `json:"position,omitempty"`
	Type         *string `json:"type,omitempty"`
}
