package arcsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArcIdentityResponseProperties struct {
	ArcApplicationClientId      *string `json:"arcApplicationClientId,omitempty"`
	ArcApplicationObjectId      *string `json:"arcApplicationObjectId,omitempty"`
	ArcApplicationTenantId      *string `json:"arcApplicationTenantId,omitempty"`
	ArcServicePrincipalObjectId *string `json:"arcServicePrincipalObjectId,omitempty"`
}
