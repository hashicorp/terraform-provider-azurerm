package managedapis

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiType string

const (
	ApiTypeNotSpecified ApiType = "NotSpecified"
	ApiTypeRest         ApiType = "Rest"
	ApiTypeSoap         ApiType = "Soap"
)

func PossibleValuesForApiType() []string {
	return []string{
		string(ApiTypeNotSpecified),
		string(ApiTypeRest),
		string(ApiTypeSoap),
	}
}

func (s *ApiType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApiType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApiType(input string) (*ApiType, error) {
	vals := map[string]ApiType{
		"notspecified": ApiTypeNotSpecified,
		"rest":         ApiTypeRest,
		"soap":         ApiTypeSoap,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApiType(input)
	return &out, nil
}

type ConnectionParameterType string

const (
	ConnectionParameterTypeArray        ConnectionParameterType = "array"
	ConnectionParameterTypeBool         ConnectionParameterType = "bool"
	ConnectionParameterTypeConnection   ConnectionParameterType = "connection"
	ConnectionParameterTypeInt          ConnectionParameterType = "int"
	ConnectionParameterTypeOauthSetting ConnectionParameterType = "oauthSetting"
	ConnectionParameterTypeObject       ConnectionParameterType = "object"
	ConnectionParameterTypeSecureobject ConnectionParameterType = "secureobject"
	ConnectionParameterTypeSecurestring ConnectionParameterType = "securestring"
	ConnectionParameterTypeString       ConnectionParameterType = "string"
)

func PossibleValuesForConnectionParameterType() []string {
	return []string{
		string(ConnectionParameterTypeArray),
		string(ConnectionParameterTypeBool),
		string(ConnectionParameterTypeConnection),
		string(ConnectionParameterTypeInt),
		string(ConnectionParameterTypeOauthSetting),
		string(ConnectionParameterTypeObject),
		string(ConnectionParameterTypeSecureobject),
		string(ConnectionParameterTypeSecurestring),
		string(ConnectionParameterTypeString),
	}
}

func (s *ConnectionParameterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionParameterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionParameterType(input string) (*ConnectionParameterType, error) {
	vals := map[string]ConnectionParameterType{
		"array":        ConnectionParameterTypeArray,
		"bool":         ConnectionParameterTypeBool,
		"connection":   ConnectionParameterTypeConnection,
		"int":          ConnectionParameterTypeInt,
		"oauthsetting": ConnectionParameterTypeOauthSetting,
		"object":       ConnectionParameterTypeObject,
		"secureobject": ConnectionParameterTypeSecureobject,
		"securestring": ConnectionParameterTypeSecurestring,
		"string":       ConnectionParameterTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionParameterType(input)
	return &out, nil
}

type WsdlImportMethod string

const (
	WsdlImportMethodNotSpecified    WsdlImportMethod = "NotSpecified"
	WsdlImportMethodSoapPassThrough WsdlImportMethod = "SoapPassThrough"
	WsdlImportMethodSoapToRest      WsdlImportMethod = "SoapToRest"
)

func PossibleValuesForWsdlImportMethod() []string {
	return []string{
		string(WsdlImportMethodNotSpecified),
		string(WsdlImportMethodSoapPassThrough),
		string(WsdlImportMethodSoapToRest),
	}
}

func (s *WsdlImportMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWsdlImportMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWsdlImportMethod(input string) (*WsdlImportMethod, error) {
	vals := map[string]WsdlImportMethod{
		"notspecified":    WsdlImportMethodNotSpecified,
		"soappassthrough": WsdlImportMethodSoapPassThrough,
		"soaptorest":      WsdlImportMethodSoapToRest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WsdlImportMethod(input)
	return &out, nil
}
