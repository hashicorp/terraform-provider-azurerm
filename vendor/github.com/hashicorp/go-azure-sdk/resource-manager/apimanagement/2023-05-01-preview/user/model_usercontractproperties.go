package user

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserContractProperties struct {
	Email            *string                    `json:"email,omitempty"`
	FirstName        *string                    `json:"firstName,omitempty"`
	Groups           *[]GroupContractProperties `json:"groups,omitempty"`
	Identities       *[]UserIdentityContract    `json:"identities,omitempty"`
	LastName         *string                    `json:"lastName,omitempty"`
	Note             *string                    `json:"note,omitempty"`
	RegistrationDate *string                    `json:"registrationDate,omitempty"`
	State            *UserState                 `json:"state,omitempty"`
}

func (o *UserContractProperties) GetRegistrationDateAsTime() (*time.Time, error) {
	if o.RegistrationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RegistrationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *UserContractProperties) SetRegistrationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RegistrationDate = &formatted
}
