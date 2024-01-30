package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomHostnameSitesProperties struct {
	CustomHostname  *string       `json:"customHostname,omitempty"`
	Region          *string       `json:"region,omitempty"`
	SiteResourceIds *[]Identifier `json:"siteResourceIds,omitempty"`
}
