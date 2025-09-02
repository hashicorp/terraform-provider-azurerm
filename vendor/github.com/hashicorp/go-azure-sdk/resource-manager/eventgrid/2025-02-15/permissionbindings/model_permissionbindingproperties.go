package permissionbindings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PermissionBindingProperties struct {
	ClientGroupName   *string                             `json:"clientGroupName,omitempty"`
	Description       *string                             `json:"description,omitempty"`
	Permission        *PermissionType                     `json:"permission,omitempty"`
	ProvisioningState *PermissionBindingProvisioningState `json:"provisioningState,omitempty"`
	TopicSpaceName    *string                             `json:"topicSpaceName,omitempty"`
}
