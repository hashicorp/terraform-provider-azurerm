package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentCapacitySettings struct {
	DesignatedCapacity *int64 `json:"designatedCapacity,omitempty"`
	Priority           *int64 `json:"priority,omitempty"`
}
