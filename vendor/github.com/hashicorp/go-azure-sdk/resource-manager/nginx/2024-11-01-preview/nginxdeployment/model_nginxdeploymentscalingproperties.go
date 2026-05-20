package nginxdeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentScalingProperties struct {
	AutoScaleSettings *NginxDeploymentScalingPropertiesAutoScaleSettings `json:"autoScaleSettings,omitempty"`
	Capacity          *int64                                             `json:"capacity,omitempty"`
}
