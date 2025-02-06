package componentsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkScopedResource struct {
	ResourceId *string `json:"ResourceId,omitempty"`
	ScopeId    *string `json:"ScopeId,omitempty"`
}
