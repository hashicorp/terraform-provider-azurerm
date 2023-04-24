package liveevents

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AsyncOperationStatus string

const (
	AsyncOperationStatusFailed     AsyncOperationStatus = "Failed"
	AsyncOperationStatusInProgress AsyncOperationStatus = "InProgress"
	AsyncOperationStatusSucceeded  AsyncOperationStatus = "Succeeded"
)

func PossibleValuesForAsyncOperationStatus() []string {
	return []string{
		string(AsyncOperationStatusFailed),
		string(AsyncOperationStatusInProgress),
		string(AsyncOperationStatusSucceeded),
	}
}

func (s *AsyncOperationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAsyncOperationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAsyncOperationStatus(input string) (*AsyncOperationStatus, error) {
	vals := map[string]AsyncOperationStatus{
		"failed":     AsyncOperationStatusFailed,
		"inprogress": AsyncOperationStatusInProgress,
		"succeeded":  AsyncOperationStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AsyncOperationStatus(input)
	return &out, nil
}

type LiveEventEncodingType string

const (
	LiveEventEncodingTypeNone                     LiveEventEncodingType = "None"
	LiveEventEncodingTypePassthroughBasic         LiveEventEncodingType = "PassthroughBasic"
	LiveEventEncodingTypePassthroughStandard      LiveEventEncodingType = "PassthroughStandard"
	LiveEventEncodingTypePremiumOneZeroEightZerop LiveEventEncodingType = "Premium1080p"
	LiveEventEncodingTypeStandard                 LiveEventEncodingType = "Standard"
)

func PossibleValuesForLiveEventEncodingType() []string {
	return []string{
		string(LiveEventEncodingTypeNone),
		string(LiveEventEncodingTypePassthroughBasic),
		string(LiveEventEncodingTypePassthroughStandard),
		string(LiveEventEncodingTypePremiumOneZeroEightZerop),
		string(LiveEventEncodingTypeStandard),
	}
}

func (s *LiveEventEncodingType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLiveEventEncodingType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLiveEventEncodingType(input string) (*LiveEventEncodingType, error) {
	vals := map[string]LiveEventEncodingType{
		"none":                LiveEventEncodingTypeNone,
		"passthroughbasic":    LiveEventEncodingTypePassthroughBasic,
		"passthroughstandard": LiveEventEncodingTypePassthroughStandard,
		"premium1080p":        LiveEventEncodingTypePremiumOneZeroEightZerop,
		"standard":            LiveEventEncodingTypeStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LiveEventEncodingType(input)
	return &out, nil
}

type LiveEventInputProtocol string

const (
	LiveEventInputProtocolFragmentedMPFour LiveEventInputProtocol = "FragmentedMP4"
	LiveEventInputProtocolRTMP             LiveEventInputProtocol = "RTMP"
)

func PossibleValuesForLiveEventInputProtocol() []string {
	return []string{
		string(LiveEventInputProtocolFragmentedMPFour),
		string(LiveEventInputProtocolRTMP),
	}
}

func (s *LiveEventInputProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLiveEventInputProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLiveEventInputProtocol(input string) (*LiveEventInputProtocol, error) {
	vals := map[string]LiveEventInputProtocol{
		"fragmentedmp4": LiveEventInputProtocolFragmentedMPFour,
		"rtmp":          LiveEventInputProtocolRTMP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LiveEventInputProtocol(input)
	return &out, nil
}

type LiveEventResourceState string

const (
	LiveEventResourceStateAllocating LiveEventResourceState = "Allocating"
	LiveEventResourceStateDeleting   LiveEventResourceState = "Deleting"
	LiveEventResourceStateRunning    LiveEventResourceState = "Running"
	LiveEventResourceStateStandBy    LiveEventResourceState = "StandBy"
	LiveEventResourceStateStarting   LiveEventResourceState = "Starting"
	LiveEventResourceStateStopped    LiveEventResourceState = "Stopped"
	LiveEventResourceStateStopping   LiveEventResourceState = "Stopping"
)

func PossibleValuesForLiveEventResourceState() []string {
	return []string{
		string(LiveEventResourceStateAllocating),
		string(LiveEventResourceStateDeleting),
		string(LiveEventResourceStateRunning),
		string(LiveEventResourceStateStandBy),
		string(LiveEventResourceStateStarting),
		string(LiveEventResourceStateStopped),
		string(LiveEventResourceStateStopping),
	}
}

func (s *LiveEventResourceState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLiveEventResourceState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLiveEventResourceState(input string) (*LiveEventResourceState, error) {
	vals := map[string]LiveEventResourceState{
		"allocating": LiveEventResourceStateAllocating,
		"deleting":   LiveEventResourceStateDeleting,
		"running":    LiveEventResourceStateRunning,
		"standby":    LiveEventResourceStateStandBy,
		"starting":   LiveEventResourceStateStarting,
		"stopped":    LiveEventResourceStateStopped,
		"stopping":   LiveEventResourceStateStopping,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LiveEventResourceState(input)
	return &out, nil
}

type StreamOptionsFlag string

const (
	StreamOptionsFlagDefault        StreamOptionsFlag = "Default"
	StreamOptionsFlagLowLatency     StreamOptionsFlag = "LowLatency"
	StreamOptionsFlagLowLatencyVTwo StreamOptionsFlag = "LowLatencyV2"
)

func PossibleValuesForStreamOptionsFlag() []string {
	return []string{
		string(StreamOptionsFlagDefault),
		string(StreamOptionsFlagLowLatency),
		string(StreamOptionsFlagLowLatencyVTwo),
	}
}

func (s *StreamOptionsFlag) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStreamOptionsFlag(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStreamOptionsFlag(input string) (*StreamOptionsFlag, error) {
	vals := map[string]StreamOptionsFlag{
		"default":      StreamOptionsFlagDefault,
		"lowlatency":   StreamOptionsFlagLowLatency,
		"lowlatencyv2": StreamOptionsFlagLowLatencyVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StreamOptionsFlag(input)
	return &out, nil
}

type StretchMode string

const (
	StretchModeAutoFit  StretchMode = "AutoFit"
	StretchModeAutoSize StretchMode = "AutoSize"
	StretchModeNone     StretchMode = "None"
)

func PossibleValuesForStretchMode() []string {
	return []string{
		string(StretchModeAutoFit),
		string(StretchModeAutoSize),
		string(StretchModeNone),
	}
}

func (s *StretchMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStretchMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStretchMode(input string) (*StretchMode, error) {
	vals := map[string]StretchMode{
		"autofit":  StretchModeAutoFit,
		"autosize": StretchModeAutoSize,
		"none":     StretchModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StretchMode(input)
	return &out, nil
}
