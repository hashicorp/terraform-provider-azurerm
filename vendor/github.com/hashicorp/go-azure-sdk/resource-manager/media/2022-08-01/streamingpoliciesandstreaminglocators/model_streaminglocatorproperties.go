package streamingpoliciesandstreaminglocators

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingLocatorProperties struct {
	AlternativeMediaId          *string                       `json:"alternativeMediaId,omitempty"`
	AssetName                   string                        `json:"assetName"`
	ContentKeys                 *[]StreamingLocatorContentKey `json:"contentKeys,omitempty"`
	Created                     *string                       `json:"created,omitempty"`
	DefaultContentKeyPolicyName *string                       `json:"defaultContentKeyPolicyName,omitempty"`
	EndTime                     *string                       `json:"endTime,omitempty"`
	Filters                     *[]string                     `json:"filters,omitempty"`
	StartTime                   *string                       `json:"startTime,omitempty"`
	StreamingLocatorId          *string                       `json:"streamingLocatorId,omitempty"`
	StreamingPolicyName         string                        `json:"streamingPolicyName"`
}

func (o *StreamingLocatorProperties) GetCreatedAsTime() (*time.Time, error) {
	if o.Created == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Created, "2006-01-02T15:04:05Z07:00")
}

func (o *StreamingLocatorProperties) SetCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Created = &formatted
}

func (o *StreamingLocatorProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *StreamingLocatorProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *StreamingLocatorProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *StreamingLocatorProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
