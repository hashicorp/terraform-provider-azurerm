package tenantaccess

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessInformationSecretsContract struct {
	Enabled      *bool   `json:"enabled,omitempty"`
	Id           *string `json:"id,omitempty"`
	PrimaryKey   *string `json:"primaryKey,omitempty"`
	PrincipalId  *string `json:"principalId,omitempty"`
	SecondaryKey *string `json:"secondaryKey,omitempty"`
}
