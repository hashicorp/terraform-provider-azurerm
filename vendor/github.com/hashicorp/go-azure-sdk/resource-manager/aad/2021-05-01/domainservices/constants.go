package domainservices

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExternalAccess string

const (
	ExternalAccessDisabled ExternalAccess = "Disabled"
	ExternalAccessEnabled  ExternalAccess = "Enabled"
)

func PossibleValuesForExternalAccess() []string {
	return []string{
		string(ExternalAccessDisabled),
		string(ExternalAccessEnabled),
	}
}

func parseExternalAccess(input string) (*ExternalAccess, error) {
	vals := map[string]ExternalAccess{
		"disabled": ExternalAccessDisabled,
		"enabled":  ExternalAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExternalAccess(input)
	return &out, nil
}

type FilteredSync string

const (
	FilteredSyncDisabled FilteredSync = "Disabled"
	FilteredSyncEnabled  FilteredSync = "Enabled"
)

func PossibleValuesForFilteredSync() []string {
	return []string{
		string(FilteredSyncDisabled),
		string(FilteredSyncEnabled),
	}
}

func parseFilteredSync(input string) (*FilteredSync, error) {
	vals := map[string]FilteredSync{
		"disabled": FilteredSyncDisabled,
		"enabled":  FilteredSyncEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FilteredSync(input)
	return &out, nil
}

type KerberosArmoring string

const (
	KerberosArmoringDisabled KerberosArmoring = "Disabled"
	KerberosArmoringEnabled  KerberosArmoring = "Enabled"
)

func PossibleValuesForKerberosArmoring() []string {
	return []string{
		string(KerberosArmoringDisabled),
		string(KerberosArmoringEnabled),
	}
}

func parseKerberosArmoring(input string) (*KerberosArmoring, error) {
	vals := map[string]KerberosArmoring{
		"disabled": KerberosArmoringDisabled,
		"enabled":  KerberosArmoringEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KerberosArmoring(input)
	return &out, nil
}

type KerberosRc4Encryption string

const (
	KerberosRc4EncryptionDisabled KerberosRc4Encryption = "Disabled"
	KerberosRc4EncryptionEnabled  KerberosRc4Encryption = "Enabled"
)

func PossibleValuesForKerberosRc4Encryption() []string {
	return []string{
		string(KerberosRc4EncryptionDisabled),
		string(KerberosRc4EncryptionEnabled),
	}
}

func parseKerberosRc4Encryption(input string) (*KerberosRc4Encryption, error) {
	vals := map[string]KerberosRc4Encryption{
		"disabled": KerberosRc4EncryptionDisabled,
		"enabled":  KerberosRc4EncryptionEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KerberosRc4Encryption(input)
	return &out, nil
}

type Ldaps string

const (
	LdapsDisabled Ldaps = "Disabled"
	LdapsEnabled  Ldaps = "Enabled"
)

func PossibleValuesForLdaps() []string {
	return []string{
		string(LdapsDisabled),
		string(LdapsEnabled),
	}
}

func parseLdaps(input string) (*Ldaps, error) {
	vals := map[string]Ldaps{
		"disabled": LdapsDisabled,
		"enabled":  LdapsEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Ldaps(input)
	return &out, nil
}

type NotifyDcAdmins string

const (
	NotifyDcAdminsDisabled NotifyDcAdmins = "Disabled"
	NotifyDcAdminsEnabled  NotifyDcAdmins = "Enabled"
)

func PossibleValuesForNotifyDcAdmins() []string {
	return []string{
		string(NotifyDcAdminsDisabled),
		string(NotifyDcAdminsEnabled),
	}
}

func parseNotifyDcAdmins(input string) (*NotifyDcAdmins, error) {
	vals := map[string]NotifyDcAdmins{
		"disabled": NotifyDcAdminsDisabled,
		"enabled":  NotifyDcAdminsEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NotifyDcAdmins(input)
	return &out, nil
}

type NotifyGlobalAdmins string

const (
	NotifyGlobalAdminsDisabled NotifyGlobalAdmins = "Disabled"
	NotifyGlobalAdminsEnabled  NotifyGlobalAdmins = "Enabled"
)

func PossibleValuesForNotifyGlobalAdmins() []string {
	return []string{
		string(NotifyGlobalAdminsDisabled),
		string(NotifyGlobalAdminsEnabled),
	}
}

func parseNotifyGlobalAdmins(input string) (*NotifyGlobalAdmins, error) {
	vals := map[string]NotifyGlobalAdmins{
		"disabled": NotifyGlobalAdminsDisabled,
		"enabled":  NotifyGlobalAdminsEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NotifyGlobalAdmins(input)
	return &out, nil
}

type NtlmV1 string

const (
	NtlmV1Disabled NtlmV1 = "Disabled"
	NtlmV1Enabled  NtlmV1 = "Enabled"
)

func PossibleValuesForNtlmV1() []string {
	return []string{
		string(NtlmV1Disabled),
		string(NtlmV1Enabled),
	}
}

func parseNtlmV1(input string) (*NtlmV1, error) {
	vals := map[string]NtlmV1{
		"disabled": NtlmV1Disabled,
		"enabled":  NtlmV1Enabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NtlmV1(input)
	return &out, nil
}

type Status string

const (
	StatusFailure Status = "Failure"
	StatusNone    Status = "None"
	StatusOK      Status = "OK"
	StatusRunning Status = "Running"
	StatusSkipped Status = "Skipped"
	StatusWarning Status = "Warning"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusFailure),
		string(StatusNone),
		string(StatusOK),
		string(StatusRunning),
		string(StatusSkipped),
		string(StatusWarning),
	}
}

func parseStatus(input string) (*Status, error) {
	vals := map[string]Status{
		"failure": StatusFailure,
		"none":    StatusNone,
		"ok":      StatusOK,
		"running": StatusRunning,
		"skipped": StatusSkipped,
		"warning": StatusWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}

type SyncKerberosPasswords string

const (
	SyncKerberosPasswordsDisabled SyncKerberosPasswords = "Disabled"
	SyncKerberosPasswordsEnabled  SyncKerberosPasswords = "Enabled"
)

func PossibleValuesForSyncKerberosPasswords() []string {
	return []string{
		string(SyncKerberosPasswordsDisabled),
		string(SyncKerberosPasswordsEnabled),
	}
}

func parseSyncKerberosPasswords(input string) (*SyncKerberosPasswords, error) {
	vals := map[string]SyncKerberosPasswords{
		"disabled": SyncKerberosPasswordsDisabled,
		"enabled":  SyncKerberosPasswordsEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SyncKerberosPasswords(input)
	return &out, nil
}

type SyncNtlmPasswords string

const (
	SyncNtlmPasswordsDisabled SyncNtlmPasswords = "Disabled"
	SyncNtlmPasswordsEnabled  SyncNtlmPasswords = "Enabled"
)

func PossibleValuesForSyncNtlmPasswords() []string {
	return []string{
		string(SyncNtlmPasswordsDisabled),
		string(SyncNtlmPasswordsEnabled),
	}
}

func parseSyncNtlmPasswords(input string) (*SyncNtlmPasswords, error) {
	vals := map[string]SyncNtlmPasswords{
		"disabled": SyncNtlmPasswordsDisabled,
		"enabled":  SyncNtlmPasswordsEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SyncNtlmPasswords(input)
	return &out, nil
}

type SyncOnPremPasswords string

const (
	SyncOnPremPasswordsDisabled SyncOnPremPasswords = "Disabled"
	SyncOnPremPasswordsEnabled  SyncOnPremPasswords = "Enabled"
)

func PossibleValuesForSyncOnPremPasswords() []string {
	return []string{
		string(SyncOnPremPasswordsDisabled),
		string(SyncOnPremPasswordsEnabled),
	}
}

func parseSyncOnPremPasswords(input string) (*SyncOnPremPasswords, error) {
	vals := map[string]SyncOnPremPasswords{
		"disabled": SyncOnPremPasswordsDisabled,
		"enabled":  SyncOnPremPasswordsEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SyncOnPremPasswords(input)
	return &out, nil
}

type TlsV1 string

const (
	TlsV1Disabled TlsV1 = "Disabled"
	TlsV1Enabled  TlsV1 = "Enabled"
)

func PossibleValuesForTlsV1() []string {
	return []string{
		string(TlsV1Disabled),
		string(TlsV1Enabled),
	}
}

func parseTlsV1(input string) (*TlsV1, error) {
	vals := map[string]TlsV1{
		"disabled": TlsV1Disabled,
		"enabled":  TlsV1Enabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TlsV1(input)
	return &out, nil
}
