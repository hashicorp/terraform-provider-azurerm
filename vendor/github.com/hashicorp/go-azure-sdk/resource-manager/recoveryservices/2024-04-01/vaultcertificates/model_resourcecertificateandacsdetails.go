package vaultcertificates

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ResourceCertificateDetails = ResourceCertificateAndAcsDetails{}

type ResourceCertificateAndAcsDetails struct {
	GlobalAcsHostName  string `json:"globalAcsHostName"`
	GlobalAcsNamespace string `json:"globalAcsNamespace"`
	GlobalAcsRPRealm   string `json:"globalAcsRPRealm"`

	// Fields inherited from ResourceCertificateDetails

	AuthType     string  `json:"authType"`
	Certificate  *string `json:"certificate,omitempty"`
	FriendlyName *string `json:"friendlyName,omitempty"`
	Issuer       *string `json:"issuer,omitempty"`
	ResourceId   *int64  `json:"resourceId,omitempty"`
	Subject      *string `json:"subject,omitempty"`
	Thumbprint   *string `json:"thumbprint,omitempty"`
	ValidFrom    *string `json:"validFrom,omitempty"`
	ValidTo      *string `json:"validTo,omitempty"`
}

func (s ResourceCertificateAndAcsDetails) ResourceCertificateDetails() BaseResourceCertificateDetailsImpl {
	return BaseResourceCertificateDetailsImpl{
		AuthType:     s.AuthType,
		Certificate:  s.Certificate,
		FriendlyName: s.FriendlyName,
		Issuer:       s.Issuer,
		ResourceId:   s.ResourceId,
		Subject:      s.Subject,
		Thumbprint:   s.Thumbprint,
		ValidFrom:    s.ValidFrom,
		ValidTo:      s.ValidTo,
	}
}

func (o *ResourceCertificateAndAcsDetails) GetValidFromAsTime() (*time.Time, error) {
	if o.ValidFrom == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ValidFrom, "2006-01-02T15:04:05Z07:00")
}

func (o *ResourceCertificateAndAcsDetails) SetValidFromAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ValidFrom = &formatted
}

func (o *ResourceCertificateAndAcsDetails) GetValidToAsTime() (*time.Time, error) {
	if o.ValidTo == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ValidTo, "2006-01-02T15:04:05Z07:00")
}

func (o *ResourceCertificateAndAcsDetails) SetValidToAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ValidTo = &formatted
}

var _ json.Marshaler = ResourceCertificateAndAcsDetails{}

func (s ResourceCertificateAndAcsDetails) MarshalJSON() ([]byte, error) {
	type wrapper ResourceCertificateAndAcsDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ResourceCertificateAndAcsDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ResourceCertificateAndAcsDetails: %+v", err)
	}

	decoded["authType"] = "AccessControlService"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ResourceCertificateAndAcsDetails: %+v", err)
	}

	return encoded, nil
}
