package networkmanagers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkManagerDeploymentStatus struct {
	CommitTime       *string            `json:"commitTime,omitempty"`
	ConfigurationIds *[]string          `json:"configurationIds,omitempty"`
	DeploymentStatus *DeploymentStatus  `json:"deploymentStatus,omitempty"`
	DeploymentType   *ConfigurationType `json:"deploymentType,omitempty"`
	ErrorMessage     *string            `json:"errorMessage,omitempty"`
	Region           *string            `json:"region,omitempty"`
}

func (o *NetworkManagerDeploymentStatus) GetCommitTimeAsTime() (*time.Time, error) {
	if o.CommitTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CommitTime, "2006-01-02T15:04:05Z07:00")
}

func (o *NetworkManagerDeploymentStatus) SetCommitTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CommitTime = &formatted
}
