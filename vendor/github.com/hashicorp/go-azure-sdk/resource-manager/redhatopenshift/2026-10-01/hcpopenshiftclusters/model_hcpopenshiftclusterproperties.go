package hcpopenshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HcpOpenShiftClusterProperties struct {
	Api                     *ApiProfile                  `json:"api,omitempty"`
	Autoscaling             *ClusterAutoscalingProfile   `json:"autoscaling,omitempty"`
	ClusterImageRegistry    *ClusterImageRegistryProfile `json:"clusterImageRegistry,omitempty"`
	Console                 *ConsoleProfile              `json:"console,omitempty"`
	CryptoRestrictions      *CryptoRestrictions          `json:"cryptoRestrictions,omitempty"`
	Dns                     *DnsProfile                  `json:"dns,omitempty"`
	Etcd                    *EtcdProfile                 `json:"etcd,omitempty"`
	ImageDigestMirrors      *[]ImageDigestMirror         `json:"imageDigestMirrors,omitempty"`
	Ingress                 *IngressProfile              `json:"ingress,omitempty"`
	Network                 *NetworkProfile              `json:"network,omitempty"`
	NodeDrainTimeoutMinutes *int64                       `json:"nodeDrainTimeoutMinutes,omitempty"`
	Platform                PlatformProfile              `json:"platform"`
	ProvisioningState       *ProvisioningState           `json:"provisioningState,omitempty"`
	Status                  *ClusterResourceStatus       `json:"status,omitempty"`
	Version                 VersionProfile               `json:"version"`
}
