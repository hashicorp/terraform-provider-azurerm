package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeClaimSpec struct {
	AccessModes      *[]string                        `json:"accessModes,omitempty"`
	DataSource       *LocalKubernetesReference        `json:"dataSource,omitempty"`
	DataSourceRef    *KubernetesReference             `json:"dataSourceRef,omitempty"`
	Resources        *VolumeClaimResourceRequirements `json:"resources,omitempty"`
	Selector         *VolumeClaimSpecSelector         `json:"selector,omitempty"`
	StorageClassName *string                          `json:"storageClassName,omitempty"`
	VolumeMode       *string                          `json:"volumeMode,omitempty"`
	VolumeName       *string                          `json:"volumeName,omitempty"`
}
