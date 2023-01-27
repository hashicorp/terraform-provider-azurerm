package liveevents

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveEventProperties struct {
	Created                 *string                   `json:"created,omitempty"`
	CrossSiteAccessPolicies *CrossSiteAccessPolicies  `json:"crossSiteAccessPolicies,omitempty"`
	Description             *string                   `json:"description,omitempty"`
	Encoding                *LiveEventEncoding        `json:"encoding,omitempty"`
	HostnamePrefix          *string                   `json:"hostnamePrefix,omitempty"`
	Input                   LiveEventInput            `json:"input"`
	LastModified            *string                   `json:"lastModified,omitempty"`
	Preview                 *LiveEventPreview         `json:"preview,omitempty"`
	ProvisioningState       *string                   `json:"provisioningState,omitempty"`
	ResourceState           *LiveEventResourceState   `json:"resourceState,omitempty"`
	StreamOptions           *[]StreamOptionsFlag      `json:"streamOptions,omitempty"`
	Transcriptions          *[]LiveEventTranscription `json:"transcriptions,omitempty"`
	UseStaticHostname       *bool                     `json:"useStaticHostname,omitempty"`
}

func (o *LiveEventProperties) GetCreatedAsTime() (*time.Time, error) {
	if o.Created == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Created, "2006-01-02T15:04:05Z07:00")
}

func (o *LiveEventProperties) SetCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Created = &formatted
}

func (o *LiveEventProperties) GetLastModifiedAsTime() (*time.Time, error) {
	if o.LastModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModified, "2006-01-02T15:04:05Z07:00")
}

func (o *LiveEventProperties) SetLastModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModified = &formatted
}
