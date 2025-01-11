package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CertificateProperties = KeyVaultCertificateProperties{}

type KeyVaultCertificateProperties struct {
	AutoSync          *KeyVaultCertificateAutoSync `json:"autoSync,omitempty"`
	CertVersion       *string                      `json:"certVersion,omitempty"`
	ExcludePrivateKey *bool                        `json:"excludePrivateKey,omitempty"`
	KeyVaultCertName  string                       `json:"keyVaultCertName"`
	VaultUri          string                       `json:"vaultUri"`

	// Fields inherited from CertificateProperties

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

func (s KeyVaultCertificateProperties) CertificateProperties() BaseCertificatePropertiesImpl {
	return BaseCertificatePropertiesImpl{
		ActivateDate:      s.ActivateDate,
		DnsNames:          s.DnsNames,
		ExpirationDate:    s.ExpirationDate,
		IssuedDate:        s.IssuedDate,
		Issuer:            s.Issuer,
		ProvisioningState: s.ProvisioningState,
		SubjectName:       s.SubjectName,
		Thumbprint:        s.Thumbprint,
		Type:              s.Type,
	}
}

var _ json.Marshaler = KeyVaultCertificateProperties{}

func (s KeyVaultCertificateProperties) MarshalJSON() ([]byte, error) {
	type wrapper KeyVaultCertificateProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KeyVaultCertificateProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KeyVaultCertificateProperties: %+v", err)
	}

	decoded["type"] = "KeyVaultCertificate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KeyVaultCertificateProperties: %+v", err)
	}

	return encoded, nil
}
