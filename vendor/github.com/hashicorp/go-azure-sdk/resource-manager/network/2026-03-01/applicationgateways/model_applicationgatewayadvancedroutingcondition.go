package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayAdvancedRoutingCondition struct {
	ConditionType        ApplicationGatewayAdvancedRoutingConditionType         `json:"conditionType"`
	PropertyName         *string                                                `json:"propertyName,omitempty"`
	PropertyValueMatcher *ApplicationGatewayAdvancedRoutingPropertyValueMatcher `json:"propertyValueMatcher,omitempty"`
	PropertyValues       *[]string                                              `json:"propertyValues,omitempty"`
}
