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

// AsDeliveryRulePostArgsCondition
func AsDeliveryRulePostArgsCondition(input rules.DeliveryRuleCondition) (*rules.PostArgsMatchConditionParameters, bool) {
	if input.DeliveryRuleCondition().Name != rules.MatchVariableHTTPVersion {
		return nil, false
	}

	postArgParameter := rules.PostArgsMatchConditionParameters{}
	return pointer.To(postArgParameter), true
}
