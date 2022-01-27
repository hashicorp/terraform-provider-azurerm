package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ CustomDomainHttpsParameters = UserManagedHttpsParameters{}

type UserManagedHttpsParameters struct {
	CertificateSourceParameters KeyVaultCertificateSourceParameters `json:"certificateSourceParameters"`

	// Fields inherited from CustomDomainHttpsParameters
	MinimumTlsVersion *MinimumTlsVersion `json:"minimumTlsVersion,omitempty"`
	ProtocolType      ProtocolType       `json:"protocolType"`
}

var _ json.Marshaler = UserManagedHttpsParameters{}

func (s UserManagedHttpsParameters) MarshalJSON() ([]byte, error) {
	type wrapper UserManagedHttpsParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling UserManagedHttpsParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling UserManagedHttpsParameters: %+v", err)
	}
	decoded["certificateSource"] = "AzureKeyVault"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling UserManagedHttpsParameters: %+v", err)
	}

	return encoded, nil
}
