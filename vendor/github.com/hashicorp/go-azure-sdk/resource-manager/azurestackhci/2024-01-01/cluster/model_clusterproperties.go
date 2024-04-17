package cluster

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProperties struct {
	AadApplicationObjectId             *string                             `json:"aadApplicationObjectId,omitempty"`
	AadClientId                        *string                             `json:"aadClientId,omitempty"`
	AadServicePrincipalObjectId        *string                             `json:"aadServicePrincipalObjectId,omitempty"`
	AadTenantId                        *string                             `json:"aadTenantId,omitempty"`
	BillingModel                       *string                             `json:"billingModel,omitempty"`
	CloudId                            *string                             `json:"cloudId,omitempty"`
	CloudManagementEndpoint            *string                             `json:"cloudManagementEndpoint,omitempty"`
	ConnectivityStatus                 *ConnectivityStatus                 `json:"connectivityStatus,omitempty"`
	DesiredProperties                  *ClusterDesiredProperties           `json:"desiredProperties,omitempty"`
	IsolatedVMAttestationConfiguration *IsolatedVMAttestationConfiguration `json:"isolatedVmAttestationConfiguration,omitempty"`
	LastBillingTimestamp               *string                             `json:"lastBillingTimestamp,omitempty"`
	LastSyncTimestamp                  *string                             `json:"lastSyncTimestamp,omitempty"`
	ProvisioningState                  *ProvisioningState                  `json:"provisioningState,omitempty"`
	RegistrationTimestamp              *string                             `json:"registrationTimestamp,omitempty"`
	ReportedProperties                 *ClusterReportedProperties          `json:"reportedProperties,omitempty"`
	ResourceProviderObjectId           *string                             `json:"resourceProviderObjectId,omitempty"`
	ServiceEndpoint                    *string                             `json:"serviceEndpoint,omitempty"`
	SoftwareAssuranceProperties        *SoftwareAssuranceProperties        `json:"softwareAssuranceProperties,omitempty"`
	Status                             *Status                             `json:"status,omitempty"`
	TrialDaysRemaining                 *float64                            `json:"trialDaysRemaining,omitempty"`
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
