package customdomains

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CertificateSourceParameters = CdnCertificateSourceParameters{}

type CdnCertificateSourceParameters struct {
	CertificateType CertificateType `json:"certificateType"`

	// Fields inherited from CertificateSourceParameters

	TypeName CertificateSourceParametersType `json:"typeName"`
}

func (s CdnCertificateSourceParameters) CertificateSourceParameters() BaseCertificateSourceParametersImpl {
	return BaseCertificateSourceParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = CdnCertificateSourceParameters{}

func (s CdnCertificateSourceParameters) MarshalJSON() ([]byte, error) {
	type wrapper CdnCertificateSourceParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CdnCertificateSourceParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CdnCertificateSourceParameters: %+v", err)
	}

	decoded["typeName"] = "CdnCertificateSourceParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CdnCertificateSourceParameters: %+v", err)
	}

	return encoded, nil
}
