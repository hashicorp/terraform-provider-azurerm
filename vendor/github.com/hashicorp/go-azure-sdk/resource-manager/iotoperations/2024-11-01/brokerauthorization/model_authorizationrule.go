package brokerauthorization

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationRule struct {
	BrokerResources     []BrokerResourceRule      `json:"brokerResources"`
	Principals          PrincipalDefinition       `json:"principals"`
	StateStoreResources *[]StateStoreResourceRule `json:"stateStoreResources,omitempty"`
}
