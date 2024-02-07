package customipprefixes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomIPPrefixPropertiesFormat struct {
	Asn                   *string             `json:"asn,omitempty"`
	AuthorizationMessage  *string             `json:"authorizationMessage,omitempty"`
	ChildCustomIPPrefixes *[]SubResource      `json:"childCustomIpPrefixes,omitempty"`
	Cidr                  *string             `json:"cidr,omitempty"`
	CommissionedState     *CommissionedState  `json:"commissionedState,omitempty"`
	CustomIPPrefixParent  *SubResource        `json:"customIpPrefixParent,omitempty"`
	ExpressRouteAdvertise *bool               `json:"expressRouteAdvertise,omitempty"`
	FailedReason          *string             `json:"failedReason,omitempty"`
	Geo                   *Geo                `json:"geo,omitempty"`
	NoInternetAdvertise   *bool               `json:"noInternetAdvertise,omitempty"`
	PrefixType            *CustomIPPrefixType `json:"prefixType,omitempty"`
	ProvisioningState     *ProvisioningState  `json:"provisioningState,omitempty"`
	PublicIPPrefixes      *[]SubResource      `json:"publicIpPrefixes,omitempty"`
	ResourceGuid          *string             `json:"resourceGuid,omitempty"`
	SignedMessage         *string             `json:"signedMessage,omitempty"`
}
