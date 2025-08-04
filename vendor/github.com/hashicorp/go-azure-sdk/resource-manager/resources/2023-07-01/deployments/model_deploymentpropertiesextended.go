package deployments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentPropertiesExtended struct {
	CorrelationId      *string                    `json:"correlationId,omitempty"`
	DebugSetting       *DebugSetting              `json:"debugSetting,omitempty"`
	Dependencies       *[]Dependency              `json:"dependencies,omitempty"`
	Duration           *string                    `json:"duration,omitempty"`
	Error              *ErrorResponse             `json:"error,omitempty"`
	Mode               *DeploymentMode            `json:"mode,omitempty"`
	OnErrorDeployment  *OnErrorDeploymentExtended `json:"onErrorDeployment,omitempty"`
	OutputResources    *[]ResourceReference       `json:"outputResources,omitempty"`
	Outputs            *interface{}               `json:"outputs,omitempty"`
	Parameters         *interface{}               `json:"parameters,omitempty"`
	ParametersLink     *ParametersLink            `json:"parametersLink,omitempty"`
	Providers          *[]Provider                `json:"providers,omitempty"`
	ProvisioningState  *ProvisioningState         `json:"provisioningState,omitempty"`
	TemplateHash       *string                    `json:"templateHash,omitempty"`
	TemplateLink       *TemplateLink              `json:"templateLink,omitempty"`
	Timestamp          *string                    `json:"timestamp,omitempty"`
	ValidatedResources *[]ResourceReference       `json:"validatedResources,omitempty"`
}

func (o *DeploymentPropertiesExtended) GetTimestampAsTime() (*time.Time, error) {
	if o.Timestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *DeploymentPropertiesExtended) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = &formatted
}
