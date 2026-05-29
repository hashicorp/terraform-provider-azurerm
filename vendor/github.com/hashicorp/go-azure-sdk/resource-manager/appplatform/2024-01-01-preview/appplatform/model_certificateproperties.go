package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateProperties interface {
	CertificateProperties() BaseCertificatePropertiesImpl
}

var _ CertificateProperties = BaseCertificatePropertiesImpl{}

type BaseCertificatePropertiesImpl struct {
	ActivateDate      *string                               `json:"activateDate,omitempty"`
	DnsNames          *[]string                             `json:"dnsNames,omitempty"`
	ExpirationDate    *string                               `json:"expirationDate,omitempty"`
	IssuedDate        *string                               `json:"issuedDate,omitempty"`
	Issuer            *string                               `json:"issuer,omitempty"`
	ProvisioningState *CertificateResourceProvisioningState `json:"provisioningState,omitempty"`
	SubjectName       *string                               `json:"subjectName,omitempty"`
	Thumbprint        *string                               `json:"thumbprint,omitempty"`
	Type              string                                `json:"type"`
}

func (s BaseCertificatePropertiesImpl) CertificateProperties() BaseCertificatePropertiesImpl {
	return s
}

var _ CertificateProperties = RawCertificatePropertiesImpl{}

// RawCertificatePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawCertificatePropertiesImpl struct {
	certificateProperties BaseCertificatePropertiesImpl
	Type                  string
	Values                map[string]interface{}
}

func (s RawCertificatePropertiesImpl) CertificateProperties() BaseCertificatePropertiesImpl {
	return s.certificateProperties
}

func UnmarshalCertificatePropertiesImplementation(input []byte) (CertificateProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling CertificateProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ContentCertificate") {
		var out ContentCertificateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContentCertificateProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KeyVaultCertificate") {
		var out KeyVaultCertificateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KeyVaultCertificateProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseCertificatePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseCertificatePropertiesImpl: %+v", err)
	}

	return RawCertificatePropertiesImpl{
		certificateProperties: parent,
		Type:                  value,
		Values:                temp,
	}, nil

}
