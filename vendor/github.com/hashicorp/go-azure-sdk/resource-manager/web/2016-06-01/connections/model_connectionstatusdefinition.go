package connections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionStatusDefinition struct {
	Error  *ConnectionError `json:"error,omitempty"`
	Status *string          `json:"status,omitempty"`
	Target *string          `json:"target,omitempty"`
}
