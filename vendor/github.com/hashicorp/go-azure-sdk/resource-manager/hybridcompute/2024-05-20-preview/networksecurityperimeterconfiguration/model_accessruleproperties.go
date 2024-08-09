package networksecurityperimeterconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessRuleProperties struct {
	AddressPrefixes *[]string            `json:"addressPrefixes,omitempty"`
	Direction       *AccessRuleDirection `json:"direction,omitempty"`
}
