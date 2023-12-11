package tagrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateDeleted      ProvisioningState = "Deleted"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateNotSpecified ProvisioningState = "NotSpecified"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateNotSpecified),
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
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"creating":     ProvisioningStateCreating,
		"deleted":      ProvisioningStateDeleted,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"notspecified": ProvisioningStateNotSpecified,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type SendAadLogsStatus string

const (
	SendAadLogsStatusDisabled SendAadLogsStatus = "Disabled"
	SendAadLogsStatusEnabled  SendAadLogsStatus = "Enabled"
)

func PossibleValuesForSendAadLogsStatus() []string {
	return []string{
		string(SendAadLogsStatusDisabled),
		string(SendAadLogsStatusEnabled),
	}
}

func (s *SendAadLogsStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSendAadLogsStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSendAadLogsStatus(input string) (*SendAadLogsStatus, error) {
	vals := map[string]SendAadLogsStatus{
		"disabled": SendAadLogsStatusDisabled,
		"enabled":  SendAadLogsStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SendAadLogsStatus(input)
	return &out, nil
}

type SendActivityLogsStatus string

const (
	SendActivityLogsStatusDisabled SendActivityLogsStatus = "Disabled"
	SendActivityLogsStatusEnabled  SendActivityLogsStatus = "Enabled"
)

func PossibleValuesForSendActivityLogsStatus() []string {
	return []string{
		string(SendActivityLogsStatusDisabled),
		string(SendActivityLogsStatusEnabled),
	}
}

func (s *SendActivityLogsStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSendActivityLogsStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSendActivityLogsStatus(input string) (*SendActivityLogsStatus, error) {
	vals := map[string]SendActivityLogsStatus{
		"disabled": SendActivityLogsStatusDisabled,
		"enabled":  SendActivityLogsStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SendActivityLogsStatus(input)
	return &out, nil
}

type SendMetricsStatus string

const (
	SendMetricsStatusDisabled SendMetricsStatus = "Disabled"
	SendMetricsStatusEnabled  SendMetricsStatus = "Enabled"
)

func PossibleValuesForSendMetricsStatus() []string {
	return []string{
		string(SendMetricsStatusDisabled),
		string(SendMetricsStatusEnabled),
	}
}

func (s *SendMetricsStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSendMetricsStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSendMetricsStatus(input string) (*SendMetricsStatus, error) {
	vals := map[string]SendMetricsStatus{
		"disabled": SendMetricsStatusDisabled,
		"enabled":  SendMetricsStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SendMetricsStatus(input)
	return &out, nil
}

type SendSubscriptionLogsStatus string

const (
	SendSubscriptionLogsStatusDisabled SendSubscriptionLogsStatus = "Disabled"
	SendSubscriptionLogsStatusEnabled  SendSubscriptionLogsStatus = "Enabled"
)

func PossibleValuesForSendSubscriptionLogsStatus() []string {
	return []string{
		string(SendSubscriptionLogsStatusDisabled),
		string(SendSubscriptionLogsStatusEnabled),
	}
}

func (s *SendSubscriptionLogsStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSendSubscriptionLogsStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSendSubscriptionLogsStatus(input string) (*SendSubscriptionLogsStatus, error) {
	vals := map[string]SendSubscriptionLogsStatus{
		"disabled": SendSubscriptionLogsStatusDisabled,
		"enabled":  SendSubscriptionLogsStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SendSubscriptionLogsStatus(input)
	return &out, nil
}

type TagAction string

const (
	TagActionExclude TagAction = "Exclude"
	TagActionInclude TagAction = "Include"
)

func PossibleValuesForTagAction() []string {
	return []string{
		string(TagActionExclude),
		string(TagActionInclude),
	}
}

func (s *TagAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTagAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTagAction(input string) (*TagAction, error) {
	vals := map[string]TagAction{
		"exclude": TagActionExclude,
		"include": TagActionInclude,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TagAction(input)
	return &out, nil
}
