package querypacks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogAnalyticsQueryPack struct {
	Id         *string                         `json:"id,omitempty"`
	Location   string                          `json:"location"`
	Name       *string                         `json:"name,omitempty"`
	Properties LogAnalyticsQueryPackProperties `json:"properties"`
	Tags       *map[string]string              `json:"tags,omitempty"`
	Type       *string                         `json:"type,omitempty"`
}
