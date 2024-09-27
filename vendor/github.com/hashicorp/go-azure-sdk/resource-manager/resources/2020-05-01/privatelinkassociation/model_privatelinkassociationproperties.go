package privatelinkassociation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkAssociationProperties struct {
	PrivateLink         *string                     `json:"privateLink,omitempty"`
	PublicNetworkAccess *PublicNetworkAccessOptions `json:"publicNetworkAccess,omitempty"`
}
