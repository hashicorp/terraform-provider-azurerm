package eventsubscriptions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

func (s *AdvancedFilterOperatorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAdvancedFilterOperatorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type DeadLetterEndPointType string

const (
	DeadLetterEndPointTypeStorageBlob DeadLetterEndPointType = "StorageBlob"
)

func PossibleValuesForDeadLetterEndPointType() []string {
	return []string{
		string(DeadLetterEndPointTypeStorageBlob),
	}
}

func (s *DeadLetterEndPointType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeadLetterEndPointType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *DeliveryAttributeMappingType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeliveryAttributeMappingType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type DeliveryMode string

const (
	DeliveryModePush  DeliveryMode = "Push"
	DeliveryModeQueue DeliveryMode = "Queue"
)

func PossibleValuesForDeliveryMode() []string {
	return []string{
		string(DeliveryModePush),
		string(DeliveryModeQueue),
	}
}

func (s *DeliveryMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeliveryMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeliveryMode(input string) (*DeliveryMode, error) {
	vals := map[string]DeliveryMode{
		"push":  DeliveryModePush,
		"queue": DeliveryModeQueue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeliveryMode(input)
	return &out, nil
}

type DeliverySchema string

const (
	DeliverySchemaCloudEventSchemaVOneZero DeliverySchema = "CloudEventSchemaV1_0"
)

func PossibleValuesForDeliverySchema() []string {
	return []string{
		string(DeliverySchemaCloudEventSchemaVOneZero),
	}
}

func (s *DeliverySchema) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeliverySchema(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeliverySchema(input string) (*DeliverySchema, error) {
	vals := map[string]DeliverySchema{
		"cloudeventschemav1_0": DeliverySchemaCloudEventSchemaVOneZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeliverySchema(input)
	return &out, nil
}

type EndpointType string

const (
	EndpointTypeAzureFunction    EndpointType = "AzureFunction"
	EndpointTypeEventHub         EndpointType = "EventHub"
	EndpointTypeHybridConnection EndpointType = "HybridConnection"
	EndpointTypeMonitorAlert     EndpointType = "MonitorAlert"
	EndpointTypeNamespaceTopic   EndpointType = "NamespaceTopic"
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
		string(EndpointTypeMonitorAlert),
		string(EndpointTypeNamespaceTopic),
		string(EndpointTypeServiceBusQueue),
		string(EndpointTypeServiceBusTopic),
		string(EndpointTypeStorageQueue),
		string(EndpointTypeWebHook),
	}
}

func (s *EndpointType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEndpointType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEndpointType(input string) (*EndpointType, error) {
	vals := map[string]EndpointType{
		"azurefunction":    EndpointTypeAzureFunction,
		"eventhub":         EndpointTypeEventHub,
		"hybridconnection": EndpointTypeHybridConnection,
		"monitoralert":     EndpointTypeMonitorAlert,
		"namespacetopic":   EndpointTypeNamespaceTopic,
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

func (s *EventDeliverySchema) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEventDeliverySchema(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
	EventSubscriptionIdentityTypeSystemAssigned EventSubscriptionIdentityType = "SystemAssigned"
	EventSubscriptionIdentityTypeUserAssigned   EventSubscriptionIdentityType = "UserAssigned"
)

func PossibleValuesForEventSubscriptionIdentityType() []string {
	return []string{
		string(EventSubscriptionIdentityTypeSystemAssigned),
		string(EventSubscriptionIdentityTypeUserAssigned),
	}
}

func (s *EventSubscriptionIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEventSubscriptionIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEventSubscriptionIdentityType(input string) (*EventSubscriptionIdentityType, error) {
	vals := map[string]EventSubscriptionIdentityType{
		"systemassigned": EventSubscriptionIdentityTypeSystemAssigned,
		"userassigned":   EventSubscriptionIdentityTypeUserAssigned,
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

func (s *EventSubscriptionProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEventSubscriptionProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type FilterOperatorType string

const (
	FilterOperatorTypeBoolEquals                FilterOperatorType = "BoolEquals"
	FilterOperatorTypeIsNotNull                 FilterOperatorType = "IsNotNull"
	FilterOperatorTypeIsNullOrUndefined         FilterOperatorType = "IsNullOrUndefined"
	FilterOperatorTypeNumberGreaterThan         FilterOperatorType = "NumberGreaterThan"
	FilterOperatorTypeNumberGreaterThanOrEquals FilterOperatorType = "NumberGreaterThanOrEquals"
	FilterOperatorTypeNumberIn                  FilterOperatorType = "NumberIn"
	FilterOperatorTypeNumberInRange             FilterOperatorType = "NumberInRange"
	FilterOperatorTypeNumberLessThan            FilterOperatorType = "NumberLessThan"
	FilterOperatorTypeNumberLessThanOrEquals    FilterOperatorType = "NumberLessThanOrEquals"
	FilterOperatorTypeNumberNotIn               FilterOperatorType = "NumberNotIn"
	FilterOperatorTypeNumberNotInRange          FilterOperatorType = "NumberNotInRange"
	FilterOperatorTypeStringBeginsWith          FilterOperatorType = "StringBeginsWith"
	FilterOperatorTypeStringContains            FilterOperatorType = "StringContains"
	FilterOperatorTypeStringEndsWith            FilterOperatorType = "StringEndsWith"
	FilterOperatorTypeStringIn                  FilterOperatorType = "StringIn"
	FilterOperatorTypeStringNotBeginsWith       FilterOperatorType = "StringNotBeginsWith"
	FilterOperatorTypeStringNotContains         FilterOperatorType = "StringNotContains"
	FilterOperatorTypeStringNotEndsWith         FilterOperatorType = "StringNotEndsWith"
	FilterOperatorTypeStringNotIn               FilterOperatorType = "StringNotIn"
)

func PossibleValuesForFilterOperatorType() []string {
	return []string{
		string(FilterOperatorTypeBoolEquals),
		string(FilterOperatorTypeIsNotNull),
		string(FilterOperatorTypeIsNullOrUndefined),
		string(FilterOperatorTypeNumberGreaterThan),
		string(FilterOperatorTypeNumberGreaterThanOrEquals),
		string(FilterOperatorTypeNumberIn),
		string(FilterOperatorTypeNumberInRange),
		string(FilterOperatorTypeNumberLessThan),
		string(FilterOperatorTypeNumberLessThanOrEquals),
		string(FilterOperatorTypeNumberNotIn),
		string(FilterOperatorTypeNumberNotInRange),
		string(FilterOperatorTypeStringBeginsWith),
		string(FilterOperatorTypeStringContains),
		string(FilterOperatorTypeStringEndsWith),
		string(FilterOperatorTypeStringIn),
		string(FilterOperatorTypeStringNotBeginsWith),
		string(FilterOperatorTypeStringNotContains),
		string(FilterOperatorTypeStringNotEndsWith),
		string(FilterOperatorTypeStringNotIn),
	}
}

func (s *FilterOperatorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFilterOperatorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFilterOperatorType(input string) (*FilterOperatorType, error) {
	vals := map[string]FilterOperatorType{
		"boolequals":                FilterOperatorTypeBoolEquals,
		"isnotnull":                 FilterOperatorTypeIsNotNull,
		"isnullorundefined":         FilterOperatorTypeIsNullOrUndefined,
		"numbergreaterthan":         FilterOperatorTypeNumberGreaterThan,
		"numbergreaterthanorequals": FilterOperatorTypeNumberGreaterThanOrEquals,
		"numberin":                  FilterOperatorTypeNumberIn,
		"numberinrange":             FilterOperatorTypeNumberInRange,
		"numberlessthan":            FilterOperatorTypeNumberLessThan,
		"numberlessthanorequals":    FilterOperatorTypeNumberLessThanOrEquals,
		"numbernotin":               FilterOperatorTypeNumberNotIn,
		"numbernotinrange":          FilterOperatorTypeNumberNotInRange,
		"stringbeginswith":          FilterOperatorTypeStringBeginsWith,
		"stringcontains":            FilterOperatorTypeStringContains,
		"stringendswith":            FilterOperatorTypeStringEndsWith,
		"stringin":                  FilterOperatorTypeStringIn,
		"stringnotbeginswith":       FilterOperatorTypeStringNotBeginsWith,
		"stringnotcontains":         FilterOperatorTypeStringNotContains,
		"stringnotendswith":         FilterOperatorTypeStringNotEndsWith,
		"stringnotin":               FilterOperatorTypeStringNotIn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FilterOperatorType(input)
	return &out, nil
}

type MonitorAlertSeverity string

const (
	MonitorAlertSeveritySevFour  MonitorAlertSeverity = "Sev4"
	MonitorAlertSeveritySevOne   MonitorAlertSeverity = "Sev1"
	MonitorAlertSeveritySevThree MonitorAlertSeverity = "Sev3"
	MonitorAlertSeveritySevTwo   MonitorAlertSeverity = "Sev2"
	MonitorAlertSeveritySevZero  MonitorAlertSeverity = "Sev0"
)

func PossibleValuesForMonitorAlertSeverity() []string {
	return []string{
		string(MonitorAlertSeveritySevFour),
		string(MonitorAlertSeveritySevOne),
		string(MonitorAlertSeveritySevThree),
		string(MonitorAlertSeveritySevTwo),
		string(MonitorAlertSeveritySevZero),
	}
}

func (s *MonitorAlertSeverity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMonitorAlertSeverity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMonitorAlertSeverity(input string) (*MonitorAlertSeverity, error) {
	vals := map[string]MonitorAlertSeverity{
		"sev4": MonitorAlertSeveritySevFour,
		"sev1": MonitorAlertSeveritySevOne,
		"sev3": MonitorAlertSeveritySevThree,
		"sev2": MonitorAlertSeveritySevTwo,
		"sev0": MonitorAlertSeveritySevZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MonitorAlertSeverity(input)
	return &out, nil
}

type SubscriptionProvisioningState string

const (
	SubscriptionProvisioningStateAwaitingManualAction SubscriptionProvisioningState = "AwaitingManualAction"
	SubscriptionProvisioningStateCanceled             SubscriptionProvisioningState = "Canceled"
	SubscriptionProvisioningStateCreateFailed         SubscriptionProvisioningState = "CreateFailed"
	SubscriptionProvisioningStateCreating             SubscriptionProvisioningState = "Creating"
	SubscriptionProvisioningStateDeleteFailed         SubscriptionProvisioningState = "DeleteFailed"
	SubscriptionProvisioningStateDeleted              SubscriptionProvisioningState = "Deleted"
	SubscriptionProvisioningStateDeleting             SubscriptionProvisioningState = "Deleting"
	SubscriptionProvisioningStateFailed               SubscriptionProvisioningState = "Failed"
	SubscriptionProvisioningStateSucceeded            SubscriptionProvisioningState = "Succeeded"
	SubscriptionProvisioningStateUpdatedFailed        SubscriptionProvisioningState = "UpdatedFailed"
	SubscriptionProvisioningStateUpdating             SubscriptionProvisioningState = "Updating"
)

func PossibleValuesForSubscriptionProvisioningState() []string {
	return []string{
		string(SubscriptionProvisioningStateAwaitingManualAction),
		string(SubscriptionProvisioningStateCanceled),
		string(SubscriptionProvisioningStateCreateFailed),
		string(SubscriptionProvisioningStateCreating),
		string(SubscriptionProvisioningStateDeleteFailed),
		string(SubscriptionProvisioningStateDeleted),
		string(SubscriptionProvisioningStateDeleting),
		string(SubscriptionProvisioningStateFailed),
		string(SubscriptionProvisioningStateSucceeded),
		string(SubscriptionProvisioningStateUpdatedFailed),
		string(SubscriptionProvisioningStateUpdating),
	}
}

func (s *SubscriptionProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSubscriptionProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSubscriptionProvisioningState(input string) (*SubscriptionProvisioningState, error) {
	vals := map[string]SubscriptionProvisioningState{
		"awaitingmanualaction": SubscriptionProvisioningStateAwaitingManualAction,
		"canceled":             SubscriptionProvisioningStateCanceled,
		"createfailed":         SubscriptionProvisioningStateCreateFailed,
		"creating":             SubscriptionProvisioningStateCreating,
		"deletefailed":         SubscriptionProvisioningStateDeleteFailed,
		"deleted":              SubscriptionProvisioningStateDeleted,
		"deleting":             SubscriptionProvisioningStateDeleting,
		"failed":               SubscriptionProvisioningStateFailed,
		"succeeded":            SubscriptionProvisioningStateSucceeded,
		"updatedfailed":        SubscriptionProvisioningStateUpdatedFailed,
		"updating":             SubscriptionProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SubscriptionProvisioningState(input)
	return &out, nil
}

type TlsVersion string

const (
	TlsVersionOnePointOne  TlsVersion = "1.1"
	TlsVersionOnePointTwo  TlsVersion = "1.2"
	TlsVersionOnePointZero TlsVersion = "1.0"
)

func PossibleValuesForTlsVersion() []string {
	return []string{
		string(TlsVersionOnePointOne),
		string(TlsVersionOnePointTwo),
		string(TlsVersionOnePointZero),
	}
}

func (s *TlsVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTlsVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTlsVersion(input string) (*TlsVersion, error) {
	vals := map[string]TlsVersion{
		"1.1": TlsVersionOnePointOne,
		"1.2": TlsVersionOnePointTwo,
		"1.0": TlsVersionOnePointZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TlsVersion(input)
	return &out, nil
}
