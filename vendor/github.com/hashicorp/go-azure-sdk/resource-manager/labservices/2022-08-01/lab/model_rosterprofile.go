package lab

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RosterProfile struct {
	ActiveDirectoryGroupId *string `json:"activeDirectoryGroupId,omitempty"`
	LmsInstance            *string `json:"lmsInstance,omitempty"`
	LtiClientId            *string `json:"ltiClientId,omitempty"`
	LtiContextId           *string `json:"ltiContextId,omitempty"`
	LtiRosterEndpoint      *string `json:"ltiRosterEndpoint,omitempty"`
}
