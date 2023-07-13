package actiongroupsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogicAppReceiver struct {
	CallbackUrl          string `json:"callbackUrl"`
	Name                 string `json:"name"`
	ResourceId           string `json:"resourceId"`
	UseCommonAlertSchema *bool  `json:"useCommonAlertSchema,omitempty"`
}
