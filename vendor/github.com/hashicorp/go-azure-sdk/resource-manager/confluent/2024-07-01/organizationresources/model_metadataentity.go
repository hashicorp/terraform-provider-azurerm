package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetadataEntity struct {
	CreatedAt    *string `json:"created_at,omitempty"`
	DeletedAt    *string `json:"deleted_at,omitempty"`
	ResourceName *string `json:"resource_name,omitempty"`
	Self         *string `json:"self,omitempty"`
	UpdatedAt    *string `json:"updated_at,omitempty"`
}
