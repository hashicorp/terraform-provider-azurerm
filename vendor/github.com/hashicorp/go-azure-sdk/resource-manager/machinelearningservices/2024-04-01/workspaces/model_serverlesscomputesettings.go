package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerlessComputeSettings struct {
	ServerlessComputeCustomSubnet *string `json:"serverlessComputeCustomSubnet,omitempty"`
	ServerlessComputeNoPublicIP   *bool   `json:"serverlessComputeNoPublicIP,omitempty"`
}
