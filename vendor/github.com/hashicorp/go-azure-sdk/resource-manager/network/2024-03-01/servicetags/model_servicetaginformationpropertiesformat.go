package servicetags

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceTagInformationPropertiesFormat struct {
	AddressPrefixes *[]string `json:"addressPrefixes,omitempty"`
	ChangeNumber    *string   `json:"changeNumber,omitempty"`
	Region          *string   `json:"region,omitempty"`
	State           *string   `json:"state,omitempty"`
	SystemService   *string   `json:"systemService,omitempty"`
}
