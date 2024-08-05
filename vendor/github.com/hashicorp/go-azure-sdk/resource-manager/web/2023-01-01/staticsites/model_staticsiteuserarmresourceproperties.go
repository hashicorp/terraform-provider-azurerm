package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteUserARMResourceProperties struct {
	DisplayName *string `json:"displayName,omitempty"`
	Provider    *string `json:"provider,omitempty"`
	Roles       *string `json:"roles,omitempty"`
	UserId      *string `json:"userId,omitempty"`
}
