package logprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogProfileResource struct {
	Id         *string              `json:"id,omitempty"`
	Location   string               `json:"location"`
	Name       *string              `json:"name,omitempty"`
	Properties LogProfileProperties `json:"properties"`
	Tags       *map[string]string   `json:"tags,omitempty"`
	Type       *string              `json:"type,omitempty"`
}
