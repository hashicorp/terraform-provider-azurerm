package storagediscoveryworkspaces

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceProvisioningState string

const (
	ResourceProvisioningStateCanceled  ResourceProvisioningState = "Canceled"
	ResourceProvisioningStateFailed    ResourceProvisioningState = "Failed"
	ResourceProvisioningStateSucceeded ResourceProvisioningState = "Succeeded"
)

func PossibleValuesForResourceProvisioningState() []string {
	return []string{
		string(ResourceProvisioningStateCanceled),
		string(ResourceProvisioningStateFailed),
		string(ResourceProvisioningStateSucceeded),
	}
}

func (s *ResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceProvisioningState(input string) (*ResourceProvisioningState, error) {
	vals := map[string]ResourceProvisioningState{
		"canceled":  ResourceProvisioningStateCanceled,
		"failed":    ResourceProvisioningStateFailed,
		"succeeded": ResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceProvisioningState(input)
	return &out, nil
}

type StorageDiscoveryResourceType string

const (
	StorageDiscoveryResourceTypeMicrosoftPointStorageStorageAccounts StorageDiscoveryResourceType = "Microsoft.Storage/storageAccounts"
)

func PossibleValuesForStorageDiscoveryResourceType() []string {
	return []string{
		string(StorageDiscoveryResourceTypeMicrosoftPointStorageStorageAccounts),
	}
}

func (s *StorageDiscoveryResourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageDiscoveryResourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageDiscoveryResourceType(input string) (*StorageDiscoveryResourceType, error) {
	vals := map[string]StorageDiscoveryResourceType{
		"microsoft.storage/storageaccounts": StorageDiscoveryResourceTypeMicrosoftPointStorageStorageAccounts,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageDiscoveryResourceType(input)
	return &out, nil
}

type StorageDiscoverySku string

const (
	StorageDiscoverySkuFree     StorageDiscoverySku = "Free"
	StorageDiscoverySkuStandard StorageDiscoverySku = "Standard"
)

func PossibleValuesForStorageDiscoverySku() []string {
	return []string{
		string(StorageDiscoverySkuFree),
		string(StorageDiscoverySkuStandard),
	}
}

func (s *StorageDiscoverySku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageDiscoverySku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageDiscoverySku(input string) (*StorageDiscoverySku, error) {
	vals := map[string]StorageDiscoverySku{
		"free":     StorageDiscoverySkuFree,
		"standard": StorageDiscoverySkuStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageDiscoverySku(input)
	return &out, nil
}
