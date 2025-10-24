package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessInviteUserAccountModel struct {
	Email              *string                   `json:"email,omitempty"`
	InvitedUserDetails *AccessInvitedUserDetails `json:"invitedUserDetails,omitempty"`
	OrganizationId     *string                   `json:"organizationId,omitempty"`
	Upn                *string                   `json:"upn,omitempty"`
}
