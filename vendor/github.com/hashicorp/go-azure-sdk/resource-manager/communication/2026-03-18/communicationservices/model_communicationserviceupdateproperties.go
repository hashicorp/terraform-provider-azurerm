package communicationservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommunicationServiceUpdateProperties struct {
	DisableLocalAuth    *bool                `json:"disableLocalAuth,omitempty"`
	LinkedDomains       *[]string            `json:"linkedDomains,omitempty"`
	PublicNetworkAccess *PublicNetworkAccess `json:"publicNetworkAccess,omitempty"`
}
