package jobdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobDefinitionUpdateProperties struct {
	AgentName   *string   `json:"agentName,omitempty"`
	CopyMode    *CopyMode `json:"copyMode,omitempty"`
	Description *string   `json:"description,omitempty"`
}
