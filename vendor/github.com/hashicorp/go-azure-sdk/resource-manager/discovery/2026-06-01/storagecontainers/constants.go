package storagecontainers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobStorageMountProtocol string

const (
	BlobStorageMountProtocolBlobfuseCaching BlobStorageMountProtocol = "BlobfuseCaching"
	BlobStorageMountProtocolNFS             BlobStorageMountProtocol = "NFS"
)

func PossibleValuesForBlobStorageMountProtocol() []string {
	return []string{
		string(BlobStorageMountProtocolBlobfuseCaching),
		string(BlobStorageMountProtocolNFS),
	}
}

func (s *BlobStorageMountProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBlobStorageMountProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBlobStorageMountProtocol(input string) (*BlobStorageMountProtocol, error) {
	vals := map[string]BlobStorageMountProtocol{
		"blobfusecaching": BlobStorageMountProtocolBlobfuseCaching,
		"nfs":             BlobStorageMountProtocolNFS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BlobStorageMountProtocol(input)
	return &out, nil
}

type NetAppMountProtocol string

const (
	NetAppMountProtocolNFS NetAppMountProtocol = "NFS"
)

func PossibleValuesForNetAppMountProtocol() []string {
	return []string{
		string(NetAppMountProtocolNFS),
	}
}

func (s *NetAppMountProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetAppMountProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetAppMountProtocol(input string) (*NetAppMountProtocol, error) {
	vals := map[string]NetAppMountProtocol{
		"nfs": NetAppMountProtocolNFS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetAppMountProtocol(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type StorageStoreType string

const (
	StorageStoreTypeAzureNetAppFiles StorageStoreType = "AzureNetAppFiles"
	StorageStoreTypeAzureStorageBlob StorageStoreType = "AzureStorageBlob"
)

func PossibleValuesForStorageStoreType() []string {
	return []string{
		string(StorageStoreTypeAzureNetAppFiles),
		string(StorageStoreTypeAzureStorageBlob),
	}
}

func (s *StorageStoreType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageStoreType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageStoreType(input string) (*StorageStoreType, error) {
	vals := map[string]StorageStoreType{
		"azurenetappfiles": StorageStoreTypeAzureNetAppFiles,
		"azurestorageblob": StorageStoreTypeAzureStorageBlob,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageStoreType(input)
	return &out, nil
}
