package integrationruntimedisableinteractivequery

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeDataFlowProperties struct {
	Cleanup          *bool                                                          `json:"cleanup,omitempty"`
	ComputeType      *DataFlowComputeType                                           `json:"computeType,omitempty"`
	CoreCount        *int64                                                         `json:"coreCount,omitempty"`
	CustomProperties *[]IntegrationRuntimeDataFlowPropertiesCustomPropertiesInlined `json:"customProperties,omitempty"`
	TimeToLive       *int64                                                         `json:"timeToLive,omitempty"`
}
