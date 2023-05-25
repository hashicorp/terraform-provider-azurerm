package prometheusrulegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrometheusRuleGroupAction struct {
	ActionGroupId    *string            `json:"actionGroupId,omitempty"`
	ActionProperties *map[string]string `json:"actionProperties,omitempty"`
}
