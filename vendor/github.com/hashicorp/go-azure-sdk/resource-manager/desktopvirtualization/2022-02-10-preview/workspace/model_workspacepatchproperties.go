package workspace

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacePatchProperties struct {
	ApplicationGroupReferences *[]string            `json:"applicationGroupReferences,omitempty"`
	Description                *string              `json:"description,omitempty"`
	FriendlyName               *string              `json:"friendlyName,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess `json:"publicNetworkAccess,omitempty"`
}
