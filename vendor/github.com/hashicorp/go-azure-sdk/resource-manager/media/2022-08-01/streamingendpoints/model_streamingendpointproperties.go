package streamingendpoints

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingEndpointProperties struct {
	AccessControl           *StreamingEndpointAccessControl `json:"accessControl,omitempty"`
	AvailabilitySetName     *string                         `json:"availabilitySetName,omitempty"`
	CdnEnabled              *bool                           `json:"cdnEnabled,omitempty"`
	CdnProfile              *string                         `json:"cdnProfile,omitempty"`
	CdnProvider             *string                         `json:"cdnProvider,omitempty"`
	Created                 *string                         `json:"created,omitempty"`
	CrossSiteAccessPolicies *CrossSiteAccessPolicies        `json:"crossSiteAccessPolicies,omitempty"`
	CustomHostNames         *[]string                       `json:"customHostNames,omitempty"`
	Description             *string                         `json:"description,omitempty"`
	FreeTrialEndTime        *string                         `json:"freeTrialEndTime,omitempty"`
	HostName                *string                         `json:"hostName,omitempty"`
	LastModified            *string                         `json:"lastModified,omitempty"`
	MaxCacheAge             *int64                          `json:"maxCacheAge,omitempty"`
	ProvisioningState       *string                         `json:"provisioningState,omitempty"`
	ResourceState           *StreamingEndpointResourceState `json:"resourceState,omitempty"`
	ScaleUnits              int64                           `json:"scaleUnits"`
}

func (o *StreamingEndpointProperties) GetCreatedAsTime() (*time.Time, error) {
	if o.Created == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Created, "2006-01-02T15:04:05Z07:00")
}

func (o *StreamingEndpointProperties) SetCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Created = &formatted
}

func (o *StreamingEndpointProperties) GetFreeTrialEndTimeAsTime() (*time.Time, error) {
	if o.FreeTrialEndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.FreeTrialEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *StreamingEndpointProperties) SetFreeTrialEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.FreeTrialEndTime = &formatted
}

func (o *StreamingEndpointProperties) GetLastModifiedAsTime() (*time.Time, error) {
	if o.LastModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModified, "2006-01-02T15:04:05Z07:00")
}

func (o *StreamingEndpointProperties) SetLastModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModified = &formatted
}
