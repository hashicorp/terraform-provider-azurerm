package webtestsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTest struct {
	Id         *string            `json:"id,omitempty"`
	Kind       *WebTestKind       `json:"kind,omitempty"`
	Location   string             `json:"location"`
	Name       *string            `json:"name,omitempty"`
	Properties *WebTestProperties `json:"properties,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
	Type       *string            `json:"type,omitempty"`
}
