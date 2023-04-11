package assetsandassetfilters

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetContainerPermission string

const (
	AssetContainerPermissionRead            AssetContainerPermission = "Read"
	AssetContainerPermissionReadWrite       AssetContainerPermission = "ReadWrite"
	AssetContainerPermissionReadWriteDelete AssetContainerPermission = "ReadWriteDelete"
)

func PossibleValuesForAssetContainerPermission() []string {
	return []string{
		string(AssetContainerPermissionRead),
		string(AssetContainerPermissionReadWrite),
		string(AssetContainerPermissionReadWriteDelete),
	}
}

func parseAssetContainerPermission(input string) (*AssetContainerPermission, error) {
	vals := map[string]AssetContainerPermission{
		"read":            AssetContainerPermissionRead,
		"readwrite":       AssetContainerPermissionReadWrite,
		"readwritedelete": AssetContainerPermissionReadWriteDelete,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssetContainerPermission(input)
	return &out, nil
}

type AssetStorageEncryptionFormat string

const (
	AssetStorageEncryptionFormatMediaStorageClientEncryption AssetStorageEncryptionFormat = "MediaStorageClientEncryption"
	AssetStorageEncryptionFormatNone                         AssetStorageEncryptionFormat = "None"
)

func PossibleValuesForAssetStorageEncryptionFormat() []string {
	return []string{
		string(AssetStorageEncryptionFormatMediaStorageClientEncryption),
		string(AssetStorageEncryptionFormatNone),
	}
}

func parseAssetStorageEncryptionFormat(input string) (*AssetStorageEncryptionFormat, error) {
	vals := map[string]AssetStorageEncryptionFormat{
		"mediastorageclientencryption": AssetStorageEncryptionFormatMediaStorageClientEncryption,
		"none":                         AssetStorageEncryptionFormatNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssetStorageEncryptionFormat(input)
	return &out, nil
}

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

type ProvisioningState string

const (
	ProvisioningStateFailed     ProvisioningState = "Failed"
	ProvisioningStateInProgress ProvisioningState = "InProgress"
	ProvisioningStateSucceeded  ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateFailed),
		string(ProvisioningStateInProgress),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"failed":     ProvisioningStateFailed,
		"inprogress": ProvisioningStateInProgress,
		"succeeded":  ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type Visibility string

const (
	VisibilityHidden  Visibility = "Hidden"
	VisibilityVisible Visibility = "Visible"
)

func PossibleValuesForVisibility() []string {
	return []string{
		string(VisibilityHidden),
		string(VisibilityVisible),
	}
}

func parseVisibility(input string) (*Visibility, error) {
	vals := map[string]Visibility{
		"hidden":  VisibilityHidden,
		"visible": VisibilityVisible,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Visibility(input)
	return &out, nil
}
