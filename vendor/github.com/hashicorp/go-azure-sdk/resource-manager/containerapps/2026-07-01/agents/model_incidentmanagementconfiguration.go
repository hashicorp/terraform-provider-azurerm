package agents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IncidentManagementConfiguration struct {
	ConnectionKey  *string `json:"connectionKey,omitempty"`
	ConnectionName *string `json:"connectionName,omitempty"`
	ConnectionURL  *string `json:"connectionUrl,omitempty"`
	OboUser        *string `json:"oboUser,omitempty"`
	Type           *string `json:"type,omitempty"`
}
