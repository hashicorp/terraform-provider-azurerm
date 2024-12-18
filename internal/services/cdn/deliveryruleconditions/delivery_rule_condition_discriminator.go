// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package deliveryruleconditions

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rules"
)

func AsDeliveryRuleCookiesCondition(input rules.DeliveryRuleCondition) (*rules.CookiesMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableCookies {
		return nil, false
	}

	cookieParameter := rules.CookiesMatchConditionParameters{}
	return pointer.To(cookieParameter), true
}

func AsDeliveryRuleIsDeviceCondition(input rules.DeliveryRuleCondition) (*rules.IsDeviceMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableIsDevice {
		return nil, false
	}

	deviceParameter := rules.IsDeviceMatchConditionParameters{}
	return pointer.To(deviceParameter), true
}

func AsDeliveryRuleHTTPVersionCondition(input rules.DeliveryRuleCondition) (*rules.HTTPVersionMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableHTTPVersion {
		return nil, false
	}

	deviceParameter := rules.HTTPVersionMatchConditionParameters{}
	return pointer.To(deviceParameter), true
}

func AsDeliveryRulePostArgsCondition(input rules.DeliveryRuleCondition) (*rules.PostArgsMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariablePostArgs {
		return nil, false
	}

	postArgParameter := rules.PostArgsMatchConditionParameters{}
	return pointer.To(postArgParameter), true
}

func AsDeliveryRuleQueryStringCondition(input rules.DeliveryRuleCondition) (*rules.QueryStringMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableQueryString {
		return nil, false
	}

	postArgParameter := rules.QueryStringMatchConditionParameters{}
	return pointer.To(postArgParameter), true
}

func AsDeliveryRuleRemoteAddressCondition(input rules.DeliveryRuleCondition) (*rules.RemoteAddressMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableRemoteAddress {
		return nil, false
	}

	remoteAddressParameter := rules.RemoteAddressMatchConditionParameters{}
	return pointer.To(remoteAddressParameter), true
}

func AsDeliveryRuleRequestBodyCondition(input rules.DeliveryRuleCondition) (*rules.RequestBodyMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableRequestBody {
		return nil, false
	}

	requestBodyParameter := rules.RequestBodyMatchConditionParameters{}
	return pointer.To(requestBodyParameter), true
}

func AsDeliveryRuleRequestHeaderCondition(input rules.DeliveryRuleCondition) (*rules.RequestHeaderMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableRequestHeader {
		return nil, false
	}

	requestHeaderParameter := rules.RequestHeaderMatchConditionParameters{}
	return pointer.To(requestHeaderParameter), true
}

func AsDeliveryRuleRequestMethodCondition(input rules.DeliveryRuleCondition) (*rules.RequestMethodMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableRequestMethod {
		return nil, false
	}

	requestMethodParameter := rules.RequestMethodMatchConditionParameters{}
	return pointer.To(requestMethodParameter), true
}

func AsDeliveryRuleRequestSchemeCondition(input rules.DeliveryRuleCondition) (*rules.RequestSchemeMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableRequestScheme {
		return nil, false
	}

	requestSchemeParameter := rules.RequestSchemeMatchConditionParameters{}
	return pointer.To(requestSchemeParameter), true
}

func AsDeliveryRuleRequestUriCondition(input rules.DeliveryRuleCondition) (*rules.RequestUriMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableRequestUri {
		return nil, false
	}

	requestUriParameter := rules.RequestUriMatchConditionParameters{}
	return pointer.To(requestUriParameter), true
}

func AsDeliveryRuleURLFileExtensionCondition(input rules.DeliveryRuleCondition) (*rules.URLFileExtensionMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableURLFileExtension {
		return nil, false
	}

	urlFileExtensionParameter := rules.URLFileExtensionMatchConditionParameters{}
	return pointer.To(urlFileExtensionParameter), true
}

func AsDeliveryRuleURLFileNameCondition(input rules.DeliveryRuleCondition) (*rules.URLFileNameMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableURLFileName {
		return nil, false
	}

	urlFileNameParameter := rules.URLFileNameMatchConditionParameters{}
	return pointer.To(urlFileNameParameter), true
}

func AsDeliveryRuleURLPathCondition(input rules.DeliveryRuleCondition) (*rules.URLPathMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableURLPath {
		return nil, false
	}

	urlPathParameter := rules.URLPathMatchConditionParameters{}
	return pointer.To(urlPathParameter), true
}
