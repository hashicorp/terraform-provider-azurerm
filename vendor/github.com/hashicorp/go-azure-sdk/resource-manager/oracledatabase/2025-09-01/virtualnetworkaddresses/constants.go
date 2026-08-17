package virtualnetworkaddresses

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureResourceProvisioningState string

const (
	AzureResourceProvisioningStateCanceled     AzureResourceProvisioningState = "Canceled"
	AzureResourceProvisioningStateFailed       AzureResourceProvisioningState = "Failed"
	AzureResourceProvisioningStateProvisioning AzureResourceProvisioningState = "Provisioning"
	AzureResourceProvisioningStateSucceeded    AzureResourceProvisioningState = "Succeeded"
)

func PossibleValuesForAzureResourceProvisioningState() []string {
	return []string{
		string(AzureResourceProvisioningStateCanceled),
		string(AzureResourceProvisioningStateFailed),
		string(AzureResourceProvisioningStateProvisioning),
		string(AzureResourceProvisioningStateSucceeded),
	}
}

func (s *AzureResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureResourceProvisioningState(input string) (*AzureResourceProvisioningState, error) {
	vals := map[string]AzureResourceProvisioningState{
		"canceled":     AzureResourceProvisioningStateCanceled,
		"failed":       AzureResourceProvisioningStateFailed,
		"provisioning": AzureResourceProvisioningStateProvisioning,
		"succeeded":    AzureResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureResourceProvisioningState(input)
	return &out, nil
}

type VirtualNetworkAddressLifecycleState string

const (
	VirtualNetworkAddressLifecycleStateAvailable    VirtualNetworkAddressLifecycleState = "Available"
	VirtualNetworkAddressLifecycleStateFailed       VirtualNetworkAddressLifecycleState = "Failed"
	VirtualNetworkAddressLifecycleStateProvisioning VirtualNetworkAddressLifecycleState = "Provisioning"
	VirtualNetworkAddressLifecycleStateTerminated   VirtualNetworkAddressLifecycleState = "Terminated"
	VirtualNetworkAddressLifecycleStateTerminating  VirtualNetworkAddressLifecycleState = "Terminating"
)

func PossibleValuesForVirtualNetworkAddressLifecycleState() []string {
	return []string{
		string(VirtualNetworkAddressLifecycleStateAvailable),
		string(VirtualNetworkAddressLifecycleStateFailed),
		string(VirtualNetworkAddressLifecycleStateProvisioning),
		string(VirtualNetworkAddressLifecycleStateTerminated),
		string(VirtualNetworkAddressLifecycleStateTerminating),
	}
}

func (s *VirtualNetworkAddressLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkAddressLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkAddressLifecycleState(input string) (*VirtualNetworkAddressLifecycleState, error) {
	vals := map[string]VirtualNetworkAddressLifecycleState{
		"available":    VirtualNetworkAddressLifecycleStateAvailable,
		"failed":       VirtualNetworkAddressLifecycleStateFailed,
		"provisioning": VirtualNetworkAddressLifecycleStateProvisioning,
		"terminated":   VirtualNetworkAddressLifecycleStateTerminated,
		"terminating":  VirtualNetworkAddressLifecycleStateTerminating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkAddressLifecycleState(input)
	return &out, nil
}
