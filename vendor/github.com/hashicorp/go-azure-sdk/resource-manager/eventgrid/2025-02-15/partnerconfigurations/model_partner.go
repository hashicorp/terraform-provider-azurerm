package partnerconfigurations

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Partner struct {
	AuthorizationExpirationTimeInUtc *string `json:"authorizationExpirationTimeInUtc,omitempty"`
	PartnerName                      *string `json:"partnerName,omitempty"`
	PartnerRegistrationImmutableId   *string `json:"partnerRegistrationImmutableId,omitempty"`
}

func (o *Partner) GetAuthorizationExpirationTimeInUtcAsTime() (*time.Time, error) {
	if o.AuthorizationExpirationTimeInUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.AuthorizationExpirationTimeInUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *Partner) SetAuthorizationExpirationTimeInUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.AuthorizationExpirationTimeInUtc = &formatted
}
