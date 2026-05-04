package privatelinkscopes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkScopeValidationDetails struct {
	ConnectionDetails   *[]ConnectionDetail      `json:"connectionDetails,omitempty"`
	Id                  *string                  `json:"id,omitempty"`
	PublicNetworkAccess *PublicNetworkAccessType `json:"publicNetworkAccess,omitempty"`
}
