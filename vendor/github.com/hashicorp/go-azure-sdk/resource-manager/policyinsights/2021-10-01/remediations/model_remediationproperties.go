package remediations

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationProperties struct {
	CorrelationId               *string                                `json:"correlationId,omitempty"`
	CreatedOn                   *string                                `json:"createdOn,omitempty"`
	DeploymentStatus            *RemediationDeploymentSummary          `json:"deploymentStatus,omitempty"`
	FailureThreshold            *RemediationPropertiesFailureThreshold `json:"failureThreshold,omitempty"`
	Filters                     *RemediationFilters                    `json:"filters,omitempty"`
	LastUpdatedOn               *string                                `json:"lastUpdatedOn,omitempty"`
	ParallelDeployments         *int64                                 `json:"parallelDeployments,omitempty"`
	PolicyAssignmentId          *string                                `json:"policyAssignmentId,omitempty"`
	PolicyDefinitionReferenceId *string                                `json:"policyDefinitionReferenceId,omitempty"`
	ProvisioningState           *string                                `json:"provisioningState,omitempty"`
	ResourceCount               *int64                                 `json:"resourceCount,omitempty"`
	ResourceDiscoveryMode       *ResourceDiscoveryMode                 `json:"resourceDiscoveryMode,omitempty"`
	StatusMessage               *string                                `json:"statusMessage,omitempty"`
}

func (o *RemediationProperties) GetCreatedOnAsTime() (*time.Time, error) {
	if o.CreatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *RemediationProperties) SetCreatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedOn = &formatted
}

func (o *RemediationProperties) GetLastUpdatedOnAsTime() (*time.Time, error) {
	if o.LastUpdatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *RemediationProperties) SetLastUpdatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedOn = &formatted
}
