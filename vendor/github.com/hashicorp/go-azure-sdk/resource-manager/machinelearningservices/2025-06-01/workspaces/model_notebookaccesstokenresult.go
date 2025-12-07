package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotebookAccessTokenResult struct {
	AccessToken        *string `json:"accessToken,omitempty"`
	ExpiresIn          *int64  `json:"expiresIn,omitempty"`
	HostName           *string `json:"hostName,omitempty"`
	NotebookResourceId *string `json:"notebookResourceId,omitempty"`
	PublicDns          *string `json:"publicDns,omitempty"`
	RefreshToken       *string `json:"refreshToken,omitempty"`
	Scope              *string `json:"scope,omitempty"`
	TokenType          *string `json:"tokenType,omitempty"`
}
