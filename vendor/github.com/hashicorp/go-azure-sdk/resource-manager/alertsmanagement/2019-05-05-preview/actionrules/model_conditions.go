package actionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Conditions struct {
	AlertContext       *Condition `json:"alertContext,omitempty"`
	AlertRuleId        *Condition `json:"alertRuleId,omitempty"`
	AlertRuleName      *Condition `json:"alertRuleName,omitempty"`
	Description        *Condition `json:"description,omitempty"`
	MonitorCondition   *Condition `json:"monitorCondition,omitempty"`
	MonitorService     *Condition `json:"monitorService,omitempty"`
	Severity           *Condition `json:"severity,omitempty"`
	TargetResourceType *Condition `json:"targetResourceType,omitempty"`
}
