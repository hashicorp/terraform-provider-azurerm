package secrets

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SecretParameters interface {
}

func unmarshalSecretParametersImplementation(input []byte) (SecretParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SecretParameters into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureFirstPartyManagedCertificate") {
		var out AzureFirstPartyManagedCertificateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFirstPartyManagedCertificateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CustomerCertificate") {
		var out CustomerCertificateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomerCertificateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ManagedCertificate") {
		var out ManagedCertificateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ManagedCertificateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "UrlSigningKey") {
		var out UrlSigningKeyParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UrlSigningKeyParameters: %+v", err)
		}
		return out, nil
	}

	type RawSecretParametersImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSecretParametersImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
