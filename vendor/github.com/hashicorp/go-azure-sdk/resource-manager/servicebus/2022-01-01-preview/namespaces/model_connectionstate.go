package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionState struct {
	Description *string                      `json:"description,omitempty"`
	Status      *PrivateLinkConnectionStatus `json:"status,omitempty"`
}
