package referencedatasets

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataStringComparisonBehavior string

const (
	DataStringComparisonBehaviorOrdinal           DataStringComparisonBehavior = "Ordinal"
	DataStringComparisonBehaviorOrdinalIgnoreCase DataStringComparisonBehavior = "OrdinalIgnoreCase"
)

func PossibleValuesForDataStringComparisonBehavior() []string {
	return []string{
		string(DataStringComparisonBehaviorOrdinal),
		string(DataStringComparisonBehaviorOrdinalIgnoreCase),
	}
}

func parseDataStringComparisonBehavior(input string) (*DataStringComparisonBehavior, error) {
	vals := map[string]DataStringComparisonBehavior{
		"ordinal":           DataStringComparisonBehaviorOrdinal,
		"ordinalignorecase": DataStringComparisonBehaviorOrdinalIgnoreCase,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataStringComparisonBehavior(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":  ProvisioningStateAccepted,
		"creating":  ProvisioningStateCreating,
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

type ReferenceDataKeyPropertyType string

const (
	ReferenceDataKeyPropertyTypeBool     ReferenceDataKeyPropertyType = "Bool"
	ReferenceDataKeyPropertyTypeDateTime ReferenceDataKeyPropertyType = "DateTime"
	ReferenceDataKeyPropertyTypeDouble   ReferenceDataKeyPropertyType = "Double"
	ReferenceDataKeyPropertyTypeString   ReferenceDataKeyPropertyType = "String"
)

func PossibleValuesForReferenceDataKeyPropertyType() []string {
	return []string{
		string(ReferenceDataKeyPropertyTypeBool),
		string(ReferenceDataKeyPropertyTypeDateTime),
		string(ReferenceDataKeyPropertyTypeDouble),
		string(ReferenceDataKeyPropertyTypeString),
	}
}

func parseReferenceDataKeyPropertyType(input string) (*ReferenceDataKeyPropertyType, error) {
	vals := map[string]ReferenceDataKeyPropertyType{
		"bool":     ReferenceDataKeyPropertyTypeBool,
		"datetime": ReferenceDataKeyPropertyTypeDateTime,
		"double":   ReferenceDataKeyPropertyTypeDouble,
		"string":   ReferenceDataKeyPropertyTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReferenceDataKeyPropertyType(input)
	return &out, nil
}
