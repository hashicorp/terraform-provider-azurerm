package operations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RawGraphicalRunbookContent struct {
	RunbookDefinition *string           `json:"runbookDefinition,omitempty"`
	RunbookType       *GraphRunbookType `json:"runbookType,omitempty"`
	SchemaVersion     *string           `json:"schemaVersion,omitempty"`
}
