package communicationsgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrimaryRegionProperties struct {
	AllowedMediaSourceAddressPrefixes     *[]string `json:"allowedMediaSourceAddressPrefixes,omitempty"`
	AllowedSignalingSourceAddressPrefixes *[]string `json:"allowedSignalingSourceAddressPrefixes,omitempty"`
	EsrpAddresses                         *[]string `json:"esrpAddresses,omitempty"`
	OperatorAddresses                     []string  `json:"operatorAddresses"`
}
