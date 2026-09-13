package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InitContainerPropertiesDefinition struct {
	Command              *[]string                                      `json:"command,omitempty"`
	EnvironmentVariables *[]EnvironmentVariable                         `json:"environmentVariables,omitempty"`
	Image                *string                                        `json:"image,omitempty"`
	InstanceView         *InitContainerPropertiesDefinitionInstanceView `json:"instanceView,omitempty"`
	SecurityContext      *SecurityContextDefinition                     `json:"securityContext,omitempty"`
	VolumeMounts         *[]VolumeMount                                 `json:"volumeMounts,omitempty"`
}
