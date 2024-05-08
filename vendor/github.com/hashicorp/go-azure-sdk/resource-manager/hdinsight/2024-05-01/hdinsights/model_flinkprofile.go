package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FlinkProfile struct {
	CatalogOptions *FlinkCatalogOptions       `json:"catalogOptions,omitempty"`
	DeploymentMode *DeploymentMode            `json:"deploymentMode,omitempty"`
	HistoryServer  *ComputeResourceDefinition `json:"historyServer,omitempty"`
	JobManager     ComputeResourceDefinition  `json:"jobManager"`
	JobSpec        *FlinkJobProfile           `json:"jobSpec,omitempty"`
	NumReplicas    *int64                     `json:"numReplicas,omitempty"`
	Storage        FlinkStorageProfile        `json:"storage"`
	TaskManager    ComputeResourceDefinition  `json:"taskManager"`
}
