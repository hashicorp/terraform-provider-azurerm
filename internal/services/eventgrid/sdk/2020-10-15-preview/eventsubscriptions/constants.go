package eventsubscriptions

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

type DeadLetterEndPointType string

const (
	DeadLetterEndPointTypeStorageBlob DeadLetterEndPointType = "StorageBlob"
)

func PossibleValuesForDeadLetterEndPointType() []string {
	return []string{
		string(DeadLetterEndPointTypeStorageBlob),
	}
}

func parseDeadLetterEndPointType(input string) (*DeadLetterEndPointType, error) {
	vals := map[string]DeadLetterEndPointType{
		"storageblob": DeadLetterEndPointTypeStorageBlob,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeadLetterEndPointType(input)
	return &out, nil
}

type DeliveryAttributeMappingType string

const (
	DeliveryAttributeMappingTypeDynamic DeliveryAttributeMappingType = "Dynamic"
	DeliveryAttributeMappingTypeStatic  DeliveryAttributeMappingType = "Static"
)

func PossibleValuesForDeliveryAttributeMappingType() []string {
	return []string{
		string(DeliveryAttributeMappingTypeDynamic),
		string(DeliveryAttributeMappingTypeStatic),
	}
}

func parseDeliveryAttributeMappingType(input string) (*DeliveryAttributeMappingType, error) {
	vals := map[string]DeliveryAttributeMappingType{
		"dynamic": DeliveryAttributeMappingTypeDynamic,
		"static":  DeliveryAttributeMappingTypeStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeliveryAttributeMappingType(input)
	return &out, nil
}

type EndpointType string

const (
	EndpointTypeAzureFunction    EndpointType = "AzureFunction"
	EndpointTypeEventHub         EndpointType = "EventHub"
	EndpointTypeHybridConnection EndpointType = "HybridConnection"
	EndpointTypeServiceBusQueue  EndpointType = "ServiceBusQueue"
	EndpointTypeServiceBusTopic  EndpointType = "ServiceBusTopic"
	EndpointTypeStorageQueue     EndpointType = "StorageQueue"
	EndpointTypeWebHook          EndpointType = "WebHook"
)

func PossibleValuesForEndpointType() []string {
	return []string{
		string(EndpointTypeAzureFunction),
		string(EndpointTypeEventHub),
		string(EndpointTypeHybridConnection),
		string(EndpointTypeServiceBusQueue),
		string(EndpointTypeServiceBusTopic),
		string(EndpointTypeStorageQueue),
		string(EndpointTypeWebHook),
	}
}

func parseEndpointType(input string) (*EndpointType, error) {
	vals := map[string]EndpointType{
		"azurefunction":    EndpointTypeAzureFunction,
		"eventhub":         EndpointTypeEventHub,
		"hybridconnection": EndpointTypeHybridConnection,
		"servicebusqueue":  EndpointTypeServiceBusQueue,
		"servicebustopic":  EndpointTypeServiceBusTopic,
		"storagequeue":     EndpointTypeStorageQueue,
		"webhook":          EndpointTypeWebHook,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointType(input)
	return &out, nil
}

type EventDeliverySchema string

const (
	EventDeliverySchemaCloudEventSchemaVOneZero EventDeliverySchema = "CloudEventSchemaV1_0"
	EventDeliverySchemaCustomInputSchema        EventDeliverySchema = "CustomInputSchema"
	EventDeliverySchemaEventGridSchema          EventDeliverySchema = "EventGridSchema"
)

func PossibleValuesForEventDeliverySchema() []string {
	return []string{
		string(EventDeliverySchemaCloudEventSchemaVOneZero),
		string(EventDeliverySchemaCustomInputSchema),
		string(EventDeliverySchemaEventGridSchema),
	}
}

func parseEventDeliverySchema(input string) (*EventDeliverySchema, error) {
	vals := map[string]EventDeliverySchema{
		"cloudeventschemav1_0": EventDeliverySchemaCloudEventSchemaVOneZero,
		"custominputschema":    EventDeliverySchemaCustomInputSchema,
		"eventgridschema":      EventDeliverySchemaEventGridSchema,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventDeliverySchema(input)
	return &out, nil
}

type EventSubscriptionIdentityType string

const (
	EventSubscriptionIdentityTypeSystemAssigned             EventSubscriptionIdentityType = "SystemAssigned"
	EventSubscriptionIdentityTypeSystemAssignedUserAssigned EventSubscriptionIdentityType = "SystemAssigned, UserAssigned"
	EventSubscriptionIdentityTypeUserAssigned               EventSubscriptionIdentityType = "UserAssigned"
)

func PossibleValuesForEventSubscriptionIdentityType() []string {
	return []string{
		string(EventSubscriptionIdentityTypeSystemAssigned),
		string(EventSubscriptionIdentityTypeSystemAssignedUserAssigned),
		string(EventSubscriptionIdentityTypeUserAssigned),
	}
}

func parseEventSubscriptionIdentityType(input string) (*EventSubscriptionIdentityType, error) {
	vals := map[string]EventSubscriptionIdentityType{
		"systemassigned":               EventSubscriptionIdentityTypeSystemAssigned,
		"systemassigned, userassigned": EventSubscriptionIdentityTypeSystemAssignedUserAssigned,
		"userassigned":                 EventSubscriptionIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventSubscriptionIdentityType(input)
	return &out, nil
}

type EventSubscriptionProvisioningState string

const (
	EventSubscriptionProvisioningStateAwaitingManualAction EventSubscriptionProvisioningState = "AwaitingManualAction"
	EventSubscriptionProvisioningStateCanceled             EventSubscriptionProvisioningState = "Canceled"
	EventSubscriptionProvisioningStateCreating             EventSubscriptionProvisioningState = "Creating"
	EventSubscriptionProvisioningStateDeleting             EventSubscriptionProvisioningState = "Deleting"
	EventSubscriptionProvisioningStateFailed               EventSubscriptionProvisioningState = "Failed"
	EventSubscriptionProvisioningStateSucceeded            EventSubscriptionProvisioningState = "Succeeded"
	EventSubscriptionProvisioningStateUpdating             EventSubscriptionProvisioningState = "Updating"
)

func PossibleValuesForEventSubscriptionProvisioningState() []string {
	return []string{
		string(EventSubscriptionProvisioningStateAwaitingManualAction),
		string(EventSubscriptionProvisioningStateCanceled),
		string(EventSubscriptionProvisioningStateCreating),
		string(EventSubscriptionProvisioningStateDeleting),
		string(EventSubscriptionProvisioningStateFailed),
		string(EventSubscriptionProvisioningStateSucceeded),
		string(EventSubscriptionProvisioningStateUpdating),
	}
}

func parseEventSubscriptionProvisioningState(input string) (*EventSubscriptionProvisioningState, error) {
	vals := map[string]EventSubscriptionProvisioningState{
		"awaitingmanualaction": EventSubscriptionProvisioningStateAwaitingManualAction,
		"canceled":             EventSubscriptionProvisioningStateCanceled,
		"creating":             EventSubscriptionProvisioningStateCreating,
		"deleting":             EventSubscriptionProvisioningStateDeleting,
		"failed":               EventSubscriptionProvisioningStateFailed,
		"succeeded":            EventSubscriptionProvisioningStateSucceeded,
		"updating":             EventSubscriptionProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventSubscriptionProvisioningState(input)
	return &out, nil
}
