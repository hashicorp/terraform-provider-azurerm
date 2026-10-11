package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomerManagedKeyDetails struct {
	KekIdentity *KekIdentityProperties `json:"kekIdentity,omitempty"`
	Key         *WorkspaceKeyDetails   `json:"key,omitempty"`
	Status      *string                `json:"status,omitempty"`
}
