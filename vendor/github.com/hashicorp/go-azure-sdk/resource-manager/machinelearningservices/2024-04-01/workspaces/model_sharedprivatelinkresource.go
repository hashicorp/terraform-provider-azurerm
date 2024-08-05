package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharedPrivateLinkResource struct {
	Name       *string                            `json:"name,omitempty"`
	Properties *SharedPrivateLinkResourceProperty `json:"properties,omitempty"`
}
