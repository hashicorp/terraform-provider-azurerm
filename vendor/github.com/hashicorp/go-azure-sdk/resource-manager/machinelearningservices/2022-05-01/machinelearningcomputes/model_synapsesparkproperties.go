package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SynapseSparkProperties struct {
	AutoPauseProperties *AutoPauseProperties `json:"autoPauseProperties,omitempty"`
	AutoScaleProperties *AutoScaleProperties `json:"autoScaleProperties,omitempty"`
	NodeCount           *int64               `json:"nodeCount,omitempty"`
	NodeSize            *string              `json:"nodeSize,omitempty"`
	NodeSizeFamily      *string              `json:"nodeSizeFamily,omitempty"`
	PoolName            *string              `json:"poolName,omitempty"`
	ResourceGroup       *string              `json:"resourceGroup,omitempty"`
	SparkVersion        *string              `json:"sparkVersion,omitempty"`
	SubscriptionId      *string              `json:"subscriptionId,omitempty"`
	WorkspaceName       *string              `json:"workspaceName,omitempty"`
}
