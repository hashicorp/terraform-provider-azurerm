package caches

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainJoinedType string

const (
	DomainJoinedTypeError DomainJoinedType = "Error"
	DomainJoinedTypeNo    DomainJoinedType = "No"
	DomainJoinedTypeYes   DomainJoinedType = "Yes"
)

func PossibleValuesForDomainJoinedType() []string {
	return []string{
		string(DomainJoinedTypeError),
		string(DomainJoinedTypeNo),
		string(DomainJoinedTypeYes),
	}
}

func parseDomainJoinedType(input string) (*DomainJoinedType, error) {
	vals := map[string]DomainJoinedType{
		"error": DomainJoinedTypeError,
		"no":    DomainJoinedTypeNo,
		"yes":   DomainJoinedTypeYes,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DomainJoinedType(input)
	return &out, nil
}

type FirmwareStatusType string

const (
	FirmwareStatusTypeAvailable   FirmwareStatusType = "available"
	FirmwareStatusTypeUnavailable FirmwareStatusType = "unavailable"
)

func PossibleValuesForFirmwareStatusType() []string {
	return []string{
		string(FirmwareStatusTypeAvailable),
		string(FirmwareStatusTypeUnavailable),
	}
}

func parseFirmwareStatusType(input string) (*FirmwareStatusType, error) {
	vals := map[string]FirmwareStatusType{
		"available":   FirmwareStatusTypeAvailable,
		"unavailable": FirmwareStatusTypeUnavailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirmwareStatusType(input)
	return &out, nil
}

type HealthStateType string

const (
	HealthStateTypeDegraded      HealthStateType = "Degraded"
	HealthStateTypeDown          HealthStateType = "Down"
	HealthStateTypeFlushing      HealthStateType = "Flushing"
	HealthStateTypeHealthy       HealthStateType = "Healthy"
	HealthStateTypeStopped       HealthStateType = "Stopped"
	HealthStateTypeStopping      HealthStateType = "Stopping"
	HealthStateTypeTransitioning HealthStateType = "Transitioning"
	HealthStateTypeUnknown       HealthStateType = "Unknown"
	HealthStateTypeUpgrading     HealthStateType = "Upgrading"
)

func PossibleValuesForHealthStateType() []string {
	return []string{
		string(HealthStateTypeDegraded),
		string(HealthStateTypeDown),
		string(HealthStateTypeFlushing),
		string(HealthStateTypeHealthy),
		string(HealthStateTypeStopped),
		string(HealthStateTypeStopping),
		string(HealthStateTypeTransitioning),
		string(HealthStateTypeUnknown),
		string(HealthStateTypeUpgrading),
	}
}

func parseHealthStateType(input string) (*HealthStateType, error) {
	vals := map[string]HealthStateType{
		"degraded":      HealthStateTypeDegraded,
		"down":          HealthStateTypeDown,
		"flushing":      HealthStateTypeFlushing,
		"healthy":       HealthStateTypeHealthy,
		"stopped":       HealthStateTypeStopped,
		"stopping":      HealthStateTypeStopping,
		"transitioning": HealthStateTypeTransitioning,
		"unknown":       HealthStateTypeUnknown,
		"upgrading":     HealthStateTypeUpgrading,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthStateType(input)
	return &out, nil
}

type NfsAccessRuleAccess string

const (
	NfsAccessRuleAccessNo NfsAccessRuleAccess = "no"
	NfsAccessRuleAccessRo NfsAccessRuleAccess = "ro"
	NfsAccessRuleAccessRw NfsAccessRuleAccess = "rw"
)

func PossibleValuesForNfsAccessRuleAccess() []string {
	return []string{
		string(NfsAccessRuleAccessNo),
		string(NfsAccessRuleAccessRo),
		string(NfsAccessRuleAccessRw),
	}
}

func parseNfsAccessRuleAccess(input string) (*NfsAccessRuleAccess, error) {
	vals := map[string]NfsAccessRuleAccess{
		"no": NfsAccessRuleAccessNo,
		"ro": NfsAccessRuleAccessRo,
		"rw": NfsAccessRuleAccessRw,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NfsAccessRuleAccess(input)
	return &out, nil
}

type NfsAccessRuleScope string

const (
	NfsAccessRuleScopeDefault NfsAccessRuleScope = "default"
	NfsAccessRuleScopeHost    NfsAccessRuleScope = "host"
	NfsAccessRuleScopeNetwork NfsAccessRuleScope = "network"
)

func PossibleValuesForNfsAccessRuleScope() []string {
	return []string{
		string(NfsAccessRuleScopeDefault),
		string(NfsAccessRuleScopeHost),
		string(NfsAccessRuleScopeNetwork),
	}
}

func parseNfsAccessRuleScope(input string) (*NfsAccessRuleScope, error) {
	vals := map[string]NfsAccessRuleScope{
		"default": NfsAccessRuleScopeDefault,
		"host":    NfsAccessRuleScopeHost,
		"network": NfsAccessRuleScopeNetwork,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NfsAccessRuleScope(input)
	return &out, nil
}

type ProvisioningStateType string

const (
	ProvisioningStateTypeCancelled ProvisioningStateType = "Cancelled"
	ProvisioningStateTypeCreating  ProvisioningStateType = "Creating"
	ProvisioningStateTypeDeleting  ProvisioningStateType = "Deleting"
	ProvisioningStateTypeFailed    ProvisioningStateType = "Failed"
	ProvisioningStateTypeSucceeded ProvisioningStateType = "Succeeded"
	ProvisioningStateTypeUpdating  ProvisioningStateType = "Updating"
)

func PossibleValuesForProvisioningStateType() []string {
	return []string{
		string(ProvisioningStateTypeCancelled),
		string(ProvisioningStateTypeCreating),
		string(ProvisioningStateTypeDeleting),
		string(ProvisioningStateTypeFailed),
		string(ProvisioningStateTypeSucceeded),
		string(ProvisioningStateTypeUpdating),
	}
}

func parseProvisioningStateType(input string) (*ProvisioningStateType, error) {
	vals := map[string]ProvisioningStateType{
		"cancelled": ProvisioningStateTypeCancelled,
		"creating":  ProvisioningStateTypeCreating,
		"deleting":  ProvisioningStateTypeDeleting,
		"failed":    ProvisioningStateTypeFailed,
		"succeeded": ProvisioningStateTypeSucceeded,
		"updating":  ProvisioningStateTypeUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningStateType(input)
	return &out, nil
}

type UsernameDownloadedType string

const (
	UsernameDownloadedTypeError UsernameDownloadedType = "Error"
	UsernameDownloadedTypeNo    UsernameDownloadedType = "No"
	UsernameDownloadedTypeYes   UsernameDownloadedType = "Yes"
)

func PossibleValuesForUsernameDownloadedType() []string {
	return []string{
		string(UsernameDownloadedTypeError),
		string(UsernameDownloadedTypeNo),
		string(UsernameDownloadedTypeYes),
	}
}

func parseUsernameDownloadedType(input string) (*UsernameDownloadedType, error) {
	vals := map[string]UsernameDownloadedType{
		"error": UsernameDownloadedTypeError,
		"no":    UsernameDownloadedTypeNo,
		"yes":   UsernameDownloadedTypeYes,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsernameDownloadedType(input)
	return &out, nil
}

type UsernameSource string

const (
	UsernameSourceAD   UsernameSource = "AD"
	UsernameSourceFile UsernameSource = "File"
	UsernameSourceLDAP UsernameSource = "LDAP"
	UsernameSourceNone UsernameSource = "None"
)

func PossibleValuesForUsernameSource() []string {
	return []string{
		string(UsernameSourceAD),
		string(UsernameSourceFile),
		string(UsernameSourceLDAP),
		string(UsernameSourceNone),
	}
}

func parseUsernameSource(input string) (*UsernameSource, error) {
	vals := map[string]UsernameSource{
		"ad":   UsernameSourceAD,
		"file": UsernameSourceFile,
		"ldap": UsernameSourceLDAP,
		"none": UsernameSourceNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsernameSource(input)
	return &out, nil
}
