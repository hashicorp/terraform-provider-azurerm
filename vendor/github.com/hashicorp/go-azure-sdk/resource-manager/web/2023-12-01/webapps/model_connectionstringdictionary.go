package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionStringDictionary struct {
	Id         *string                             `json:"id,omitempty"`
	Kind       *string                             `json:"kind,omitempty"`
	Name       *string                             `json:"name,omitempty"`
	Properties *map[string]ConnStringValueTypePair `json:"properties,omitempty"`
	Type       *string                             `json:"type,omitempty"`
}
