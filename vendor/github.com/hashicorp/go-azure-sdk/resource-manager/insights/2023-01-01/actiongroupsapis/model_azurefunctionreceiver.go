package actiongroupsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFunctionReceiver struct {
	FunctionAppResourceId string `json:"functionAppResourceId"`
	FunctionName          string `json:"functionName"`
	HTTPTriggerURL        string `json:"httpTriggerUrl"`
	Name                  string `json:"name"`
	UseCommonAlertSchema  *bool  `json:"useCommonAlertSchema,omitempty"`
}
