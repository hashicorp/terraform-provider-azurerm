package activitylogalertsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRuleAllOfCondition struct {
	AllOf []AlertRuleAnyOfOrLeafCondition `json:"allOf"`
}
