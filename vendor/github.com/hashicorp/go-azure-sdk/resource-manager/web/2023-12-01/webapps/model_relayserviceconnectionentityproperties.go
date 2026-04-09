package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RelayServiceConnectionEntityProperties struct {
	BiztalkUri               *string `json:"biztalkUri,omitempty"`
	EntityConnectionString   *string `json:"entityConnectionString,omitempty"`
	EntityName               *string `json:"entityName,omitempty"`
	Hostname                 *string `json:"hostname,omitempty"`
	Port                     *int64  `json:"port,omitempty"`
	ResourceConnectionString *string `json:"resourceConnectionString,omitempty"`
	ResourceType             *string `json:"resourceType,omitempty"`
}
