package organizationresources

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrganizationResourceProperties struct {
	CreatedTime       *string           `json:"createdTime,omitempty"`
	LinkOrganization  *LinkOrganization `json:"linkOrganization,omitempty"`
	OfferDetail       OfferDetail       `json:"offerDetail"`
	OrganizationId    *string           `json:"organizationId,omitempty"`
	ProvisioningState *ProvisionState   `json:"provisioningState,omitempty"`
	SsoURL            *string           `json:"ssoUrl,omitempty"`
	UserDetail        UserDetail        `json:"userDetail"`
}

func (o *OrganizationResourceProperties) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *OrganizationResourceProperties) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}
