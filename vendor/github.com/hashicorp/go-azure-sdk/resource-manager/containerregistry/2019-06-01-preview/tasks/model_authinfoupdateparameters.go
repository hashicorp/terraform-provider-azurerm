package tasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthInfoUpdateParameters struct {
	ExpiresIn    *int64     `json:"expiresIn,omitempty"`
	RefreshToken *string    `json:"refreshToken,omitempty"`
	Scope        *string    `json:"scope,omitempty"`
	Token        *string    `json:"token,omitempty"`
	TokenType    *TokenType `json:"tokenType,omitempty"`
}
