package secrets

import (
	"encoding/json"
	"fmt"
)

var _ SecretParameters = AzureFirstPartyManagedCertificateParameters{}

type AzureFirstPartyManagedCertificateParameters struct {

	// Fields inherited from SecretParameters
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
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureFirstPartyManagedCertificateParameters: %+v", err)
	}
	decoded["type"] = "AzureFirstPartyManagedCertificate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureFirstPartyManagedCertificateParameters: %+v", err)
	}

	return encoded, nil
}
