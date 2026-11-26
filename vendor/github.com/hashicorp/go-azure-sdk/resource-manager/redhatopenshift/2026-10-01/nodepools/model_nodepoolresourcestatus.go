package nodepools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodePoolResourceStatus struct {
	ActiveVersions *[]NodePoolActiveVersion `json:"activeVersions,omitempty"`
	Conditions     *[]Condition             `json:"conditions,omitempty"`
}
