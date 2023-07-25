package workflowrunactions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowRunActionRepetitionDefinitionCollection struct {
	NextLink *string                                  `json:"nextLink,omitempty"`
	Value    *[]WorkflowRunActionRepetitionDefinition `json:"value,omitempty"`
}
