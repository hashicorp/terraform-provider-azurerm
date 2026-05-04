package webapps

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SiteContainerProperties struct {
	AuthType                    *AuthType              `json:"authType,omitempty"`
	CreatedTime                 *string                `json:"createdTime,omitempty"`
	EnvironmentVariables        *[]EnvironmentVariable `json:"environmentVariables,omitempty"`
	Image                       string                 `json:"image"`
	IsMain                      bool                   `json:"isMain"`
	LastModifiedTime            *string                `json:"lastModifiedTime,omitempty"`
	PasswordSecret              *string                `json:"passwordSecret,omitempty"`
	StartUpCommand              *string                `json:"startUpCommand,omitempty"`
	TargetPort                  *string                `json:"targetPort,omitempty"`
	UserManagedIdentityClientId *string                `json:"userManagedIdentityClientId,omitempty"`
	UserName                    *string                `json:"userName,omitempty"`
	VolumeMounts                *[]VolumeMount         `json:"volumeMounts,omitempty"`
}

func (o *SiteContainerProperties) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SiteContainerProperties) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}

func (o *SiteContainerProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SiteContainerProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
