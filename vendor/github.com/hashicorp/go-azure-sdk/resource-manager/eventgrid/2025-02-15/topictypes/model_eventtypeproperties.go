package topictypes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventTypeProperties struct {
	Description    *string `json:"description,omitempty"`
	DisplayName    *string `json:"displayName,omitempty"`
	IsInDefaultSet *bool   `json:"isInDefaultSet,omitempty"`
	SchemaURL      *string `json:"schemaUrl,omitempty"`
}
