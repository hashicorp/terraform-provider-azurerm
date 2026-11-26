package secrets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SecretParameters = CustomerCertificateParameters{}

type CustomerCertificateParameters struct {
	CertificateAuthority    *string           `json:"certificateAuthority,omitempty"`
	ExpirationDate          *string           `json:"expirationDate,omitempty"`
	SecretSource            ResourceReference `json:"secretSource"`
	SecretVersion           *string           `json:"secretVersion,omitempty"`
	Subject                 *string           `json:"subject,omitempty"`
	SubjectAlternativeNames *[]string         `json:"subjectAlternativeNames,omitempty"`
	Thumbprint              *string           `json:"thumbprint,omitempty"`
	UseLatestVersion        *bool             `json:"useLatestVersion,omitempty"`

	// Fields inherited from SecretParameters

	Type SecretType `json:"type"`
}

func (s CustomerCertificateParameters) SecretParameters() BaseSecretParametersImpl {
	return BaseSecretParametersImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = CustomerCertificateParameters{}

func (s CustomerCertificateParameters) MarshalJSON() ([]byte, error) {
	type wrapper CustomerCertificateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomerCertificateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomerCertificateParameters: %+v", err)
	}

	decoded["type"] = "CustomerCertificate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomerCertificateParameters: %+v", err)
	}

	return encoded, nil
}
