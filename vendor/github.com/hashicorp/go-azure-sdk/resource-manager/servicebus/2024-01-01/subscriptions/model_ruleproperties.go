package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Ruleproperties struct {
	Action            *Action            `json:"action,omitempty"`
	CorrelationFilter *CorrelationFilter `json:"correlationFilter,omitempty"`
	FilterType        *FilterType        `json:"filterType,omitempty"`
	SqlFilter         *SqlFilter         `json:"sqlFilter,omitempty"`
}
