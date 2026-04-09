package topicspaces

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicSpaceProvisioningState string

const (
	TopicSpaceProvisioningStateCanceled  TopicSpaceProvisioningState = "Canceled"
	TopicSpaceProvisioningStateCreating  TopicSpaceProvisioningState = "Creating"
	TopicSpaceProvisioningStateDeleted   TopicSpaceProvisioningState = "Deleted"
	TopicSpaceProvisioningStateDeleting  TopicSpaceProvisioningState = "Deleting"
	TopicSpaceProvisioningStateFailed    TopicSpaceProvisioningState = "Failed"
	TopicSpaceProvisioningStateSucceeded TopicSpaceProvisioningState = "Succeeded"
	TopicSpaceProvisioningStateUpdating  TopicSpaceProvisioningState = "Updating"
)

func PossibleValuesForTopicSpaceProvisioningState() []string {
	return []string{
		string(TopicSpaceProvisioningStateCanceled),
		string(TopicSpaceProvisioningStateCreating),
		string(TopicSpaceProvisioningStateDeleted),
		string(TopicSpaceProvisioningStateDeleting),
		string(TopicSpaceProvisioningStateFailed),
		string(TopicSpaceProvisioningStateSucceeded),
		string(TopicSpaceProvisioningStateUpdating),
	}
}

func (s *TopicSpaceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTopicSpaceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTopicSpaceProvisioningState(input string) (*TopicSpaceProvisioningState, error) {
	vals := map[string]TopicSpaceProvisioningState{
		"canceled":  TopicSpaceProvisioningStateCanceled,
		"creating":  TopicSpaceProvisioningStateCreating,
		"deleted":   TopicSpaceProvisioningStateDeleted,
		"deleting":  TopicSpaceProvisioningStateDeleting,
		"failed":    TopicSpaceProvisioningStateFailed,
		"succeeded": TopicSpaceProvisioningStateSucceeded,
		"updating":  TopicSpaceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TopicSpaceProvisioningState(input)
	return &out, nil
}
