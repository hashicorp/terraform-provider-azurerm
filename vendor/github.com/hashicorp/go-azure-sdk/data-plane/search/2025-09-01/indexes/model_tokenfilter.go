package indexes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TokenFilter interface {
	TokenFilter() BaseTokenFilterImpl
}

var _ TokenFilter = BaseTokenFilterImpl{}

type BaseTokenFilterImpl struct {
	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s BaseTokenFilterImpl) TokenFilter() BaseTokenFilterImpl {
	return s
}

var _ TokenFilter = RawTokenFilterImpl{}

// RawTokenFilterImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawTokenFilterImpl struct {
	tokenFilter BaseTokenFilterImpl
	Type        string
	Values      map[string]interface{}
}

func (s RawTokenFilterImpl) TokenFilter() BaseTokenFilterImpl {
	return s.tokenFilter
}

func UnmarshalTokenFilterImplementation(input []byte) (TokenFilter, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TokenFilter into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["@odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.AsciiFoldingTokenFilter") {
		var out AsciiFoldingTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AsciiFoldingTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.CjkBigramTokenFilter") {
		var out CjkBigramTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CjkBigramTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.CommonGramTokenFilter") {
		var out CommonGramTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CommonGramTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.DictionaryDecompounderTokenFilter") {
		var out DictionaryDecompounderTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DictionaryDecompounderTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.EdgeNGramTokenFilter") {
		var out EdgeNGramTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EdgeNGramTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.EdgeNGramTokenFilterV2") {
		var out EdgeNGramTokenFilterV2
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EdgeNGramTokenFilterV2: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.ElisionTokenFilter") {
		var out ElisionTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ElisionTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.KeepTokenFilter") {
		var out KeepTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KeepTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.KeywordMarkerTokenFilter") {
		var out KeywordMarkerTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KeywordMarkerTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.LengthTokenFilter") {
		var out LengthTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LengthTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.LimitTokenFilter") {
		var out LimitTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LimitTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.NGramTokenFilter") {
		var out NGramTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NGramTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.NGramTokenFilterV2") {
		var out NGramTokenFilterV2
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NGramTokenFilterV2: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.PatternCaptureTokenFilter") {
		var out PatternCaptureTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PatternCaptureTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.PatternReplaceTokenFilter") {
		var out PatternReplaceTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PatternReplaceTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.PhoneticTokenFilter") {
		var out PhoneticTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PhoneticTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.ShingleTokenFilter") {
		var out ShingleTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ShingleTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.SnowballTokenFilter") {
		var out SnowballTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SnowballTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.StemmerOverrideTokenFilter") {
		var out StemmerOverrideTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StemmerOverrideTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.StemmerTokenFilter") {
		var out StemmerTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StemmerTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.StopwordsTokenFilter") {
		var out StopwordsTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StopwordsTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.SynonymTokenFilter") {
		var out SynonymTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SynonymTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.TruncateTokenFilter") {
		var out TruncateTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TruncateTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.UniqueTokenFilter") {
		var out UniqueTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UniqueTokenFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.WordDelimiterTokenFilter") {
		var out WordDelimiterTokenFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WordDelimiterTokenFilter: %+v", err)
		}
		return out, nil
	}

	var parent BaseTokenFilterImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseTokenFilterImpl: %+v", err)
	}

	return RawTokenFilterImpl{
		tokenFilter: parent,
		Type:        value,
		Values:      temp,
	}, nil

}
