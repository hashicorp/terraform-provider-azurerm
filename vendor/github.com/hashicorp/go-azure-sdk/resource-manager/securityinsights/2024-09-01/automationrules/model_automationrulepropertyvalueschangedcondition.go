package automationrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationRulePropertyValuesChangedCondition struct {
	ChangeType     *AutomationRulePropertyChangedConditionSupportedChangedType  `json:"changeType,omitempty"`
	Operator       *AutomationRulePropertyConditionSupportedOperator            `json:"operator,omitempty"`
	PropertyName   *AutomationRulePropertyChangedConditionSupportedPropertyType `json:"propertyName,omitempty"`
	PropertyValues *[]string                                                    `json:"propertyValues,omitempty"`
}
