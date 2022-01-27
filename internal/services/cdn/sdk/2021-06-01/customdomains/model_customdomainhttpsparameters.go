package customdomains

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CustomDomainHttpsParameters interface {
}

func unmarshalCustomDomainHttpsParametersImplementation(input []byte) (CustomDomainHttpsParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomDomainHttpsParameters into map[string]interface: %+v", err)
	}

	value, ok := temp["certificateSource"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Cdn") {
		var out CdnManagedHttpsParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CdnManagedHttpsParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureKeyVault") {
		var out UserManagedHttpsParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UserManagedHttpsParameters: %+v", err)
		}
		return out, nil
	}

	type RawCustomDomainHttpsParametersImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawCustomDomainHttpsParametersImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
