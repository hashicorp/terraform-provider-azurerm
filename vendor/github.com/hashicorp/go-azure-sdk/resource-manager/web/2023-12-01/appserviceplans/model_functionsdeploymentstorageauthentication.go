package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionsDeploymentStorageAuthentication struct {
	StorageAccountConnectionStringName *string             `json:"storageAccountConnectionStringName,omitempty"`
	Type                               *AuthenticationType `json:"type,omitempty"`
	UserAssignedIdentityResourceId     *string             `json:"userAssignedIdentityResourceId,omitempty"`
}
