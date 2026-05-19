package projectconnectionresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionManagedIdentity struct {
	ClientId   *string `json:"clientId,omitempty"`
	ResourceId *string `json:"resourceId,omitempty"`
}
