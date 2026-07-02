package onlineendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointAuthToken struct {
	AccessToken         *string `json:"accessToken,omitempty"`
	ExpiryTimeUtc       *int64  `json:"expiryTimeUtc,omitempty"`
	RefreshAfterTimeUtc *int64  `json:"refreshAfterTimeUtc,omitempty"`
	TokenType           *string `json:"tokenType,omitempty"`
}
