package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserRecord struct {
	AuthType *string         `json:"auth_type,omitempty"`
	Email    *string         `json:"email,omitempty"`
	FullName *string         `json:"full_name,omitempty"`
	Id       *string         `json:"id,omitempty"`
	Kind     *string         `json:"kind,omitempty"`
	Metadata *MetadataEntity `json:"metadata,omitempty"`
}
