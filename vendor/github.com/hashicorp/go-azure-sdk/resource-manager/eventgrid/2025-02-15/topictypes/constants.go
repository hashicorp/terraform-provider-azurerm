package topictypes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceRegionType string

const (
	ResourceRegionTypeGlobalResource   ResourceRegionType = "GlobalResource"
	ResourceRegionTypeRegionalResource ResourceRegionType = "RegionalResource"
)

func PossibleValuesForResourceRegionType() []string {
	return []string{
		string(ResourceRegionTypeGlobalResource),
		string(ResourceRegionTypeRegionalResource),
	}
}

func (s *ResourceRegionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceRegionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceRegionType(input string) (*ResourceRegionType, error) {
	vals := map[string]ResourceRegionType{
		"globalresource":   ResourceRegionTypeGlobalResource,
		"regionalresource": ResourceRegionTypeRegionalResource,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceRegionType(input)
	return &out, nil
}

type TopicTypeProvisioningState string

const (
	TopicTypeProvisioningStateCanceled  TopicTypeProvisioningState = "Canceled"
	TopicTypeProvisioningStateCreating  TopicTypeProvisioningState = "Creating"
	TopicTypeProvisioningStateDeleting  TopicTypeProvisioningState = "Deleting"
	TopicTypeProvisioningStateFailed    TopicTypeProvisioningState = "Failed"
	TopicTypeProvisioningStateSucceeded TopicTypeProvisioningState = "Succeeded"
	TopicTypeProvisioningStateUpdating  TopicTypeProvisioningState = "Updating"
)

func PossibleValuesForTopicTypeProvisioningState() []string {
	return []string{
		string(TopicTypeProvisioningStateCanceled),
		string(TopicTypeProvisioningStateCreating),
		string(TopicTypeProvisioningStateDeleting),
		string(TopicTypeProvisioningStateFailed),
		string(TopicTypeProvisioningStateSucceeded),
		string(TopicTypeProvisioningStateUpdating),
	}
}

func (s *TopicTypeProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTopicTypeProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTopicTypeProvisioningState(input string) (*TopicTypeProvisioningState, error) {
	vals := map[string]TopicTypeProvisioningState{
		"canceled":  TopicTypeProvisioningStateCanceled,
		"creating":  TopicTypeProvisioningStateCreating,
		"deleting":  TopicTypeProvisioningStateDeleting,
		"failed":    TopicTypeProvisioningStateFailed,
		"succeeded": TopicTypeProvisioningStateSucceeded,
		"updating":  TopicTypeProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TopicTypeProvisioningState(input)
	return &out, nil
}

type TopicTypeSourceScope string

const (
	TopicTypeSourceScopeAzureSubscription TopicTypeSourceScope = "AzureSubscription"
	TopicTypeSourceScopeManagementGroup   TopicTypeSourceScope = "ManagementGroup"
	TopicTypeSourceScopeResource          TopicTypeSourceScope = "Resource"
	TopicTypeSourceScopeResourceGroup     TopicTypeSourceScope = "ResourceGroup"
)

func PossibleValuesForTopicTypeSourceScope() []string {
	return []string{
		string(TopicTypeSourceScopeAzureSubscription),
		string(TopicTypeSourceScopeManagementGroup),
		string(TopicTypeSourceScopeResource),
		string(TopicTypeSourceScopeResourceGroup),
	}
}

func (s *TopicTypeSourceScope) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTopicTypeSourceScope(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTopicTypeSourceScope(input string) (*TopicTypeSourceScope, error) {
	vals := map[string]TopicTypeSourceScope{
		"azuresubscription": TopicTypeSourceScopeAzureSubscription,
		"managementgroup":   TopicTypeSourceScopeManagementGroup,
		"resource":          TopicTypeSourceScopeResource,
		"resourcegroup":     TopicTypeSourceScopeResourceGroup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TopicTypeSourceScope(input)
	return &out, nil
}
