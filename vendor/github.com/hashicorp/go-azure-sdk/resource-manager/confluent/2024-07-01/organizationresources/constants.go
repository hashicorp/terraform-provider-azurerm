package organizationresources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisionState string

const (
	ProvisionStateAccepted     ProvisionState = "Accepted"
	ProvisionStateCanceled     ProvisionState = "Canceled"
	ProvisionStateCreating     ProvisionState = "Creating"
	ProvisionStateDeleted      ProvisionState = "Deleted"
	ProvisionStateDeleting     ProvisionState = "Deleting"
	ProvisionStateFailed       ProvisionState = "Failed"
	ProvisionStateNotSpecified ProvisionState = "NotSpecified"
	ProvisionStateSucceeded    ProvisionState = "Succeeded"
	ProvisionStateUpdating     ProvisionState = "Updating"
)

func PossibleValuesForProvisionState() []string {
	return []string{
		string(ProvisionStateAccepted),
		string(ProvisionStateCanceled),
		string(ProvisionStateCreating),
		string(ProvisionStateDeleted),
		string(ProvisionStateDeleting),
		string(ProvisionStateFailed),
		string(ProvisionStateNotSpecified),
		string(ProvisionStateSucceeded),
		string(ProvisionStateUpdating),
	}
}

func (s *ProvisionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisionState(input string) (*ProvisionState, error) {
	vals := map[string]ProvisionState{
		"accepted":     ProvisionStateAccepted,
		"canceled":     ProvisionStateCanceled,
		"creating":     ProvisionStateCreating,
		"deleted":      ProvisionStateDeleted,
		"deleting":     ProvisionStateDeleting,
		"failed":       ProvisionStateFailed,
		"notspecified": ProvisionStateNotSpecified,
		"succeeded":    ProvisionStateSucceeded,
		"updating":     ProvisionStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisionState(input)
	return &out, nil
}

type SaaSOfferStatus string

const (
	SaaSOfferStatusFailed                  SaaSOfferStatus = "Failed"
	SaaSOfferStatusInProgress              SaaSOfferStatus = "InProgress"
	SaaSOfferStatusPendingFulfillmentStart SaaSOfferStatus = "PendingFulfillmentStart"
	SaaSOfferStatusReinstated              SaaSOfferStatus = "Reinstated"
	SaaSOfferStatusStarted                 SaaSOfferStatus = "Started"
	SaaSOfferStatusSubscribed              SaaSOfferStatus = "Subscribed"
	SaaSOfferStatusSucceeded               SaaSOfferStatus = "Succeeded"
	SaaSOfferStatusSuspended               SaaSOfferStatus = "Suspended"
	SaaSOfferStatusUnsubscribed            SaaSOfferStatus = "Unsubscribed"
	SaaSOfferStatusUpdating                SaaSOfferStatus = "Updating"
)

func PossibleValuesForSaaSOfferStatus() []string {
	return []string{
		string(SaaSOfferStatusFailed),
		string(SaaSOfferStatusInProgress),
		string(SaaSOfferStatusPendingFulfillmentStart),
		string(SaaSOfferStatusReinstated),
		string(SaaSOfferStatusStarted),
		string(SaaSOfferStatusSubscribed),
		string(SaaSOfferStatusSucceeded),
		string(SaaSOfferStatusSuspended),
		string(SaaSOfferStatusUnsubscribed),
		string(SaaSOfferStatusUpdating),
	}
}

func (s *SaaSOfferStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSaaSOfferStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSaaSOfferStatus(input string) (*SaaSOfferStatus, error) {
	vals := map[string]SaaSOfferStatus{
		"failed":                  SaaSOfferStatusFailed,
		"inprogress":              SaaSOfferStatusInProgress,
		"pendingfulfillmentstart": SaaSOfferStatusPendingFulfillmentStart,
		"reinstated":              SaaSOfferStatusReinstated,
		"started":                 SaaSOfferStatusStarted,
		"subscribed":              SaaSOfferStatusSubscribed,
		"succeeded":               SaaSOfferStatusSucceeded,
		"suspended":               SaaSOfferStatusSuspended,
		"unsubscribed":            SaaSOfferStatusUnsubscribed,
		"updating":                SaaSOfferStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SaaSOfferStatus(input)
	return &out, nil
}
