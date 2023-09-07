package sim

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SimPropertiesFormat struct {
	AuthenticationKey                     *string                           `json:"authenticationKey,omitempty"`
	DeviceType                            *string                           `json:"deviceType,omitempty"`
	IntegratedCircuitCardIdentifier       *string                           `json:"integratedCircuitCardIdentifier,omitempty"`
	InternationalMobileSubscriberIdentity string                            `json:"internationalMobileSubscriberIdentity"`
	OperatorKeyCode                       *string                           `json:"operatorKeyCode,omitempty"`
	ProvisioningState                     *ProvisioningState                `json:"provisioningState,omitempty"`
	SimPolicy                             *SimPolicyResourceId              `json:"simPolicy,omitempty"`
	SimState                              *SimState                         `json:"simState,omitempty"`
	SiteProvisioningState                 *map[string]SiteProvisioningState `json:"siteProvisioningState,omitempty"`
	StaticIPConfiguration                 *[]SimStaticIPProperties          `json:"staticIpConfiguration,omitempty"`
	VendorKeyFingerprint                  *string                           `json:"vendorKeyFingerprint,omitempty"`
	VendorName                            *string                           `json:"vendorName,omitempty"`
}
