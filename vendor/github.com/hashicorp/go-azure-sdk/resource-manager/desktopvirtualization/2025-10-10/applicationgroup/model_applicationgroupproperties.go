package applicationgroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGroupProperties struct {
	ApplicationGroupType ApplicationGroupType `json:"applicationGroupType"`
	CloudPcResource      *bool                `json:"cloudPcResource,omitempty"`
	Description          *string              `json:"description,omitempty"`
	FriendlyName         *string              `json:"friendlyName,omitempty"`
	HostPoolArmPath      string               `json:"hostPoolArmPath"`
	ObjectId             *string              `json:"objectId,omitempty"`
	ShowInFeed           *bool                `json:"showInFeed,omitempty"`
	WorkspaceArmPath     *string              `json:"workspaceArmPath,omitempty"`
}
