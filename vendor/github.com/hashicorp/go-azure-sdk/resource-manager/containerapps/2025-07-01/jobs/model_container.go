package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Container struct {
	Args         *[]string            `json:"args,omitempty"`
	Command      *[]string            `json:"command,omitempty"`
	Env          *[]EnvironmentVar    `json:"env,omitempty"`
	Image        *string              `json:"image,omitempty"`
	Name         *string              `json:"name,omitempty"`
	Probes       *[]ContainerAppProbe `json:"probes,omitempty"`
	Resources    *ContainerResources  `json:"resources,omitempty"`
	VolumeMounts *[]VolumeMount       `json:"volumeMounts,omitempty"`
}
