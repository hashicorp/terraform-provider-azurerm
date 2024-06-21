package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayRewriteRuleActionSet struct {
	RequestHeaderConfigurations  *[]ApplicationGatewayHeaderConfiguration `json:"requestHeaderConfigurations,omitempty"`
	ResponseHeaderConfigurations *[]ApplicationGatewayHeaderConfiguration `json:"responseHeaderConfigurations,omitempty"`
	UrlConfiguration             *ApplicationGatewayUrlConfiguration      `json:"urlConfiguration,omitempty"`
}
