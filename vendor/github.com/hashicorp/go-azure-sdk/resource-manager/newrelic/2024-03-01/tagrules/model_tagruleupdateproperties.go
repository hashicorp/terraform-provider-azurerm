package tagrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagRuleUpdateProperties struct {
	LogRules    *LogRules    `json:"logRules,omitempty"`
	MetricRules *MetricRules `json:"metricRules,omitempty"`
}
