package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InvitationRecord struct {
	AcceptedAt *string         `json:"accepted_at,omitempty"`
	AuthType   *string         `json:"auth_type,omitempty"`
	Email      *string         `json:"email,omitempty"`
	ExpiresAt  *string         `json:"expires_at,omitempty"`
	Id         *string         `json:"id,omitempty"`
	Kind       *string         `json:"kind,omitempty"`
	Metadata   *MetadataEntity `json:"metadata,omitempty"`
	Status     *string         `json:"status,omitempty"`
}
