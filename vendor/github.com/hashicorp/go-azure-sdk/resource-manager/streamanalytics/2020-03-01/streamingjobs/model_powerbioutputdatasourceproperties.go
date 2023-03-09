package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PowerBIOutputDataSourceProperties struct {
	AuthenticationMode     *AuthenticationMode `json:"authenticationMode,omitempty"`
	Dataset                *string             `json:"dataset,omitempty"`
	GroupId                *string             `json:"groupId,omitempty"`
	GroupName              *string             `json:"groupName,omitempty"`
	RefreshToken           *string             `json:"refreshToken,omitempty"`
	Table                  *string             `json:"table,omitempty"`
	TokenUserDisplayName   *string             `json:"tokenUserDisplayName,omitempty"`
	TokenUserPrincipalName *string             `json:"tokenUserPrincipalName,omitempty"`
}
