package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Role struct {
	Autoscale             *Autoscale             `json:"autoscale,omitempty"`
	DataDisksGroups       *[]DataDisksGroups     `json:"dataDisksGroups,omitempty"`
	EncryptDataDisks      *bool                  `json:"encryptDataDisks,omitempty"`
	HardwareProfile       *HardwareProfile       `json:"hardwareProfile,omitempty"`
	MinInstanceCount      *int64                 `json:"minInstanceCount,omitempty"`
	Name                  *string                `json:"name,omitempty"`
	OsProfile             *OsProfile             `json:"osProfile,omitempty"`
	ScriptActions         *[]ScriptAction        `json:"scriptActions,omitempty"`
	TargetInstanceCount   *int64                 `json:"targetInstanceCount,omitempty"`
	VMGroupName           *string                `json:"VMGroupName,omitempty"`
	VirtualNetworkProfile *VirtualNetworkProfile `json:"virtualNetworkProfile,omitempty"`
}
