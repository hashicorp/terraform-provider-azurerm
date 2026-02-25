package runbook

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunbookCreateOrUpdateParameters struct {
	Location   *string                         `json:"location,omitempty"`
	Name       *string                         `json:"name,omitempty"`
	Properties RunbookCreateOrUpdateProperties `json:"properties"`
	Tags       *map[string]string              `json:"tags,omitempty"`
}
