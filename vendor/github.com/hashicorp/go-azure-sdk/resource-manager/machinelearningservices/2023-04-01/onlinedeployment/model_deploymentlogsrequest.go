package onlinedeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentLogsRequest struct {
	ContainerType *ContainerType `json:"containerType,omitempty"`
	Tail          *int64         `json:"tail,omitempty"`
}
