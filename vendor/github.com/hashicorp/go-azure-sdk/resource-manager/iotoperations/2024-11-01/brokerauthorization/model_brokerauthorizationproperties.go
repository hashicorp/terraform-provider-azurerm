package brokerauthorization

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BrokerAuthorizationProperties struct {
	AuthorizationPolicies AuthorizationConfig `json:"authorizationPolicies"`
	ProvisioningState     *ProvisioningState  `json:"provisioningState,omitempty"`
}
