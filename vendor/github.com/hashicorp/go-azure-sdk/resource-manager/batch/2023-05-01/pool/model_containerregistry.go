package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerRegistry struct {
	IdentityReference *ComputeNodeIdentityReference `json:"identityReference,omitempty"`
	Password          *string                       `json:"password,omitempty"`
	RegistryServer    *string                       `json:"registryServer,omitempty"`
	Username          *string                       `json:"username,omitempty"`
}
