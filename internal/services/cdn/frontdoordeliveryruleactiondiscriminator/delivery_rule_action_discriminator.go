// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frontdoordeliveryruleactiondiscriminator

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rules"
)

func AsDeliveryRuleCacheExpirationAction(input rules.DeliveryRuleAction) (*rules.DeliveryRuleCacheExpirationAction, bool) {
	if input.DeliveryRuleAction().Name != rules.DeliveryRuleActionNameRouteConfigurationOverride {
		return nil, false
	}

	cacheExpiration := rules.DeliveryRuleCacheExpirationAction{}
	return pointer.To(cacheExpiration), true
}

func AsDeliveryRuleRouteConfigurationOverrideAction(input rules.DeliveryRuleAction) (*rules.DeliveryRuleRouteConfigurationOverrideAction, bool) {
	if input.DeliveryRuleAction().Name != rules.DeliveryRuleActionNameRouteConfigurationOverride {
		return nil, false
	}

	routeConfigurationOverride := rules.DeliveryRuleRouteConfigurationOverrideAction{}
	return pointer.To(routeConfigurationOverride), true
}

func AsDeliveryRuleResponseHeaderAction(input rules.DeliveryRuleAction) (*rules.DeliveryRuleResponseHeaderAction, bool) {
	if input.DeliveryRuleAction().Name != rules.DeliveryRuleActionNameModifyResponseHeader {
		return nil, false
	}

	responseHeader := rules.DeliveryRuleResponseHeaderAction{}
	return pointer.To(responseHeader), true
}

func AsDeliveryRuleRequestHeaderAction(input rules.DeliveryRuleAction) (*rules.DeliveryRuleRequestHeaderAction, bool) {
	if input.DeliveryRuleAction().Name != rules.DeliveryRuleActionNameModifyRequestHeader {
		return nil, false
	}

	requestHeader := rules.DeliveryRuleRequestHeaderAction{}
	return pointer.To(requestHeader), true
}

func AsDeliveryRuleUrlRewriteAction(input rules.DeliveryRuleAction) (*rules.URLRewriteAction, bool) {
	if input.DeliveryRuleAction().Name != rules.DeliveryRuleActionNameURLRewrite {
		return nil, false
	}

	urlRewrite := rules.URLRewriteAction{}
	return pointer.To(urlRewrite), true
}

func AsDeliveryRuleUrlRedirectAction(input rules.DeliveryRuleAction) (*rules.URLRedirectAction, bool) {
	if input.DeliveryRuleAction().Name != rules.DeliveryRuleActionNameURLRedirect {
		return nil, false
	}

	urlRedirect := rules.URLRedirectAction{}
	return pointer.To(urlRedirect), true
}
