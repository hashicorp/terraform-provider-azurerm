package appservicecertificateorders

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServiceCertificateOrderProperties struct {
	AppServiceCertificateNotRenewableReasons *[]ResourceNotRenewableReason     `json:"appServiceCertificateNotRenewableReasons,omitempty"`
	AutoRenew                                *bool                             `json:"autoRenew,omitempty"`
	Certificates                             *map[string]AppServiceCertificate `json:"certificates,omitempty"`
	Contact                                  *CertificateOrderContact          `json:"contact,omitempty"`
	Csr                                      *string                           `json:"csr,omitempty"`
	DistinguishedName                        *string                           `json:"distinguishedName,omitempty"`
	DomainVerificationToken                  *string                           `json:"domainVerificationToken,omitempty"`
	ExpirationTime                           *string                           `json:"expirationTime,omitempty"`
	Intermediate                             *CertificateDetails               `json:"intermediate,omitempty"`
	IsPrivateKeyExternal                     *bool                             `json:"isPrivateKeyExternal,omitempty"`
	KeySize                                  *int64                            `json:"keySize,omitempty"`
	LastCertificateIssuanceTime              *string                           `json:"lastCertificateIssuanceTime,omitempty"`
	NextAutoRenewalTimeStamp                 *string                           `json:"nextAutoRenewalTimeStamp,omitempty"`
	ProductType                              CertificateProductType            `json:"productType"`
	ProvisioningState                        *ProvisioningState                `json:"provisioningState,omitempty"`
	Root                                     *CertificateDetails               `json:"root,omitempty"`
	SerialNumber                             *string                           `json:"serialNumber,omitempty"`
	SignedCertificate                        *CertificateDetails               `json:"signedCertificate,omitempty"`
	Status                                   *CertificateOrderStatus           `json:"status,omitempty"`
	ValidityInYears                          *int64                            `json:"validityInYears,omitempty"`
}

func (o *AppServiceCertificateOrderProperties) GetExpirationTimeAsTime() (*time.Time, error) {
	if o.ExpirationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AppServiceCertificateOrderProperties) SetExpirationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationTime = &formatted
}

func (o *AppServiceCertificateOrderProperties) GetLastCertificateIssuanceTimeAsTime() (*time.Time, error) {
	if o.LastCertificateIssuanceTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastCertificateIssuanceTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AppServiceCertificateOrderProperties) SetLastCertificateIssuanceTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastCertificateIssuanceTime = &formatted
}

func (o *AppServiceCertificateOrderProperties) GetNextAutoRenewalTimeStampAsTime() (*time.Time, error) {
	if o.NextAutoRenewalTimeStamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NextAutoRenewalTimeStamp, "2006-01-02T15:04:05Z07:00")
}

func (o *AppServiceCertificateOrderProperties) SetNextAutoRenewalTimeStampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NextAutoRenewalTimeStamp = &formatted
}
