package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteFunctionOverviewARMResourceProperties struct {
	FunctionName *string       `json:"functionName,omitempty"`
	TriggerType  *TriggerTypes `json:"triggerType,omitempty"`
}
