package hcpopenshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HcpOpenShiftClusterPropertiesUpdate struct {
	Api                     *ApiProfileUpdate          `json:"api,omitempty"`
	Autoscaling             *ClusterAutoscalingProfile `json:"autoscaling,omitempty"`
	Etcd                    *EtcdProfileUpdate         `json:"etcd,omitempty"`
	ImageDigestMirrors      *[]ImageDigestMirror       `json:"imageDigestMirrors,omitempty"`
	NodeDrainTimeoutMinutes *int64                     `json:"nodeDrainTimeoutMinutes,omitempty"`
	Platform                *PlatformProfileUpdate     `json:"platform,omitempty"`
	Version                 *VersionProfileUpdate      `json:"version,omitempty"`
}
