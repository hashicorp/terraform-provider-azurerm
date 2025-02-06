package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkAcls struct {
	IPRules           *[]IPRule                        `json:"ipRules,omitempty"`
	PublicNetworkRule *PublicInternetAuthorizationRule `json:"publicNetworkRule,omitempty"`
}
