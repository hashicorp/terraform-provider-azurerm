package customdomains

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CustomDomainHTTPSParameters = UserManagedHTTPSParameters{}

type UserManagedHTTPSParameters struct {
	CertificateSourceParameters KeyVaultCertificateSourceParameters `json:"certificateSourceParameters"`

	// Fields inherited from CustomDomainHTTPSParameters

	CertificateSource CertificateSource  `json:"certificateSource"`
	MinimumTlsVersion *MinimumTlsVersion `json:"minimumTlsVersion,omitempty"`
	ProtocolType      ProtocolType       `json:"protocolType"`
}

func (s UserManagedHTTPSParameters) CustomDomainHTTPSParameters() BaseCustomDomainHTTPSParametersImpl {
	return BaseCustomDomainHTTPSParametersImpl{
		CertificateSource: s.CertificateSource,
		MinimumTlsVersion: s.MinimumTlsVersion,
		ProtocolType:      s.ProtocolType,
	}
}

var _ json.Marshaler = UserManagedHTTPSParameters{}

func (s UserManagedHTTPSParameters) MarshalJSON() ([]byte, error) {
	type wrapper UserManagedHTTPSParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling UserManagedHTTPSParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling UserManagedHTTPSParameters: %+v", err)
	}

	decoded["certificateSource"] = "AzureKeyVault"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling UserManagedHTTPSParameters: %+v", err)
	}

	return encoded, nil
}
