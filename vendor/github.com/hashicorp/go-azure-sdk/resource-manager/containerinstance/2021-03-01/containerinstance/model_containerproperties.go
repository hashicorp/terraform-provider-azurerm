package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerProperties struct {
	Command              *[]string                        `json:"command,omitempty"`
	EnvironmentVariables *[]EnvironmentVariable           `json:"environmentVariables,omitempty"`
	Image                string                           `json:"image"`
	InstanceView         *ContainerPropertiesInstanceView `json:"instanceView,omitempty"`
	LivenessProbe        *ContainerProbe                  `json:"livenessProbe,omitempty"`
	Ports                *[]ContainerPort                 `json:"ports,omitempty"`
	ReadinessProbe       *ContainerProbe                  `json:"readinessProbe,omitempty"`
	Resources            ResourceRequirements             `json:"resources"`
	VolumeMounts         *[]VolumeMount                   `json:"volumeMounts,omitempty"`
}
