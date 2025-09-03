package commitmentplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommitmentPlanProperties struct {
	AutoRenew          *bool                            `json:"autoRenew,omitempty"`
	CommitmentPlanGuid *string                          `json:"commitmentPlanGuid,omitempty"`
	Current            *CommitmentPeriod                `json:"current,omitempty"`
	HostingModel       *HostingModel                    `json:"hostingModel,omitempty"`
	Last               *CommitmentPeriod                `json:"last,omitempty"`
	Next               *CommitmentPeriod                `json:"next,omitempty"`
	PlanType           *string                          `json:"planType,omitempty"`
	ProvisioningIssues *[]string                        `json:"provisioningIssues,omitempty"`
	ProvisioningState  *CommitmentPlanProvisioningState `json:"provisioningState,omitempty"`
}
