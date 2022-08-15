package remediations

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationDeployment struct {
	CreatedOn            *string          `json:"createdOn,omitempty"`
	DeploymentId         *string          `json:"deploymentId,omitempty"`
	Error                *ErrorDefinition `json:"error,omitempty"`
	LastUpdatedOn        *string          `json:"lastUpdatedOn,omitempty"`
	RemediatedResourceId *string          `json:"remediatedResourceId,omitempty"`
	ResourceLocation     *string          `json:"resourceLocation,omitempty"`
	Status               *string          `json:"status,omitempty"`
}

func (o *RemediationDeployment) GetCreatedOnAsTime() (*time.Time, error) {
	if o.CreatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *RemediationDeployment) SetCreatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedOn = &formatted
}

func (o *RemediationDeployment) GetLastUpdatedOnAsTime() (*time.Time, error) {
	if o.LastUpdatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *RemediationDeployment) SetLastUpdatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedOn = &formatted
}
