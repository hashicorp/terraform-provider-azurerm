package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeDefinition struct {
	Bind        *BindOptions          `json:"bind,omitempty"`
	Consistency *string               `json:"consistency,omitempty"`
	ReadOnly    *bool                 `json:"readOnly,omitempty"`
	Source      *string               `json:"source,omitempty"`
	Target      *string               `json:"target,omitempty"`
	Tmpfs       *TmpfsOptions         `json:"tmpfs,omitempty"`
	Type        *VolumeDefinitionType `json:"type,omitempty"`
	Volume      *VolumeOptions        `json:"volume,omitempty"`
}
