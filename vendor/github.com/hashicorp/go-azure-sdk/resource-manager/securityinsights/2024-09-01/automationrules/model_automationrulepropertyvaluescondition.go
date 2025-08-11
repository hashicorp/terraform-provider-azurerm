package automationrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationRulePropertyValuesCondition struct {
	Operator       *AutomationRulePropertyConditionSupportedOperator `json:"operator,omitempty"`
	PropertyName   *AutomationRulePropertyConditionSupportedProperty `json:"propertyName,omitempty"`
	PropertyValues *[]string                                         `json:"propertyValues,omitempty"`
}
