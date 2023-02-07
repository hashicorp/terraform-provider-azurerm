package inputs

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationMode string

const (
	AuthenticationModeConnectionString AuthenticationMode = "ConnectionString"
	AuthenticationModeMsi              AuthenticationMode = "Msi"
	AuthenticationModeUserToken        AuthenticationMode = "UserToken"
)

func PossibleValuesForAuthenticationMode() []string {
	return []string{
		string(AuthenticationModeConnectionString),
		string(AuthenticationModeMsi),
		string(AuthenticationModeUserToken),
	}
}

func parseAuthenticationMode(input string) (*AuthenticationMode, error) {
	vals := map[string]AuthenticationMode{
		"connectionstring": AuthenticationModeConnectionString,
		"msi":              AuthenticationModeMsi,
		"usertoken":        AuthenticationModeUserToken,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationMode(input)
	return &out, nil
}

type CompressionType string

const (
	CompressionTypeDeflate CompressionType = "Deflate"
	CompressionTypeGZip    CompressionType = "GZip"
	CompressionTypeNone    CompressionType = "None"
)

func PossibleValuesForCompressionType() []string {
	return []string{
		string(CompressionTypeDeflate),
		string(CompressionTypeGZip),
		string(CompressionTypeNone),
	}
}

func parseCompressionType(input string) (*CompressionType, error) {
	vals := map[string]CompressionType{
		"deflate": CompressionTypeDeflate,
		"gzip":    CompressionTypeGZip,
		"none":    CompressionTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CompressionType(input)
	return &out, nil
}

type Encoding string

const (
	EncodingUTFEight Encoding = "UTF8"
)

func PossibleValuesForEncoding() []string {
	return []string{
		string(EncodingUTFEight),
	}
}

func parseEncoding(input string) (*Encoding, error) {
	vals := map[string]Encoding{
		"utf8": EncodingUTFEight,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Encoding(input)
	return &out, nil
}

type EventSerializationType string

const (
	EventSerializationTypeAvro    EventSerializationType = "Avro"
	EventSerializationTypeCsv     EventSerializationType = "Csv"
	EventSerializationTypeJson    EventSerializationType = "Json"
	EventSerializationTypeParquet EventSerializationType = "Parquet"
)

func PossibleValuesForEventSerializationType() []string {
	return []string{
		string(EventSerializationTypeAvro),
		string(EventSerializationTypeCsv),
		string(EventSerializationTypeJson),
		string(EventSerializationTypeParquet),
	}
}

func parseEventSerializationType(input string) (*EventSerializationType, error) {
	vals := map[string]EventSerializationType{
		"avro":    EventSerializationTypeAvro,
		"csv":     EventSerializationTypeCsv,
		"json":    EventSerializationTypeJson,
		"parquet": EventSerializationTypeParquet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventSerializationType(input)
	return &out, nil
}

type JsonOutputSerializationFormat string

const (
	JsonOutputSerializationFormatArray         JsonOutputSerializationFormat = "Array"
	JsonOutputSerializationFormatLineSeparated JsonOutputSerializationFormat = "LineSeparated"
)

func PossibleValuesForJsonOutputSerializationFormat() []string {
	return []string{
		string(JsonOutputSerializationFormatArray),
		string(JsonOutputSerializationFormatLineSeparated),
	}
}

func parseJsonOutputSerializationFormat(input string) (*JsonOutputSerializationFormat, error) {
	vals := map[string]JsonOutputSerializationFormat{
		"array":         JsonOutputSerializationFormatArray,
		"lineseparated": JsonOutputSerializationFormatLineSeparated,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JsonOutputSerializationFormat(input)
	return &out, nil
}

type RefreshType string

const (
	RefreshTypeRefreshPeriodicallyWithDelta RefreshType = "RefreshPeriodicallyWithDelta"
	RefreshTypeRefreshPeriodicallyWithFull  RefreshType = "RefreshPeriodicallyWithFull"
	RefreshTypeStatic                       RefreshType = "Static"
)

func PossibleValuesForRefreshType() []string {
	return []string{
		string(RefreshTypeRefreshPeriodicallyWithDelta),
		string(RefreshTypeRefreshPeriodicallyWithFull),
		string(RefreshTypeStatic),
	}
}

func parseRefreshType(input string) (*RefreshType, error) {
	vals := map[string]RefreshType{
		"refreshperiodicallywithdelta": RefreshTypeRefreshPeriodicallyWithDelta,
		"refreshperiodicallywithfull":  RefreshTypeRefreshPeriodicallyWithFull,
		"static":                       RefreshTypeStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RefreshType(input)
	return &out, nil
}
