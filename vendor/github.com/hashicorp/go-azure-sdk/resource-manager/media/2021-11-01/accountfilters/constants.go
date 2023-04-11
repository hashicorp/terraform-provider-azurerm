package accountfilters

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FilterTrackPropertyCompareOperation string

const (
	FilterTrackPropertyCompareOperationEqual    FilterTrackPropertyCompareOperation = "Equal"
	FilterTrackPropertyCompareOperationNotEqual FilterTrackPropertyCompareOperation = "NotEqual"
)

func PossibleValuesForFilterTrackPropertyCompareOperation() []string {
	return []string{
		string(FilterTrackPropertyCompareOperationEqual),
		string(FilterTrackPropertyCompareOperationNotEqual),
	}
}

func parseFilterTrackPropertyCompareOperation(input string) (*FilterTrackPropertyCompareOperation, error) {
	vals := map[string]FilterTrackPropertyCompareOperation{
		"equal":    FilterTrackPropertyCompareOperationEqual,
		"notequal": FilterTrackPropertyCompareOperationNotEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FilterTrackPropertyCompareOperation(input)
	return &out, nil
}

type FilterTrackPropertyType string

const (
	FilterTrackPropertyTypeBitrate  FilterTrackPropertyType = "Bitrate"
	FilterTrackPropertyTypeFourCC   FilterTrackPropertyType = "FourCC"
	FilterTrackPropertyTypeLanguage FilterTrackPropertyType = "Language"
	FilterTrackPropertyTypeName     FilterTrackPropertyType = "Name"
	FilterTrackPropertyTypeType     FilterTrackPropertyType = "Type"
	FilterTrackPropertyTypeUnknown  FilterTrackPropertyType = "Unknown"
)

func PossibleValuesForFilterTrackPropertyType() []string {
	return []string{
		string(FilterTrackPropertyTypeBitrate),
		string(FilterTrackPropertyTypeFourCC),
		string(FilterTrackPropertyTypeLanguage),
		string(FilterTrackPropertyTypeName),
		string(FilterTrackPropertyTypeType),
		string(FilterTrackPropertyTypeUnknown),
	}
}

func parseFilterTrackPropertyType(input string) (*FilterTrackPropertyType, error) {
	vals := map[string]FilterTrackPropertyType{
		"bitrate":  FilterTrackPropertyTypeBitrate,
		"fourcc":   FilterTrackPropertyTypeFourCC,
		"language": FilterTrackPropertyTypeLanguage,
		"name":     FilterTrackPropertyTypeName,
		"type":     FilterTrackPropertyTypeType,
		"unknown":  FilterTrackPropertyTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FilterTrackPropertyType(input)
	return &out, nil
}
