package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RuntimeScriptAction struct {
	ApplicationName *string  `json:"applicationName,omitempty"`
	Name            string   `json:"name"`
	Parameters      *string  `json:"parameters,omitempty"`
	Roles           []string `json:"roles"`
	Uri             string   `json:"uri"`
}
