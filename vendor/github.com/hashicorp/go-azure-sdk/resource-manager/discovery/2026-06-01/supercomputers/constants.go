package supercomputers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomerManagedKeys string

const (
	CustomerManagedKeysDisabled CustomerManagedKeys = "Disabled"
	CustomerManagedKeysEnabled  CustomerManagedKeys = "Enabled"
)

func PossibleValuesForCustomerManagedKeys() []string {
	return []string{
		string(CustomerManagedKeysDisabled),
		string(CustomerManagedKeysEnabled),
	}
}

func (s *CustomerManagedKeys) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomerManagedKeys(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomerManagedKeys(input string) (*CustomerManagedKeys, error) {
	vals := map[string]CustomerManagedKeys{
		"disabled": CustomerManagedKeysDisabled,
		"enabled":  CustomerManagedKeysEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomerManagedKeys(input)
	return &out, nil
}

type NetworkEgressType string

const (
	NetworkEgressTypeLoadBalancer NetworkEgressType = "LoadBalancer"
	NetworkEgressTypeNone         NetworkEgressType = "None"
)

func PossibleValuesForNetworkEgressType() []string {
	return []string{
		string(NetworkEgressTypeLoadBalancer),
		string(NetworkEgressTypeNone),
	}
}

func (s *NetworkEgressType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkEgressType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkEgressType(input string) (*NetworkEgressType, error) {
	vals := map[string]NetworkEgressType{
		"loadbalancer": NetworkEgressTypeLoadBalancer,
		"none":         NetworkEgressTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkEgressType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
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
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
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

type SystemAssignedServiceIdentityType string

const (
	SystemAssignedServiceIdentityTypeNone           SystemAssignedServiceIdentityType = "None"
	SystemAssignedServiceIdentityTypeSystemAssigned SystemAssignedServiceIdentityType = "SystemAssigned"
)

func PossibleValuesForSystemAssignedServiceIdentityType() []string {
	return []string{
		string(SystemAssignedServiceIdentityTypeNone),
		string(SystemAssignedServiceIdentityTypeSystemAssigned),
	}
}

func (s *SystemAssignedServiceIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSystemAssignedServiceIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSystemAssignedServiceIdentityType(input string) (*SystemAssignedServiceIdentityType, error) {
	vals := map[string]SystemAssignedServiceIdentityType{
		"none":           SystemAssignedServiceIdentityTypeNone,
		"systemassigned": SystemAssignedServiceIdentityTypeSystemAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SystemAssignedServiceIdentityType(input)
	return &out, nil
}

type SystemSku string

const (
	SystemSkuStandardDFoursVFive SystemSku = "Standard_D4s_v5"
	SystemSkuStandardDFoursVFour SystemSku = "Standard_D4s_v4"
	SystemSkuStandardDFoursVSix  SystemSku = "Standard_D4s_v6"
)

func PossibleValuesForSystemSku() []string {
	return []string{
		string(SystemSkuStandardDFoursVFive),
		string(SystemSkuStandardDFoursVFour),
		string(SystemSkuStandardDFoursVSix),
	}
}

func (s *SystemSku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSystemSku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSystemSku(input string) (*SystemSku, error) {
	vals := map[string]SystemSku{
		"standard_d4s_v5": SystemSkuStandardDFoursVFive,
		"standard_d4s_v4": SystemSkuStandardDFoursVFour,
		"standard_d4s_v6": SystemSkuStandardDFoursVSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SystemSku(input)
	return &out, nil
}
