package actiongroupsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArmRoleReceiver struct {
	Name                 string `json:"name"`
	RoleId               string `json:"roleId"`
	UseCommonAlertSchema *bool  `json:"useCommonAlertSchema,omitempty"`
}
