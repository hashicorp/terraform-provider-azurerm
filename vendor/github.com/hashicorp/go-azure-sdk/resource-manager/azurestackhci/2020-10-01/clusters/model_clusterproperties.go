package clusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProperties struct {
	AadClientId           string                     `json:"aadClientId"`
	AadTenantId           string                     `json:"aadTenantId"`
	BillingModel          *string                    `json:"billingModel,omitempty"`
	CloudId               *string                    `json:"cloudId,omitempty"`
	LastBillingTimestamp  *string                    `json:"lastBillingTimestamp,omitempty"`
	LastSyncTimestamp     *string                    `json:"lastSyncTimestamp,omitempty"`
	ProvisioningState     *ProvisioningState         `json:"provisioningState,omitempty"`
	RegistrationTimestamp *string                    `json:"registrationTimestamp,omitempty"`
	ReportedProperties    *ClusterReportedProperties `json:"reportedProperties,omitempty"`
	Status                *Status                    `json:"status,omitempty"`
	TrialDaysRemaining    *float64                   `json:"trialDaysRemaining,omitempty"`
}

func (o *ClusterProperties) GetLastBillingTimestampAsTime() (*time.Time, error) {
	if o.LastBillingTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastBillingTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ClusterProperties) SetLastBillingTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastBillingTimestamp = &formatted
}

func (o *ClusterProperties) GetLastSyncTimestampAsTime() (*time.Time, error) {
	if o.LastSyncTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSyncTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ClusterProperties) SetLastSyncTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSyncTimestamp = &formatted
}

func (o *ClusterProperties) GetRegistrationTimestampAsTime() (*time.Time, error) {
	if o.RegistrationTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RegistrationTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ClusterProperties) SetRegistrationTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RegistrationTimestamp = &formatted
}
