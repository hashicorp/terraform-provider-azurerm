package containerinstance

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerGroupProfileProperties struct {
	ConfidentialComputeProperties *ConfidentialComputeProperties `json:"confidentialComputeProperties,omitempty"`
	Containers                    []Container                    `json:"containers"`
	Diagnostics                   *ContainerGroupDiagnostics     `json:"diagnostics,omitempty"`
	EncryptionProperties          *EncryptionProperties          `json:"encryptionProperties,omitempty"`
	Extensions                    *[]DeploymentExtensionSpec     `json:"extensions,omitempty"`
	IPAddress                     *IPAddress                     `json:"ipAddress,omitempty"`
	ImageRegistryCredentials      *[]ImageRegistryCredential     `json:"imageRegistryCredentials,omitempty"`
	InitContainers                *[]InitContainerDefinition     `json:"initContainers,omitempty"`
	OsType                        OperatingSystemTypes           `json:"osType"`
	Priority                      *ContainerGroupPriority        `json:"priority,omitempty"`
	RegisteredRevisions           *[]int64                       `json:"registeredRevisions,omitempty"`
	RestartPolicy                 *ContainerGroupRestartPolicy   `json:"restartPolicy,omitempty"`
	Revision                      *int64                         `json:"revision,omitempty"`
	SecurityContext               *SecurityContextDefinition     `json:"securityContext,omitempty"`
	ShutdownGracePeriod           *string                        `json:"shutdownGracePeriod,omitempty"`
	Sku                           *ContainerGroupSku             `json:"sku,omitempty"`
	TimeToLive                    *string                        `json:"timeToLive,omitempty"`
	UseKrypton                    *bool                          `json:"useKrypton,omitempty"`
	Volumes                       *[]Volume                      `json:"volumes,omitempty"`
}

func (o *ContainerGroupProfileProperties) GetShutdownGracePeriodAsTime() (*time.Time, error) {
	if o.ShutdownGracePeriod == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ShutdownGracePeriod, "2006-01-02T15:04:05Z07:00")
}

func (o *ContainerGroupProfileProperties) SetShutdownGracePeriodAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ShutdownGracePeriod = &formatted
}

func (o *ContainerGroupProfileProperties) GetTimeToLiveAsTime() (*time.Time, error) {
	if o.TimeToLive == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeToLive, "2006-01-02T15:04:05Z07:00")
}

func (o *ContainerGroupProfileProperties) SetTimeToLiveAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeToLive = &formatted
}
