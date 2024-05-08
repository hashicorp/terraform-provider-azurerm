package serviceendpointpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FlowLogPropertiesFormat struct {
	Enabled                    *bool                       `json:"enabled,omitempty"`
	FlowAnalyticsConfiguration *TrafficAnalyticsProperties `json:"flowAnalyticsConfiguration,omitempty"`
	Format                     *FlowLogFormatParameters    `json:"format,omitempty"`
	ProvisioningState          *ProvisioningState          `json:"provisioningState,omitempty"`
	RetentionPolicy            *RetentionPolicyParameters  `json:"retentionPolicy,omitempty"`
	StorageId                  string                      `json:"storageId"`
	TargetResourceGuid         *string                     `json:"targetResourceGuid,omitempty"`
	TargetResourceId           string                      `json:"targetResourceId"`
}
