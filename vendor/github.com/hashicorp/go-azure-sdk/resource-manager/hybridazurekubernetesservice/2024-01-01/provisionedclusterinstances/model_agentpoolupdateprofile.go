package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentPoolUpdateProfile struct {
	Count             *int64  `json:"count,omitempty"`
	KubernetesVersion *string `json:"kubernetesVersion,omitempty"`
	VMSize            *string `json:"vmSize,omitempty"`
}
