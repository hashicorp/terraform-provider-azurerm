package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayAdvancedRoutingPropertyValueMatcher struct {
	IgnoreCase *bool  `json:"ignoreCase,omitempty"`
	Negate     *bool  `json:"negate,omitempty"`
	Pattern    string `json:"pattern"`
}
