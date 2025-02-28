package databases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseInviteFollowerRequest struct {
	InviteeEmail                string                       `json:"inviteeEmail"`
	TableLevelSharingProperties *TableLevelSharingProperties `json:"tableLevelSharingProperties,omitempty"`
}
