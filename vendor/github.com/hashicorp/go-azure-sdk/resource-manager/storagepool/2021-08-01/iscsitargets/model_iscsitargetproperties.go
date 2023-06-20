package iscsitargets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IscsiTargetProperties struct {
	AclMode           IscsiTargetAclMode `json:"aclMode"`
	Endpoints         *[]string          `json:"endpoints,omitempty"`
	Luns              *[]IscsiLun        `json:"luns,omitempty"`
	Port              *int64             `json:"port,omitempty"`
	ProvisioningState ProvisioningStates `json:"provisioningState"`
	Sessions          *[]string          `json:"sessions,omitempty"`
	StaticAcls        *[]Acl             `json:"staticAcls,omitempty"`
	Status            OperationalStatus  `json:"status"`
	TargetIqn         string             `json:"targetIqn"`
}
