package buckets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BucketPatchPermissions string

const (
	BucketPatchPermissionsReadOnly  BucketPatchPermissions = "ReadOnly"
	BucketPatchPermissionsReadWrite BucketPatchPermissions = "ReadWrite"
)

func PossibleValuesForBucketPatchPermissions() []string {
	return []string{
		string(BucketPatchPermissionsReadOnly),
		string(BucketPatchPermissionsReadWrite),
	}
}

func (s *BucketPatchPermissions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBucketPatchPermissions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBucketPatchPermissions(input string) (*BucketPatchPermissions, error) {
	vals := map[string]BucketPatchPermissions{
		"readonly":  BucketPatchPermissionsReadOnly,
		"readwrite": BucketPatchPermissionsReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BucketPatchPermissions(input)
	return &out, nil
}

type BucketPermissions string

const (
	BucketPermissionsReadOnly  BucketPermissions = "ReadOnly"
	BucketPermissionsReadWrite BucketPermissions = "ReadWrite"
)

func PossibleValuesForBucketPermissions() []string {
	return []string{
		string(BucketPermissionsReadOnly),
		string(BucketPermissionsReadWrite),
	}
}

func (s *BucketPermissions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBucketPermissions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBucketPermissions(input string) (*BucketPermissions, error) {
	vals := map[string]BucketPermissions{
		"readonly":  BucketPermissionsReadOnly,
		"readwrite": BucketPermissionsReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BucketPermissions(input)
	return &out, nil
}

type CredentialsStatus string

const (
	CredentialsStatusActive             CredentialsStatus = "Active"
	CredentialsStatusCredentialsExpired CredentialsStatus = "CredentialsExpired"
	CredentialsStatusNoCredentialsSet   CredentialsStatus = "NoCredentialsSet"
)

func PossibleValuesForCredentialsStatus() []string {
	return []string{
		string(CredentialsStatusActive),
		string(CredentialsStatusCredentialsExpired),
		string(CredentialsStatusNoCredentialsSet),
	}
}

func (s *CredentialsStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCredentialsStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCredentialsStatus(input string) (*CredentialsStatus, error) {
	vals := map[string]CredentialsStatus{
		"active":             CredentialsStatusActive,
		"credentialsexpired": CredentialsStatusCredentialsExpired,
		"nocredentialsset":   CredentialsStatusNoCredentialsSet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CredentialsStatus(input)
	return &out, nil
}

type NetAppProvisioningState string

const (
	NetAppProvisioningStateAccepted  NetAppProvisioningState = "Accepted"
	NetAppProvisioningStateCreating  NetAppProvisioningState = "Creating"
	NetAppProvisioningStateDeleting  NetAppProvisioningState = "Deleting"
	NetAppProvisioningStateFailed    NetAppProvisioningState = "Failed"
	NetAppProvisioningStateMoving    NetAppProvisioningState = "Moving"
	NetAppProvisioningStatePatching  NetAppProvisioningState = "Patching"
	NetAppProvisioningStateSucceeded NetAppProvisioningState = "Succeeded"
	NetAppProvisioningStateUpdating  NetAppProvisioningState = "Updating"
)

func PossibleValuesForNetAppProvisioningState() []string {
	return []string{
		string(NetAppProvisioningStateAccepted),
		string(NetAppProvisioningStateCreating),
		string(NetAppProvisioningStateDeleting),
		string(NetAppProvisioningStateFailed),
		string(NetAppProvisioningStateMoving),
		string(NetAppProvisioningStatePatching),
		string(NetAppProvisioningStateSucceeded),
		string(NetAppProvisioningStateUpdating),
	}
}

func (s *NetAppProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetAppProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetAppProvisioningState(input string) (*NetAppProvisioningState, error) {
	vals := map[string]NetAppProvisioningState{
		"accepted":  NetAppProvisioningStateAccepted,
		"creating":  NetAppProvisioningStateCreating,
		"deleting":  NetAppProvisioningStateDeleting,
		"failed":    NetAppProvisioningStateFailed,
		"moving":    NetAppProvisioningStateMoving,
		"patching":  NetAppProvisioningStatePatching,
		"succeeded": NetAppProvisioningStateSucceeded,
		"updating":  NetAppProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetAppProvisioningState(input)
	return &out, nil
}

type OnCertificateConflictAction string

const (
	OnCertificateConflictActionFail   OnCertificateConflictAction = "Fail"
	OnCertificateConflictActionUpdate OnCertificateConflictAction = "Update"
)

func PossibleValuesForOnCertificateConflictAction() []string {
	return []string{
		string(OnCertificateConflictActionFail),
		string(OnCertificateConflictActionUpdate),
	}
}

func (s *OnCertificateConflictAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOnCertificateConflictAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOnCertificateConflictAction(input string) (*OnCertificateConflictAction, error) {
	vals := map[string]OnCertificateConflictAction{
		"fail":   OnCertificateConflictActionFail,
		"update": OnCertificateConflictActionUpdate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OnCertificateConflictAction(input)
	return &out, nil
}
