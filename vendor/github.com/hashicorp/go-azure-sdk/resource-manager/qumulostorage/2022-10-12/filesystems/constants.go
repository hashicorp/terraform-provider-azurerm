package filesystems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MarketplaceSubscriptionStatus string

const (
	MarketplaceSubscriptionStatusPendingFulfillmentStart MarketplaceSubscriptionStatus = "PendingFulfillmentStart"
	MarketplaceSubscriptionStatusSubscribed              MarketplaceSubscriptionStatus = "Subscribed"
	MarketplaceSubscriptionStatusSuspended               MarketplaceSubscriptionStatus = "Suspended"
	MarketplaceSubscriptionStatusUnsubscribed            MarketplaceSubscriptionStatus = "Unsubscribed"
)

func PossibleValuesForMarketplaceSubscriptionStatus() []string {
	return []string{
		string(MarketplaceSubscriptionStatusPendingFulfillmentStart),
		string(MarketplaceSubscriptionStatusSubscribed),
		string(MarketplaceSubscriptionStatusSuspended),
		string(MarketplaceSubscriptionStatusUnsubscribed),
	}
}

func (s *MarketplaceSubscriptionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMarketplaceSubscriptionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMarketplaceSubscriptionStatus(input string) (*MarketplaceSubscriptionStatus, error) {
	vals := map[string]MarketplaceSubscriptionStatus{
		"pendingfulfillmentstart": MarketplaceSubscriptionStatusPendingFulfillmentStart,
		"subscribed":              MarketplaceSubscriptionStatusSubscribed,
		"suspended":               MarketplaceSubscriptionStatusSuspended,
		"unsubscribed":            MarketplaceSubscriptionStatusUnsubscribed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MarketplaceSubscriptionStatus(input)
	return &out, nil
}

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

type StorageSku string

const (
	StorageSkuPerformance StorageSku = "Performance"
	StorageSkuStandard    StorageSku = "Standard"
)

func PossibleValuesForStorageSku() []string {
	return []string{
		string(StorageSkuPerformance),
		string(StorageSkuStandard),
	}
}

func (s *StorageSku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageSku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageSku(input string) (*StorageSku, error) {
	vals := map[string]StorageSku{
		"performance": StorageSkuPerformance,
		"standard":    StorageSkuStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageSku(input)
	return &out, nil
}
