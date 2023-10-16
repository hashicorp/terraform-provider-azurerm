package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentSettings struct {
	AddonConfigs                  *map[string]interface{} `json:"addonConfigs,omitempty"`
	Apms                          *[]ApmReference         `json:"apms,omitempty"`
	ContainerProbeSettings        *ContainerProbeSettings `json:"containerProbeSettings,omitempty"`
	EnvironmentVariables          *map[string]string      `json:"environmentVariables,omitempty"`
	LivenessProbe                 *Probe                  `json:"livenessProbe,omitempty"`
	ReadinessProbe                *Probe                  `json:"readinessProbe,omitempty"`
	ResourceRequests              *ResourceRequests       `json:"resourceRequests,omitempty"`
	Scale                         *Scale                  `json:"scale,omitempty"`
	StartupProbe                  *Probe                  `json:"startupProbe,omitempty"`
	TerminationGracePeriodSeconds *int64                  `json:"terminationGracePeriodSeconds,omitempty"`
}
