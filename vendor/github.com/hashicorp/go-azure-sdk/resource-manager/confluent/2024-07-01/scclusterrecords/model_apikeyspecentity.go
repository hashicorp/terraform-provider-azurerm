package scclusterrecords

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type APIKeySpecEntity struct {
	Description *string               `json:"description,omitempty"`
	Name        *string               `json:"name,omitempty"`
	Owner       *APIKeyOwnerEntity    `json:"owner,omitempty"`
	Resource    *APIKeyResourceEntity `json:"resource,omitempty"`
	Secret      *string               `json:"secret,omitempty"`
}
