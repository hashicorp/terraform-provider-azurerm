package managedhsmkeys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedHsmKey struct {
	Id         *string                 `json:"id,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties ManagedHsmKeyProperties `json:"properties"`
	Tags       *map[string]string      `json:"tags,omitempty"`
	Type       *string                 `json:"type,omitempty"`
}
