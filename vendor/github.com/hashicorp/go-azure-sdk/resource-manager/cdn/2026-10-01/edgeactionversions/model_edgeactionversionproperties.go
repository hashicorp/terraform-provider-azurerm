package edgeactionversions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeActionVersionProperties struct {
	DeploymentType        EdgeActionVersionDeploymentType    `json:"deploymentType"`
	IsDefaultVersion      EdgeActionIsDefaultVersion         `json:"isDefaultVersion"`
	LastPackageUpdateTime *string                            `json:"lastPackageUpdateTime,omitempty"`
	ProvisioningState     *ProvisioningState                 `json:"provisioningState,omitempty"`
	ValidationStatus      *EdgeActionVersionValidationStatus `json:"validationStatus,omitempty"`
}

func (o *EdgeActionVersionProperties) GetLastPackageUpdateTimeAsTime() (*time.Time, error) {
	if o.LastPackageUpdateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastPackageUpdateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *EdgeActionVersionProperties) SetLastPackageUpdateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastPackageUpdateTime = &formatted
}
