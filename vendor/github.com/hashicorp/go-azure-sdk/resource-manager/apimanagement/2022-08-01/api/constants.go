package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiType string

const (
	ApiTypeGraphql   ApiType = "graphql"
	ApiTypeHTTP      ApiType = "http"
	ApiTypeSoap      ApiType = "soap"
	ApiTypeWebsocket ApiType = "websocket"
)

func PossibleValuesForApiType() []string {
	return []string{
		string(ApiTypeGraphql),
		string(ApiTypeHTTP),
		string(ApiTypeSoap),
		string(ApiTypeWebsocket),
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
		"graphql":   ApiTypeGraphql,
		"http":      ApiTypeHTTP,
		"soap":      ApiTypeSoap,
		"websocket": ApiTypeWebsocket,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApiType(input)
	return &out, nil
}

type BearerTokenSendingMethods string

const (
	BearerTokenSendingMethodsAuthorizationHeader BearerTokenSendingMethods = "authorizationHeader"
	BearerTokenSendingMethodsQuery               BearerTokenSendingMethods = "query"
)

func PossibleValuesForBearerTokenSendingMethods() []string {
	return []string{
		string(BearerTokenSendingMethodsAuthorizationHeader),
		string(BearerTokenSendingMethodsQuery),
	}
}

func (s *BearerTokenSendingMethods) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBearerTokenSendingMethods(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBearerTokenSendingMethods(input string) (*BearerTokenSendingMethods, error) {
	vals := map[string]BearerTokenSendingMethods{
		"authorizationheader": BearerTokenSendingMethodsAuthorizationHeader,
		"query":               BearerTokenSendingMethodsQuery,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BearerTokenSendingMethods(input)
	return &out, nil
}

type ContentFormat string

const (
	ContentFormatGraphqlNegativelink             ContentFormat = "graphql-link"
	ContentFormatOpenapi                         ContentFormat = "openapi"
	ContentFormatOpenapiNegativelink             ContentFormat = "openapi-link"
	ContentFormatOpenapiPositivejson             ContentFormat = "openapi+json"
	ContentFormatOpenapiPositivejsonNegativelink ContentFormat = "openapi+json-link"
	ContentFormatSwaggerNegativejson             ContentFormat = "swagger-json"
	ContentFormatSwaggerNegativelinkNegativejson ContentFormat = "swagger-link-json"
	ContentFormatWadlNegativelinkNegativejson    ContentFormat = "wadl-link-json"
	ContentFormatWadlNegativexml                 ContentFormat = "wadl-xml"
	ContentFormatWsdl                            ContentFormat = "wsdl"
	ContentFormatWsdlNegativelink                ContentFormat = "wsdl-link"
)

func PossibleValuesForContentFormat() []string {
	return []string{
		string(ContentFormatGraphqlNegativelink),
		string(ContentFormatOpenapi),
		string(ContentFormatOpenapiNegativelink),
		string(ContentFormatOpenapiPositivejson),
		string(ContentFormatOpenapiPositivejsonNegativelink),
		string(ContentFormatSwaggerNegativejson),
		string(ContentFormatSwaggerNegativelinkNegativejson),
		string(ContentFormatWadlNegativelinkNegativejson),
		string(ContentFormatWadlNegativexml),
		string(ContentFormatWsdl),
		string(ContentFormatWsdlNegativelink),
	}
}

func (s *ContentFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContentFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContentFormat(input string) (*ContentFormat, error) {
	vals := map[string]ContentFormat{
		"graphql-link":      ContentFormatGraphqlNegativelink,
		"openapi":           ContentFormatOpenapi,
		"openapi-link":      ContentFormatOpenapiNegativelink,
		"openapi+json":      ContentFormatOpenapiPositivejson,
		"openapi+json-link": ContentFormatOpenapiPositivejsonNegativelink,
		"swagger-json":      ContentFormatSwaggerNegativejson,
		"swagger-link-json": ContentFormatSwaggerNegativelinkNegativejson,
		"wadl-link-json":    ContentFormatWadlNegativelinkNegativejson,
		"wadl-xml":          ContentFormatWadlNegativexml,
		"wsdl":              ContentFormatWsdl,
		"wsdl-link":         ContentFormatWsdlNegativelink,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContentFormat(input)
	return &out, nil
}

type Protocol string

const (
	ProtocolHTTP  Protocol = "http"
	ProtocolHTTPS Protocol = "https"
	ProtocolWs    Protocol = "ws"
	ProtocolWss   Protocol = "wss"
)

func PossibleValuesForProtocol() []string {
	return []string{
		string(ProtocolHTTP),
		string(ProtocolHTTPS),
		string(ProtocolWs),
		string(ProtocolWss),
	}
}

func (s *Protocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProtocol(input string) (*Protocol, error) {
	vals := map[string]Protocol{
		"http":  ProtocolHTTP,
		"https": ProtocolHTTPS,
		"ws":    ProtocolWs,
		"wss":   ProtocolWss,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Protocol(input)
	return &out, nil
}

type SoapApiType string

const (
	SoapApiTypeGraphql   SoapApiType = "graphql"
	SoapApiTypeHTTP      SoapApiType = "http"
	SoapApiTypeSoap      SoapApiType = "soap"
	SoapApiTypeWebsocket SoapApiType = "websocket"
)

func PossibleValuesForSoapApiType() []string {
	return []string{
		string(SoapApiTypeGraphql),
		string(SoapApiTypeHTTP),
		string(SoapApiTypeSoap),
		string(SoapApiTypeWebsocket),
	}
}

func (s *SoapApiType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSoapApiType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSoapApiType(input string) (*SoapApiType, error) {
	vals := map[string]SoapApiType{
		"graphql":   SoapApiTypeGraphql,
		"http":      SoapApiTypeHTTP,
		"soap":      SoapApiTypeSoap,
		"websocket": SoapApiTypeWebsocket,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SoapApiType(input)
	return &out, nil
}

type TranslateRequiredQueryParametersConduct string

const (
	TranslateRequiredQueryParametersConductQuery    TranslateRequiredQueryParametersConduct = "query"
	TranslateRequiredQueryParametersConductTemplate TranslateRequiredQueryParametersConduct = "template"
)

func PossibleValuesForTranslateRequiredQueryParametersConduct() []string {
	return []string{
		string(TranslateRequiredQueryParametersConductQuery),
		string(TranslateRequiredQueryParametersConductTemplate),
	}
}

func (s *TranslateRequiredQueryParametersConduct) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTranslateRequiredQueryParametersConduct(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTranslateRequiredQueryParametersConduct(input string) (*TranslateRequiredQueryParametersConduct, error) {
	vals := map[string]TranslateRequiredQueryParametersConduct{
		"query":    TranslateRequiredQueryParametersConductQuery,
		"template": TranslateRequiredQueryParametersConductTemplate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TranslateRequiredQueryParametersConduct(input)
	return &out, nil
}

type VersioningScheme string

const (
	VersioningSchemeHeader  VersioningScheme = "Header"
	VersioningSchemeQuery   VersioningScheme = "Query"
	VersioningSchemeSegment VersioningScheme = "Segment"
)

func PossibleValuesForVersioningScheme() []string {
	return []string{
		string(VersioningSchemeHeader),
		string(VersioningSchemeQuery),
		string(VersioningSchemeSegment),
	}
}

func (s *VersioningScheme) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVersioningScheme(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVersioningScheme(input string) (*VersioningScheme, error) {
	vals := map[string]VersioningScheme{
		"header":  VersioningSchemeHeader,
		"query":   VersioningSchemeQuery,
		"segment": VersioningSchemeSegment,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VersioningScheme(input)
	return &out, nil
}
