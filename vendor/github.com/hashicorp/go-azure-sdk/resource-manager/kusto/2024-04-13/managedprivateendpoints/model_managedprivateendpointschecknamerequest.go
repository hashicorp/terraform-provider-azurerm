package managedprivateendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedPrivateEndpointsCheckNameRequest struct {
	Name string                      `json:"name"`
	Type ManagedPrivateEndpointsType `json:"type"`
}
