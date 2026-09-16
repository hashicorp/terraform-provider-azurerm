package certificateprofiles

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateProfileStatus string

const (
	CertificateProfileStatusActive    CertificateProfileStatus = "Active"
	CertificateProfileStatusDisabled  CertificateProfileStatus = "Disabled"
	CertificateProfileStatusSuspended CertificateProfileStatus = "Suspended"
)

func PossibleValuesForCertificateProfileStatus() []string {
	return []string{
		string(CertificateProfileStatusActive),
		string(CertificateProfileStatusDisabled),
		string(CertificateProfileStatusSuspended),
	}
}

func (s *CertificateProfileStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateProfileStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificateProfileStatus(input string) (*CertificateProfileStatus, error) {
	vals := map[string]CertificateProfileStatus{
		"active":    CertificateProfileStatusActive,
		"disabled":  CertificateProfileStatusDisabled,
		"suspended": CertificateProfileStatusSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateProfileStatus(input)
	return &out, nil
}

type CertificateStatus string

const (
	CertificateStatusActive  CertificateStatus = "Active"
	CertificateStatusExpired CertificateStatus = "Expired"
	CertificateStatusRevoked CertificateStatus = "Revoked"
)

func PossibleValuesForCertificateStatus() []string {
	return []string{
		string(CertificateStatusActive),
		string(CertificateStatusExpired),
		string(CertificateStatusRevoked),
	}
}

func (s *CertificateStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificateStatus(input string) (*CertificateStatus, error) {
	vals := map[string]CertificateStatus{
		"active":  CertificateStatusActive,
		"expired": CertificateStatusExpired,
		"revoked": CertificateStatusRevoked,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateStatus(input)
	return &out, nil
}

type ProfileType string

const (
	ProfileTypePrivateTrust         ProfileType = "PrivateTrust"
	ProfileTypePrivateTrustCIPolicy ProfileType = "PrivateTrustCIPolicy"
	ProfileTypePublicTrust          ProfileType = "PublicTrust"
	ProfileTypePublicTrustTest      ProfileType = "PublicTrustTest"
	ProfileTypeVBSEnclave           ProfileType = "VBSEnclave"
)

func PossibleValuesForProfileType() []string {
	return []string{
		string(ProfileTypePrivateTrust),
		string(ProfileTypePrivateTrustCIPolicy),
		string(ProfileTypePublicTrust),
		string(ProfileTypePublicTrustTest),
		string(ProfileTypeVBSEnclave),
	}
}

func (s *ProfileType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProfileType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProfileType(input string) (*ProfileType, error) {
	vals := map[string]ProfileType{
		"privatetrust":         ProfileTypePrivateTrust,
		"privatetrustcipolicy": ProfileTypePrivateTrustCIPolicy,
		"publictrust":          ProfileTypePublicTrust,
		"publictrusttest":      ProfileTypePublicTrustTest,
		"vbsenclave":           ProfileTypeVBSEnclave,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProfileType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
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
		"accepted":  ProvisioningStateAccepted,
		"canceled":  ProvisioningStateCanceled,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type RevocationStatus string

const (
	RevocationStatusFailed     RevocationStatus = "Failed"
	RevocationStatusInProgress RevocationStatus = "InProgress"
	RevocationStatusSucceeded  RevocationStatus = "Succeeded"
)

func PossibleValuesForRevocationStatus() []string {
	return []string{
		string(RevocationStatusFailed),
		string(RevocationStatusInProgress),
		string(RevocationStatusSucceeded),
	}
}

func (s *RevocationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRevocationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRevocationStatus(input string) (*RevocationStatus, error) {
	vals := map[string]RevocationStatus{
		"failed":     RevocationStatusFailed,
		"inprogress": RevocationStatusInProgress,
		"succeeded":  RevocationStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RevocationStatus(input)
	return &out, nil
}
