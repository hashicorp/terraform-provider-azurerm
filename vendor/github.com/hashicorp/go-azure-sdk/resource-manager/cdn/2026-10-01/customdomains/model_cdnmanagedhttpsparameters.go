package customdomains

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CustomDomainHTTPSParameters = CdnManagedHTTPSParameters{}

type CdnManagedHTTPSParameters struct {
	CertificateSourceParameters CdnCertificateSourceParameters `json:"certificateSourceParameters"`

	// Fields inherited from CustomDomainHTTPSParameters

	CertificateSource CertificateSource  `json:"certificateSource"`
	MinimumTlsVersion *MinimumTlsVersion `json:"minimumTlsVersion,omitempty"`
	ProtocolType      ProtocolType       `json:"protocolType"`
}

func (s CdnManagedHTTPSParameters) CustomDomainHTTPSParameters() BaseCustomDomainHTTPSParametersImpl {
	return BaseCustomDomainHTTPSParametersImpl{
		CertificateSource: s.CertificateSource,
		MinimumTlsVersion: s.MinimumTlsVersion,
		ProtocolType:      s.ProtocolType,
	}
}

var _ json.Marshaler = CdnManagedHTTPSParameters{}

func (s CdnManagedHTTPSParameters) MarshalJSON() ([]byte, error) {
	type wrapper CdnManagedHTTPSParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CdnManagedHTTPSParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CdnManagedHTTPSParameters: %+v", err)
	}

	decoded["certificateSource"] = "Cdn"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CdnManagedHTTPSParameters: %+v", err)
	}

	return encoded, nil
}
