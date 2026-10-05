package channels

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ChannelProperties struct {
	ChannelType                     *ChannelType              `json:"channelType,omitempty"`
	ExpirationTimeIfNotActivatedUtc *string                   `json:"expirationTimeIfNotActivatedUtc,omitempty"`
	MessageForActivation            *string                   `json:"messageForActivation,omitempty"`
	PartnerTopicInfo                *PartnerTopicInfo         `json:"partnerTopicInfo,omitempty"`
	ProvisioningState               *ChannelProvisioningState `json:"provisioningState,omitempty"`
	ReadinessState                  *ReadinessState           `json:"readinessState,omitempty"`
}

func (o *ChannelProperties) GetExpirationTimeIfNotActivatedUtcAsTime() (*time.Time, error) {
	if o.ExpirationTimeIfNotActivatedUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationTimeIfNotActivatedUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ChannelProperties) SetExpirationTimeIfNotActivatedUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationTimeIfNotActivatedUtc = &formatted
}
