package dataflowendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowEndpointDataExplorerAuthentication struct {
	Method                                ManagedIdentityMethod                                        `json:"method"`
	SystemAssignedManagedIdentitySettings *DataflowEndpointAuthenticationSystemAssignedManagedIdentity `json:"systemAssignedManagedIdentitySettings,omitempty"`
	UserAssignedManagedIdentitySettings   *DataflowEndpointAuthenticationUserAssignedManagedIdentity   `json:"userAssignedManagedIdentitySettings,omitempty"`
}
