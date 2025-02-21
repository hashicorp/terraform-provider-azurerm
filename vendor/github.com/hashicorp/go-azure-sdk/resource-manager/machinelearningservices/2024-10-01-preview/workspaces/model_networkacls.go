package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkAcls struct {
	DefaultAction *DefaultActionType `json:"defaultAction,omitempty"`
	IPRules       *[]IPRule          `json:"ipRules,omitempty"`
}
