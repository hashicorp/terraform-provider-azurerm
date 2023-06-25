package user

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserProperties struct {
	AdditionalUsageQuota *string            `json:"additionalUsageQuota,omitempty"`
	DisplayName          *string            `json:"displayName,omitempty"`
	Email                string             `json:"email"`
	InvitationSent       *string            `json:"invitationSent,omitempty"`
	InvitationState      *InvitationState   `json:"invitationState,omitempty"`
	ProvisioningState    *ProvisioningState `json:"provisioningState,omitempty"`
	RegistrationState    *RegistrationState `json:"registrationState,omitempty"`
	TotalUsage           *string            `json:"totalUsage,omitempty"`
}

func (o *UserProperties) GetInvitationSentAsTime() (*time.Time, error) {
	if o.InvitationSent == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.InvitationSent, "2006-01-02T15:04:05Z07:00")
}

func (o *UserProperties) SetInvitationSentAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.InvitationSent = &formatted
}
