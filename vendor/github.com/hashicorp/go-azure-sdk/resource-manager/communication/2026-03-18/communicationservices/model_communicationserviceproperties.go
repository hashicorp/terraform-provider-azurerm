package communicationservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommunicationServiceProperties struct {
	DataLocation        string                                  `json:"dataLocation"`
	DisableLocalAuth    *bool                                   `json:"disableLocalAuth,omitempty"`
	HostName            *string                                 `json:"hostName,omitempty"`
	ImmutableResourceId *string                                 `json:"immutableResourceId,omitempty"`
	LinkedDomains       *[]string                               `json:"linkedDomains,omitempty"`
	NotificationHubId   *string                                 `json:"notificationHubId,omitempty"`
	ProvisioningState   *CommunicationServicesProvisioningState `json:"provisioningState,omitempty"`
	PublicNetworkAccess *PublicNetworkAccess                    `json:"publicNetworkAccess,omitempty"`
	Version             *string                                 `json:"version,omitempty"`
}
