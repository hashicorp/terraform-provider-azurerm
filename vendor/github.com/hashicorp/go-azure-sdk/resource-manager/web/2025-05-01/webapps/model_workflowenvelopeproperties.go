package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowEnvelopeProperties struct {
	Files     *map[string]interface{} `json:"files,omitempty"`
	FlowState *WorkflowState          `json:"flowState,omitempty"`
	Health    *WorkflowHealth         `json:"health,omitempty"`
}
