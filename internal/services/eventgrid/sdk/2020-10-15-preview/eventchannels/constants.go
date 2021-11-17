package eventchannels

import "strings"

type AdvancedFilterOperatorType string

const (
	AdvancedFilterOperatorTypeBoolEquals                AdvancedFilterOperatorType = "BoolEquals"
	AdvancedFilterOperatorTypeIsNotNull                 AdvancedFilterOperatorType = "IsNotNull"
	AdvancedFilterOperatorTypeIsNullOrUndefined         AdvancedFilterOperatorType = "IsNullOrUndefined"
	AdvancedFilterOperatorTypeNumberGreaterThan         AdvancedFilterOperatorType = "NumberGreaterThan"
	AdvancedFilterOperatorTypeNumberGreaterThanOrEquals AdvancedFilterOperatorType = "NumberGreaterThanOrEquals"
	AdvancedFilterOperatorTypeNumberIn                  AdvancedFilterOperatorType = "NumberIn"
	AdvancedFilterOperatorTypeNumberInRange             AdvancedFilterOperatorType = "NumberInRange"
	AdvancedFilterOperatorTypeNumberLessThan            AdvancedFilterOperatorType = "NumberLessThan"
	AdvancedFilterOperatorTypeNumberLessThanOrEquals    AdvancedFilterOperatorType = "NumberLessThanOrEquals"
	AdvancedFilterOperatorTypeNumberNotIn               AdvancedFilterOperatorType = "NumberNotIn"
	AdvancedFilterOperatorTypeNumberNotInRange          AdvancedFilterOperatorType = "NumberNotInRange"
	AdvancedFilterOperatorTypeStringBeginsWith          AdvancedFilterOperatorType = "StringBeginsWith"
	AdvancedFilterOperatorTypeStringContains            AdvancedFilterOperatorType = "StringContains"
	AdvancedFilterOperatorTypeStringEndsWith            AdvancedFilterOperatorType = "StringEndsWith"
	AdvancedFilterOperatorTypeStringIn                  AdvancedFilterOperatorType = "StringIn"
	AdvancedFilterOperatorTypeStringNotBeginsWith       AdvancedFilterOperatorType = "StringNotBeginsWith"
	AdvancedFilterOperatorTypeStringNotContains         AdvancedFilterOperatorType = "StringNotContains"
	AdvancedFilterOperatorTypeStringNotEndsWith         AdvancedFilterOperatorType = "StringNotEndsWith"
	AdvancedFilterOperatorTypeStringNotIn               AdvancedFilterOperatorType = "StringNotIn"
)

func PossibleValuesForAdvancedFilterOperatorType() []string {
	return []string{
		string(AdvancedFilterOperatorTypeBoolEquals),
		string(AdvancedFilterOperatorTypeIsNotNull),
		string(AdvancedFilterOperatorTypeIsNullOrUndefined),
		string(AdvancedFilterOperatorTypeNumberGreaterThan),
		string(AdvancedFilterOperatorTypeNumberGreaterThanOrEquals),
		string(AdvancedFilterOperatorTypeNumberIn),
		string(AdvancedFilterOperatorTypeNumberInRange),
		string(AdvancedFilterOperatorTypeNumberLessThan),
		string(AdvancedFilterOperatorTypeNumberLessThanOrEquals),
		string(AdvancedFilterOperatorTypeNumberNotIn),
		string(AdvancedFilterOperatorTypeNumberNotInRange),
		string(AdvancedFilterOperatorTypeStringBeginsWith),
		string(AdvancedFilterOperatorTypeStringContains),
		string(AdvancedFilterOperatorTypeStringEndsWith),
		string(AdvancedFilterOperatorTypeStringIn),
		string(AdvancedFilterOperatorTypeStringNotBeginsWith),
		string(AdvancedFilterOperatorTypeStringNotContains),
		string(AdvancedFilterOperatorTypeStringNotEndsWith),
		string(AdvancedFilterOperatorTypeStringNotIn),
	}
}

func parseAdvancedFilterOperatorType(input string) (*AdvancedFilterOperatorType, error) {
	vals := map[string]AdvancedFilterOperatorType{
		"boolequals":                AdvancedFilterOperatorTypeBoolEquals,
		"isnotnull":                 AdvancedFilterOperatorTypeIsNotNull,
		"isnullorundefined":         AdvancedFilterOperatorTypeIsNullOrUndefined,
		"numbergreaterthan":         AdvancedFilterOperatorTypeNumberGreaterThan,
		"numbergreaterthanorequals": AdvancedFilterOperatorTypeNumberGreaterThanOrEquals,
		"numberin":                  AdvancedFilterOperatorTypeNumberIn,
		"numberinrange":             AdvancedFilterOperatorTypeNumberInRange,
		"numberlessthan":            AdvancedFilterOperatorTypeNumberLessThan,
		"numberlessthanorequals":    AdvancedFilterOperatorTypeNumberLessThanOrEquals,
		"numbernotin":               AdvancedFilterOperatorTypeNumberNotIn,
		"numbernotinrange":          AdvancedFilterOperatorTypeNumberNotInRange,
		"stringbeginswith":          AdvancedFilterOperatorTypeStringBeginsWith,
		"stringcontains":            AdvancedFilterOperatorTypeStringContains,
		"stringendswith":            AdvancedFilterOperatorTypeStringEndsWith,
		"stringin":                  AdvancedFilterOperatorTypeStringIn,
		"stringnotbeginswith":       AdvancedFilterOperatorTypeStringNotBeginsWith,
		"stringnotcontains":         AdvancedFilterOperatorTypeStringNotContains,
		"stringnotendswith":         AdvancedFilterOperatorTypeStringNotEndsWith,
		"stringnotin":               AdvancedFilterOperatorTypeStringNotIn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AdvancedFilterOperatorType(input)
	return &out, nil
}

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

type EventChannelProvisioningState string

const (
	EventChannelProvisioningStateCanceled  EventChannelProvisioningState = "Canceled"
	EventChannelProvisioningStateCreating  EventChannelProvisioningState = "Creating"
	EventChannelProvisioningStateDeleting  EventChannelProvisioningState = "Deleting"
	EventChannelProvisioningStateFailed    EventChannelProvisioningState = "Failed"
	EventChannelProvisioningStateSucceeded EventChannelProvisioningState = "Succeeded"
	EventChannelProvisioningStateUpdating  EventChannelProvisioningState = "Updating"
)

func PossibleValuesForEventChannelProvisioningState() []string {
	return []string{
		string(EventChannelProvisioningStateCanceled),
		string(EventChannelProvisioningStateCreating),
		string(EventChannelProvisioningStateDeleting),
		string(EventChannelProvisioningStateFailed),
		string(EventChannelProvisioningStateSucceeded),
		string(EventChannelProvisioningStateUpdating),
	}
}

func parseEventChannelProvisioningState(input string) (*EventChannelProvisioningState, error) {
	vals := map[string]EventChannelProvisioningState{
		"canceled":  EventChannelProvisioningStateCanceled,
		"creating":  EventChannelProvisioningStateCreating,
		"deleting":  EventChannelProvisioningStateDeleting,
		"failed":    EventChannelProvisioningStateFailed,
		"succeeded": EventChannelProvisioningStateSucceeded,
		"updating":  EventChannelProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventChannelProvisioningState(input)
	return &out, nil
}

type PartnerTopicReadinessState string

const (
	PartnerTopicReadinessStateActivatedByUser       PartnerTopicReadinessState = "ActivatedByUser"
	PartnerTopicReadinessStateDeactivatedByUser     PartnerTopicReadinessState = "DeactivatedByUser"
	PartnerTopicReadinessStateDeletedByUser         PartnerTopicReadinessState = "DeletedByUser"
	PartnerTopicReadinessStateNotActivatedByUserYet PartnerTopicReadinessState = "NotActivatedByUserYet"
)

func PossibleValuesForPartnerTopicReadinessState() []string {
	return []string{
		string(PartnerTopicReadinessStateActivatedByUser),
		string(PartnerTopicReadinessStateDeactivatedByUser),
		string(PartnerTopicReadinessStateDeletedByUser),
		string(PartnerTopicReadinessStateNotActivatedByUserYet),
	}
}

func parsePartnerTopicReadinessState(input string) (*PartnerTopicReadinessState, error) {
	vals := map[string]PartnerTopicReadinessState{
		"activatedbyuser":       PartnerTopicReadinessStateActivatedByUser,
		"deactivatedbyuser":     PartnerTopicReadinessStateDeactivatedByUser,
		"deletedbyuser":         PartnerTopicReadinessStateDeletedByUser,
		"notactivatedbyuseryet": PartnerTopicReadinessStateNotActivatedByUserYet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerTopicReadinessState(input)
	return &out, nil
}
