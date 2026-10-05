package channels

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ChannelUpdateParametersProperties struct {
	ExpirationTimeIfNotActivatedUtc *string                 `json:"expirationTimeIfNotActivatedUtc,omitempty"`
	PartnerTopicInfo                *PartnerUpdateTopicInfo `json:"partnerTopicInfo,omitempty"`
}

func (o *ChannelUpdateParametersProperties) GetExpirationTimeIfNotActivatedUtcAsTime() (*time.Time, error) {
	if o.ExpirationTimeIfNotActivatedUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationTimeIfNotActivatedUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ChannelUpdateParametersProperties) SetExpirationTimeIfNotActivatedUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationTimeIfNotActivatedUtc = &formatted
}
