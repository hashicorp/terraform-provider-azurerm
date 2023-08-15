package apimanagementservice

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiManagementServiceUpdateProperties struct {
	AdditionalLocations         *[]AdditionalLocation                     `json:"additionalLocations,omitempty"`
	ApiVersionConstraint        *ApiVersionConstraint                     `json:"apiVersionConstraint,omitempty"`
	Certificates                *[]CertificateConfiguration               `json:"certificates,omitempty"`
	CreatedAtUtc                *string                                   `json:"createdAtUtc,omitempty"`
	CustomProperties            *map[string]string                        `json:"customProperties,omitempty"`
	DeveloperPortalUrl          *string                                   `json:"developerPortalUrl,omitempty"`
	DisableGateway              *bool                                     `json:"disableGateway,omitempty"`
	EnableClientCertificate     *bool                                     `json:"enableClientCertificate,omitempty"`
	GatewayRegionalUrl          *string                                   `json:"gatewayRegionalUrl,omitempty"`
	GatewayUrl                  *string                                   `json:"gatewayUrl,omitempty"`
	HostnameConfigurations      *[]HostnameConfiguration                  `json:"hostnameConfigurations,omitempty"`
	ManagementApiUrl            *string                                   `json:"managementApiUrl,omitempty"`
	NotificationSenderEmail     *string                                   `json:"notificationSenderEmail,omitempty"`
	PlatformVersion             *PlatformVersion                          `json:"platformVersion,omitempty"`
	PortalUrl                   *string                                   `json:"portalUrl,omitempty"`
	PrivateEndpointConnections  *[]RemotePrivateEndpointConnectionWrapper `json:"privateEndpointConnections,omitempty"`
	PrivateIPAddresses          *[]string                                 `json:"privateIPAddresses,omitempty"`
	ProvisioningState           *string                                   `json:"provisioningState,omitempty"`
	PublicIPAddressId           *string                                   `json:"publicIpAddressId,omitempty"`
	PublicIPAddresses           *[]string                                 `json:"publicIPAddresses,omitempty"`
	PublicNetworkAccess         *PublicNetworkAccess                      `json:"publicNetworkAccess,omitempty"`
	PublisherEmail              *string                                   `json:"publisherEmail,omitempty"`
	PublisherName               *string                                   `json:"publisherName,omitempty"`
	Restore                     *bool                                     `json:"restore,omitempty"`
	ScmUrl                      *string                                   `json:"scmUrl,omitempty"`
	TargetProvisioningState     *string                                   `json:"targetProvisioningState,omitempty"`
	VirtualNetworkConfiguration *VirtualNetworkConfiguration              `json:"virtualNetworkConfiguration,omitempty"`
	VirtualNetworkType          *VirtualNetworkType                       `json:"virtualNetworkType,omitempty"`
}

func (o *ApiManagementServiceUpdateProperties) GetCreatedAtUtcAsTime() (*time.Time, error) {
	if o.CreatedAtUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAtUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ApiManagementServiceUpdateProperties) SetCreatedAtUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAtUtc = &formatted
}
