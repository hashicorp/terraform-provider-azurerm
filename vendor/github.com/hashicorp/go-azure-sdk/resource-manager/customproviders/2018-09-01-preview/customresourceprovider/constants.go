package customresourceprovider

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionRouting string

const (
	ActionRoutingProxy ActionRouting = "Proxy"
)

func PossibleValuesForActionRouting() []string {
	return []string{
		string(ActionRoutingProxy),
	}
}

func parseActionRouting(input string) (*ActionRouting, error) {
	vals := map[string]ActionRouting{
		"proxy": ActionRoutingProxy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActionRouting(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateRunning   ProvisioningState = "Running"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateRunning),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":  ProvisioningStateAccepted,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"running":   ProvisioningStateRunning,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type ResourceTypeRouting string

const (
	ResourceTypeRoutingProxy      ResourceTypeRouting = "Proxy"
	ResourceTypeRoutingProxyCache ResourceTypeRouting = "Proxy,Cache"
)

func PossibleValuesForResourceTypeRouting() []string {
	return []string{
		string(ResourceTypeRoutingProxy),
		string(ResourceTypeRoutingProxyCache),
	}
}

func parseResourceTypeRouting(input string) (*ResourceTypeRouting, error) {
	vals := map[string]ResourceTypeRouting{
		"proxy":       ResourceTypeRoutingProxy,
		"proxy,cache": ResourceTypeRoutingProxyCache,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceTypeRouting(input)
	return &out, nil
}

type ValidationType string

const (
	ValidationTypeSwagger ValidationType = "Swagger"
)

func PossibleValuesForValidationType() []string {
	return []string{
		string(ValidationTypeSwagger),
	}
}

func parseValidationType(input string) (*ValidationType, error) {
	vals := map[string]ValidationType{
		"swagger": ValidationTypeSwagger,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ValidationType(input)
	return &out, nil
}
