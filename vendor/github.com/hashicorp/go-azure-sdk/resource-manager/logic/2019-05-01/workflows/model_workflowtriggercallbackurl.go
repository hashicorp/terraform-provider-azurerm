package workflows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowTriggerCallbackURL struct {
	BasePath               *string                                `json:"basePath,omitempty"`
	Method                 *string                                `json:"method,omitempty"`
	Queries                *WorkflowTriggerListCallbackURLQueries `json:"queries,omitempty"`
	RelativePath           *string                                `json:"relativePath,omitempty"`
	RelativePathParameters *[]string                              `json:"relativePathParameters,omitempty"`
	Value                  *string                                `json:"value,omitempty"`
}
