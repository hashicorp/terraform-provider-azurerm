package customdomains

import (
	"encoding/json"
	"fmt"
)

var _ CustomDomainHttpsParameters = CdnManagedHttpsParameters{}

type CdnManagedHttpsParameters struct {
	CertificateSourceParameters CdnCertificateSourceParameters `json:"certificateSourceParameters"`

	// Fields inherited from CustomDomainHttpsParameters
	MinimumTlsVersion *MinimumTlsVersion `json:"minimumTlsVersion,omitempty"`
	ProtocolType      ProtocolType       `json:"protocolType"`
}

var _ json.Marshaler = CdnManagedHttpsParameters{}

func (s CdnManagedHttpsParameters) MarshalJSON() ([]byte, error) {
	type wrapper CdnManagedHttpsParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CdnManagedHttpsParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CdnManagedHttpsParameters: %+v", err)
	}
	decoded["certificateSource"] = "Cdn"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CdnManagedHttpsParameters: %+v", err)
	}

	return encoded, nil
}
