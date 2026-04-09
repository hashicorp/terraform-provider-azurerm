package labs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArtifactInstallProperties struct {
	ArtifactId               *string                        `json:"artifactId,omitempty"`
	ArtifactTitle            *string                        `json:"artifactTitle,omitempty"`
	DeploymentStatusMessage  *string                        `json:"deploymentStatusMessage,omitempty"`
	InstallTime              *string                        `json:"installTime,omitempty"`
	Parameters               *[]ArtifactParameterProperties `json:"parameters,omitempty"`
	Status                   *string                        `json:"status,omitempty"`
	VMExtensionStatusMessage *string                        `json:"vmExtensionStatusMessage,omitempty"`
}

func (o *ArtifactInstallProperties) GetInstallTimeAsTime() (*time.Time, error) {
	if o.InstallTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.InstallTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ArtifactInstallProperties) SetInstallTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.InstallTime = &formatted
}
