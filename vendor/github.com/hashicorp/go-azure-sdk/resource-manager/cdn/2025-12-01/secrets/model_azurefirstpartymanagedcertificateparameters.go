package secrets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SecretParameters = AzureFirstPartyManagedCertificateParameters{}

type AzureFirstPartyManagedCertificateParameters struct {
	CertificateAuthority    *string            `json:"certificateAuthority,omitempty"`
	ExpirationDate          *string            `json:"expirationDate,omitempty"`
	SecretSource            *ResourceReference `json:"secretSource,omitempty"`
	Subject                 *string            `json:"subject,omitempty"`
	SubjectAlternativeNames *[]string          `json:"subjectAlternativeNames,omitempty"`
	Thumbprint              *string            `json:"thumbprint,omitempty"`

	// Fields inherited from SecretParameters

	Type SecretType `json:"type"`
}

func (s AzureFirstPartyManagedCertificateParameters) SecretParameters() BaseSecretParametersImpl {
	return BaseSecretParametersImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureFirstPartyManagedCertificateParameters{}

func (s AzureFirstPartyManagedCertificateParameters) MarshalJSON() ([]byte, error) {
	type wrapper AzureFirstPartyManagedCertificateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureFirstPartyManagedCertificateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureFirstPartyManagedCertificateParameters: %+v", err)
	}

	decoded["type"] = "AzureFirstPartyManagedCertificate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureFirstPartyManagedCertificateParameters: %+v", err)
	}

	return encoded, nil
}
