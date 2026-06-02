package certificate

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateProperties struct {
	DeleteCertificateError                  *DeleteCertificateError       `json:"deleteCertificateError,omitempty"`
	Format                                  *CertificateFormat            `json:"format,omitempty"`
	PreviousProvisioningState               *CertificateProvisioningState `json:"previousProvisioningState,omitempty"`
	PreviousProvisioningStateTransitionTime *string                       `json:"previousProvisioningStateTransitionTime,omitempty"`
	ProvisioningState                       *CertificateProvisioningState `json:"provisioningState,omitempty"`
	ProvisioningStateTransitionTime         *string                       `json:"provisioningStateTransitionTime,omitempty"`
	PublicData                              *string                       `json:"publicData,omitempty"`
	Thumbprint                              *string                       `json:"thumbprint,omitempty"`
	ThumbprintAlgorithm                     *string                       `json:"thumbprintAlgorithm,omitempty"`
}

func (o *CertificateProperties) GetPreviousProvisioningStateTransitionTimeAsTime() (*time.Time, error) {
	if o.PreviousProvisioningStateTransitionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PreviousProvisioningStateTransitionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *CertificateProperties) SetPreviousProvisioningStateTransitionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PreviousProvisioningStateTransitionTime = &formatted
}

func (o *CertificateProperties) GetProvisioningStateTransitionTimeAsTime() (*time.Time, error) {
	if o.ProvisioningStateTransitionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ProvisioningStateTransitionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *CertificateProperties) SetProvisioningStateTransitionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ProvisioningStateTransitionTime = &formatted
}
