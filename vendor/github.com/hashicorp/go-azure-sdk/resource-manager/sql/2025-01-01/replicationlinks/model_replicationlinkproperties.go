package replicationlinks

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationLinkProperties struct {
	IsTerminationAllowed *bool                `json:"isTerminationAllowed,omitempty"`
	LinkType             *ReplicationLinkType `json:"linkType,omitempty"`
	PartnerDatabase      *string              `json:"partnerDatabase,omitempty"`
	PartnerDatabaseId    *string              `json:"partnerDatabaseId,omitempty"`
	PartnerLocation      *string              `json:"partnerLocation,omitempty"`
	PartnerRole          *ReplicationRole     `json:"partnerRole,omitempty"`
	PartnerServer        *string              `json:"partnerServer,omitempty"`
	PercentComplete      *int64               `json:"percentComplete,omitempty"`
	ReplicationMode      *string              `json:"replicationMode,omitempty"`
	ReplicationState     *ReplicationState    `json:"replicationState,omitempty"`
	Role                 *ReplicationRole     `json:"role,omitempty"`
	StartTime            *string              `json:"startTime,omitempty"`
}

func (o *ReplicationLinkProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ReplicationLinkProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
