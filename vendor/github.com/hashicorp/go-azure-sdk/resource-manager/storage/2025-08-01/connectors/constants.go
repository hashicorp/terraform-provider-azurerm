package connectors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NativeDataSharingProvisioningState string

const (
	NativeDataSharingProvisioningStateAccepted  NativeDataSharingProvisioningState = "Accepted"
	NativeDataSharingProvisioningStateCanceled  NativeDataSharingProvisioningState = "Canceled"
	NativeDataSharingProvisioningStateCreating  NativeDataSharingProvisioningState = "Creating"
	NativeDataSharingProvisioningStateDeleting  NativeDataSharingProvisioningState = "Deleting"
	NativeDataSharingProvisioningStateFailed    NativeDataSharingProvisioningState = "Failed"
	NativeDataSharingProvisioningStateSucceeded NativeDataSharingProvisioningState = "Succeeded"
)

func PossibleValuesForNativeDataSharingProvisioningState() []string {
	return []string{
		string(NativeDataSharingProvisioningStateAccepted),
		string(NativeDataSharingProvisioningStateCanceled),
		string(NativeDataSharingProvisioningStateCreating),
		string(NativeDataSharingProvisioningStateDeleting),
		string(NativeDataSharingProvisioningStateFailed),
		string(NativeDataSharingProvisioningStateSucceeded),
	}
}

func (s *NativeDataSharingProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNativeDataSharingProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNativeDataSharingProvisioningState(input string) (*NativeDataSharingProvisioningState, error) {
	vals := map[string]NativeDataSharingProvisioningState{
		"accepted":  NativeDataSharingProvisioningStateAccepted,
		"canceled":  NativeDataSharingProvisioningStateCanceled,
		"creating":  NativeDataSharingProvisioningStateCreating,
		"deleting":  NativeDataSharingProvisioningStateDeleting,
		"failed":    NativeDataSharingProvisioningStateFailed,
		"succeeded": NativeDataSharingProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NativeDataSharingProvisioningState(input)
	return &out, nil
}

type StorageConnectorAuthType string

const (
	StorageConnectorAuthTypeManagedIdentity StorageConnectorAuthType = "ManagedIdentity"
)

func PossibleValuesForStorageConnectorAuthType() []string {
	return []string{
		string(StorageConnectorAuthTypeManagedIdentity),
	}
}

func (s *StorageConnectorAuthType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageConnectorAuthType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageConnectorAuthType(input string) (*StorageConnectorAuthType, error) {
	vals := map[string]StorageConnectorAuthType{
		"managedidentity": StorageConnectorAuthTypeManagedIdentity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageConnectorAuthType(input)
	return &out, nil
}

type StorageConnectorConnectionType string

const (
	StorageConnectorConnectionTypeDataShare StorageConnectorConnectionType = "DataShare"
)

func PossibleValuesForStorageConnectorConnectionType() []string {
	return []string{
		string(StorageConnectorConnectionTypeDataShare),
	}
}

func (s *StorageConnectorConnectionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageConnectorConnectionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageConnectorConnectionType(input string) (*StorageConnectorConnectionType, error) {
	vals := map[string]StorageConnectorConnectionType{
		"datashare": StorageConnectorConnectionTypeDataShare,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageConnectorConnectionType(input)
	return &out, nil
}

type StorageConnectorDataSourceType string

const (
	StorageConnectorDataSourceTypeAzureDataShare StorageConnectorDataSourceType = "Azure_DataShare"
)

func PossibleValuesForStorageConnectorDataSourceType() []string {
	return []string{
		string(StorageConnectorDataSourceTypeAzureDataShare),
	}
}

func (s *StorageConnectorDataSourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageConnectorDataSourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageConnectorDataSourceType(input string) (*StorageConnectorDataSourceType, error) {
	vals := map[string]StorageConnectorDataSourceType{
		"azure_datashare": StorageConnectorDataSourceTypeAzureDataShare,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageConnectorDataSourceType(input)
	return &out, nil
}

type StorageConnectorSourceType string

const (
	StorageConnectorSourceTypeDataShare StorageConnectorSourceType = "DataShare"
)

func PossibleValuesForStorageConnectorSourceType() []string {
	return []string{
		string(StorageConnectorSourceTypeDataShare),
	}
}

func (s *StorageConnectorSourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageConnectorSourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageConnectorSourceType(input string) (*StorageConnectorSourceType, error) {
	vals := map[string]StorageConnectorSourceType{
		"datashare": StorageConnectorSourceTypeDataShare,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageConnectorSourceType(input)
	return &out, nil
}

type StorageConnectorState string

const (
	StorageConnectorStateActive   StorageConnectorState = "Active"
	StorageConnectorStateInactive StorageConnectorState = "Inactive"
)

func PossibleValuesForStorageConnectorState() []string {
	return []string{
		string(StorageConnectorStateActive),
		string(StorageConnectorStateInactive),
	}
}

func (s *StorageConnectorState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageConnectorState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageConnectorState(input string) (*StorageConnectorState, error) {
	vals := map[string]StorageConnectorState{
		"active":   StorageConnectorStateActive,
		"inactive": StorageConnectorStateInactive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageConnectorState(input)
	return &out, nil
}
