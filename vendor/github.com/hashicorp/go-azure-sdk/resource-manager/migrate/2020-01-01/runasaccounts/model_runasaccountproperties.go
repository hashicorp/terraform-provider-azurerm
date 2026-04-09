package runasaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunAsAccountProperties struct {
	CreatedTimestamp *string         `json:"createdTimestamp,omitempty"`
	CredentialType   *CredentialType `json:"credentialType,omitempty"`
	DisplayName      *string         `json:"displayName,omitempty"`
	UpdatedTimestamp *string         `json:"updatedTimestamp,omitempty"`
}
