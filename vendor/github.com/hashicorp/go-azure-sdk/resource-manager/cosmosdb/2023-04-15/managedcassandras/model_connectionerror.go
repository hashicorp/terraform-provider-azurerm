package managedcassandras

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionError struct {
	ConnectionState *ConnectionState `json:"connectionState,omitempty"`
	Exception       *string          `json:"exception,omitempty"`
	IPFrom          *string          `json:"iPFrom,omitempty"`
	IPTo            *string          `json:"iPTo,omitempty"`
	Port            *int64           `json:"port,omitempty"`
}
