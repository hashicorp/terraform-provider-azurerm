package workspace

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceProperties struct {
	ApplicationGroupReferences *[]string            `json:"applicationGroupReferences,omitempty"`
	CloudPcResource            *bool                `json:"cloudPcResource,omitempty"`
	Description                *string              `json:"description,omitempty"`
	FriendlyName               *string              `json:"friendlyName,omitempty"`
	ObjectId                   *string              `json:"objectId,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess `json:"publicNetworkAccess,omitempty"`
}
