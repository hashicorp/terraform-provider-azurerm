package networksecurityperimeteraccessrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NspAccessRuleProperties struct {
	AddressPrefixes           *[]string                   `json:"addressPrefixes,omitempty"`
	Direction                 *AccessRuleDirection        `json:"direction,omitempty"`
	EmailAddresses            *[]string                   `json:"emailAddresses,omitempty"`
	FullyQualifiedDomainNames *[]string                   `json:"fullyQualifiedDomainNames,omitempty"`
	NetworkSecurityPerimeters *[]PerimeterBasedAccessRule `json:"networkSecurityPerimeters,omitempty"`
	PhoneNumbers              *[]string                   `json:"phoneNumbers,omitempty"`
	ProvisioningState         *NspProvisioningState       `json:"provisioningState,omitempty"`
	ServiceTags               *[]string                   `json:"serviceTags,omitempty"`
	Subscriptions             *[]SubscriptionId           `json:"subscriptions,omitempty"`
}
