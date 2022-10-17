package collectorpolicies

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

func PossibleValuesForCreatedByType() []string {
	return []string{
		string(CreatedByTypeApplication),
		string(CreatedByTypeKey),
		string(CreatedByTypeManagedIdentity),
		string(CreatedByTypeUser),
	}
}

func parseCreatedByType(input string) (*CreatedByType, error) {
	vals := map[string]CreatedByType{
		"application":     CreatedByTypeApplication,
		"key":             CreatedByTypeKey,
		"managedidentity": CreatedByTypeManagedIdentity,
		"user":            CreatedByTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreatedByType(input)
	return &out, nil
}

type DestinationType string

const (
	DestinationTypeAzureMonitor DestinationType = "AzureMonitor"
)

func PossibleValuesForDestinationType() []string {
	return []string{
		string(DestinationTypeAzureMonitor),
	}
}

func parseDestinationType(input string) (*DestinationType, error) {
	vals := map[string]DestinationType{
		"azuremonitor": DestinationTypeAzureMonitor,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DestinationType(input)
	return &out, nil
}

type EmissionType string

const (
	EmissionTypeIPFIX EmissionType = "IPFIX"
)

func PossibleValuesForEmissionType() []string {
	return []string{
		string(EmissionTypeIPFIX),
	}
}

func parseEmissionType(input string) (*EmissionType, error) {
	vals := map[string]EmissionType{
		"ipfix": EmissionTypeIPFIX,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EmissionType(input)
	return &out, nil
}

type IngestionType string

const (
	IngestionTypeIPFIX IngestionType = "IPFIX"
)

func PossibleValuesForIngestionType() []string {
	return []string{
		string(IngestionTypeIPFIX),
	}
}

func parseIngestionType(input string) (*IngestionType, error) {
	vals := map[string]IngestionType{
		"ipfix": IngestionTypeIPFIX,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IngestionType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
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

type SourceType string

const (
	SourceTypeResource SourceType = "Resource"
)

func PossibleValuesForSourceType() []string {
	return []string{
		string(SourceTypeResource),
	}
}

func parseSourceType(input string) (*SourceType, error) {
	vals := map[string]SourceType{
		"resource": SourceTypeResource,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SourceType(input)
	return &out, nil
}
