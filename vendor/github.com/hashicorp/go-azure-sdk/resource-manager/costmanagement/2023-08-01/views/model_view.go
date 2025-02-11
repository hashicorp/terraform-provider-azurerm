package views

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type View struct {
	ETag       *string         `json:"eTag,omitempty"`
	Id         *string         `json:"id,omitempty"`
	Name       *string         `json:"name,omitempty"`
	Properties *ViewProperties `json:"properties,omitempty"`
	Type       *string         `json:"type,omitempty"`
}
