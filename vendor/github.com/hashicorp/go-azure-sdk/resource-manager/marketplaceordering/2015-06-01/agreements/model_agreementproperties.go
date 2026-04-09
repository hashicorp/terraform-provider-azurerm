package agreements

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgreementProperties struct {
	Accepted          *bool   `json:"accepted,omitempty"`
	LicenseTextLink   *string `json:"licenseTextLink,omitempty"`
	Plan              *string `json:"plan,omitempty"`
	PrivacyPolicyLink *string `json:"privacyPolicyLink,omitempty"`
	Product           *string `json:"product,omitempty"`
	Publisher         *string `json:"publisher,omitempty"`
	RetrieveDatetime  *string `json:"retrieveDatetime,omitempty"`
	Signature         *string `json:"signature,omitempty"`
}

func (o *AgreementProperties) GetRetrieveDatetimeAsTime() (*time.Time, error) {
	if o.RetrieveDatetime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RetrieveDatetime, "2006-01-02T15:04:05Z07:00")
}

func (o *AgreementProperties) SetRetrieveDatetimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RetrieveDatetime = &formatted
}
