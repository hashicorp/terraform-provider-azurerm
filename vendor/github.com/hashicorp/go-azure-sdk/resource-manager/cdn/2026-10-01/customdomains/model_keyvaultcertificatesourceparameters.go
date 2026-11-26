package customdomains

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CertificateSourceParameters = KeyVaultCertificateSourceParameters{}

type KeyVaultCertificateSourceParameters struct {
	DeleteRule        DeleteRule `json:"deleteRule"`
	ResourceGroupName string     `json:"resourceGroupName"`
	SecretName        string     `json:"secretName"`
	SecretVersion     *string    `json:"secretVersion,omitempty"`
	SubscriptionId    string     `json:"subscriptionId"`
	UpdateRule        UpdateRule `json:"updateRule"`
	VaultName         string     `json:"vaultName"`

	// Fields inherited from CertificateSourceParameters

	TypeName CertificateSourceParametersType `json:"typeName"`
}

func (s KeyVaultCertificateSourceParameters) CertificateSourceParameters() BaseCertificateSourceParametersImpl {
	return BaseCertificateSourceParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = KeyVaultCertificateSourceParameters{}

func (s KeyVaultCertificateSourceParameters) MarshalJSON() ([]byte, error) {
	type wrapper KeyVaultCertificateSourceParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KeyVaultCertificateSourceParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KeyVaultCertificateSourceParameters: %+v", err)
	}

	decoded["typeName"] = "KeyVaultCertificateSourceParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KeyVaultCertificateSourceParameters: %+v", err)
	}

	return encoded, nil
}
