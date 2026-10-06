package deploymentstacksatresourcegroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionOnUnmanage struct {
	ManagementGroups *UnmanageActionManagementGroupMode `json:"managementGroups,omitempty"`
	ResourceGroups   *UnmanageActionResourceGroupMode   `json:"resourceGroups,omitempty"`
	Resources        UnmanageActionResourceMode         `json:"resources"`
}
