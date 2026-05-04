package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevToolPortalSsoProperties struct {
	ClientId     *string   `json:"clientId,omitempty"`
	ClientSecret *string   `json:"clientSecret,omitempty"`
	MetadataURL  *string   `json:"metadataUrl,omitempty"`
	Scopes       *[]string `json:"scopes,omitempty"`
}
