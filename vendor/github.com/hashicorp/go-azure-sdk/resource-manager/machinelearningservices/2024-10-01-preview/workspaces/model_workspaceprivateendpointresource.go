package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacePrivateEndpointResource struct {
	Id          *string `json:"id,omitempty"`
	SubnetArmId *string `json:"subnetArmId,omitempty"`
}
