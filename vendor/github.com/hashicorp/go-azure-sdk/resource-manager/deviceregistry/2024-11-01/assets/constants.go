package assets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataPointObservabilityMode string

const (
	DataPointObservabilityModeCounter   DataPointObservabilityMode = "Counter"
	DataPointObservabilityModeGauge     DataPointObservabilityMode = "Gauge"
	DataPointObservabilityModeHistogram DataPointObservabilityMode = "Histogram"
	DataPointObservabilityModeLog       DataPointObservabilityMode = "Log"
	DataPointObservabilityModeNone      DataPointObservabilityMode = "None"
)

func PossibleValuesForDataPointObservabilityMode() []string {
	return []string{
		string(DataPointObservabilityModeCounter),
		string(DataPointObservabilityModeGauge),
		string(DataPointObservabilityModeHistogram),
		string(DataPointObservabilityModeLog),
		string(DataPointObservabilityModeNone),
	}
}

func (s *DataPointObservabilityMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataPointObservabilityMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataPointObservabilityMode(input string) (*DataPointObservabilityMode, error) {
	vals := map[string]DataPointObservabilityMode{
		"counter":   DataPointObservabilityModeCounter,
		"gauge":     DataPointObservabilityModeGauge,
		"histogram": DataPointObservabilityModeHistogram,
		"log":       DataPointObservabilityModeLog,
		"none":      DataPointObservabilityModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataPointObservabilityMode(input)
	return &out, nil
}

type EventObservabilityMode string

const (
	EventObservabilityModeLog  EventObservabilityMode = "Log"
	EventObservabilityModeNone EventObservabilityMode = "None"
)

func PossibleValuesForEventObservabilityMode() []string {
	return []string{
		string(EventObservabilityModeLog),
		string(EventObservabilityModeNone),
	}
}

func (s *EventObservabilityMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEventObservabilityMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEventObservabilityMode(input string) (*EventObservabilityMode, error) {
	vals := map[string]EventObservabilityMode{
		"log":  EventObservabilityModeLog,
		"none": EventObservabilityModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventObservabilityMode(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
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
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type TopicRetainType string

const (
	TopicRetainTypeKeep  TopicRetainType = "Keep"
	TopicRetainTypeNever TopicRetainType = "Never"
)

func PossibleValuesForTopicRetainType() []string {
	return []string{
		string(TopicRetainTypeKeep),
		string(TopicRetainTypeNever),
	}
}

func (s *TopicRetainType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTopicRetainType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTopicRetainType(input string) (*TopicRetainType, error) {
	vals := map[string]TopicRetainType{
		"keep":  TopicRetainTypeKeep,
		"never": TopicRetainTypeNever,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TopicRetainType(input)
	return &out, nil
}
