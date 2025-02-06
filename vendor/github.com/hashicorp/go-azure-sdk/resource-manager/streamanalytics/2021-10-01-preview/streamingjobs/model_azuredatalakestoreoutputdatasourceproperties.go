package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureDataLakeStoreOutputDataSourceProperties struct {
	AccountName            *string             `json:"accountName,omitempty"`
	AuthenticationMode     *AuthenticationMode `json:"authenticationMode,omitempty"`
	DateFormat             *string             `json:"dateFormat,omitempty"`
	FilePathPrefix         *string             `json:"filePathPrefix,omitempty"`
	RefreshToken           *string             `json:"refreshToken,omitempty"`
	TenantId               *string             `json:"tenantId,omitempty"`
	TimeFormat             *string             `json:"timeFormat,omitempty"`
	TokenUserDisplayName   *string             `json:"tokenUserDisplayName,omitempty"`
	TokenUserPrincipalName *string             `json:"tokenUserPrincipalName,omitempty"`
}
