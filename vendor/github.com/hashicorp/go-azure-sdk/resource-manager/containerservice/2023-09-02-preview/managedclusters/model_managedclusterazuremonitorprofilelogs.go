package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterAzureMonitorProfileLogs struct {
	AppMonitoring     *ManagedClusterAzureMonitorProfileAppMonitoring     `json:"appMonitoring,omitempty"`
	ContainerInsights *ManagedClusterAzureMonitorProfileContainerInsights `json:"containerInsights,omitempty"`
}
