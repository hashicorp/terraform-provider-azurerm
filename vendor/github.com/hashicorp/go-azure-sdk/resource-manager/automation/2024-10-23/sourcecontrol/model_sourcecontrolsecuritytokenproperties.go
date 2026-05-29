package sourcecontrol

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceControlSecurityTokenProperties struct {
	AccessToken  *string    `json:"accessToken,omitempty"`
	RefreshToken *string    `json:"refreshToken,omitempty"`
	TokenType    *TokenType `json:"tokenType,omitempty"`
}
