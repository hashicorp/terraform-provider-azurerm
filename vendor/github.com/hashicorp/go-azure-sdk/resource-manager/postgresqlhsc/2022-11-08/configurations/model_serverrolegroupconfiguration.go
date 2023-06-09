package configurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerRoleGroupConfiguration struct {
	DefaultValue *string    `json:"defaultValue,omitempty"`
	Role         ServerRole `json:"role"`
	Source       *string    `json:"source,omitempty"`
	Value        string     `json:"value"`
}
