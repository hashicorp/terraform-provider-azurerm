package brokerauthorization

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BrokerResourceDefinitionMethods string

const (
	BrokerResourceDefinitionMethodsConnect   BrokerResourceDefinitionMethods = "Connect"
	BrokerResourceDefinitionMethodsPublish   BrokerResourceDefinitionMethods = "Publish"
	BrokerResourceDefinitionMethodsSubscribe BrokerResourceDefinitionMethods = "Subscribe"
)

func PossibleValuesForBrokerResourceDefinitionMethods() []string {
	return []string{
		string(BrokerResourceDefinitionMethodsConnect),
		string(BrokerResourceDefinitionMethodsPublish),
		string(BrokerResourceDefinitionMethodsSubscribe),
	}
}

func (s *BrokerResourceDefinitionMethods) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBrokerResourceDefinitionMethods(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBrokerResourceDefinitionMethods(input string) (*BrokerResourceDefinitionMethods, error) {
	vals := map[string]BrokerResourceDefinitionMethods{
		"connect":   BrokerResourceDefinitionMethodsConnect,
		"publish":   BrokerResourceDefinitionMethodsPublish,
		"subscribe": BrokerResourceDefinitionMethodsSubscribe,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BrokerResourceDefinitionMethods(input)
	return &out, nil
}

type ExtendedLocationType string

const (
	ExtendedLocationTypeCustomLocation ExtendedLocationType = "CustomLocation"
)

func PossibleValuesForExtendedLocationType() []string {
	return []string{
		string(ExtendedLocationTypeCustomLocation),
	}
}

func (s *ExtendedLocationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExtendedLocationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExtendedLocationType(input string) (*ExtendedLocationType, error) {
	vals := map[string]ExtendedLocationType{
		"customlocation": ExtendedLocationTypeCustomLocation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtendedLocationType(input)
	return &out, nil
}

type OperationalMode string

const (
	OperationalModeDisabled OperationalMode = "Disabled"
	OperationalModeEnabled  OperationalMode = "Enabled"
)

func PossibleValuesForOperationalMode() []string {
	return []string{
		string(OperationalModeDisabled),
		string(OperationalModeEnabled),
	}
}

func (s *OperationalMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationalMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationalMode(input string) (*OperationalMode, error) {
	vals := map[string]OperationalMode{
		"disabled": OperationalModeDisabled,
		"enabled":  OperationalModeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationalMode(input)
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

type StateStoreResourceDefinitionMethods string

const (
	StateStoreResourceDefinitionMethodsRead      StateStoreResourceDefinitionMethods = "Read"
	StateStoreResourceDefinitionMethodsReadWrite StateStoreResourceDefinitionMethods = "ReadWrite"
	StateStoreResourceDefinitionMethodsWrite     StateStoreResourceDefinitionMethods = "Write"
)

func PossibleValuesForStateStoreResourceDefinitionMethods() []string {
	return []string{
		string(StateStoreResourceDefinitionMethodsRead),
		string(StateStoreResourceDefinitionMethodsReadWrite),
		string(StateStoreResourceDefinitionMethodsWrite),
	}
}

func (s *StateStoreResourceDefinitionMethods) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStateStoreResourceDefinitionMethods(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStateStoreResourceDefinitionMethods(input string) (*StateStoreResourceDefinitionMethods, error) {
	vals := map[string]StateStoreResourceDefinitionMethods{
		"read":      StateStoreResourceDefinitionMethodsRead,
		"readwrite": StateStoreResourceDefinitionMethodsReadWrite,
		"write":     StateStoreResourceDefinitionMethodsWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StateStoreResourceDefinitionMethods(input)
	return &out, nil
}

type StateStoreResourceKeyTypes string

const (
	StateStoreResourceKeyTypesBinary  StateStoreResourceKeyTypes = "Binary"
	StateStoreResourceKeyTypesPattern StateStoreResourceKeyTypes = "Pattern"
	StateStoreResourceKeyTypesString  StateStoreResourceKeyTypes = "String"
)

func PossibleValuesForStateStoreResourceKeyTypes() []string {
	return []string{
		string(StateStoreResourceKeyTypesBinary),
		string(StateStoreResourceKeyTypesPattern),
		string(StateStoreResourceKeyTypesString),
	}
}

func (s *StateStoreResourceKeyTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStateStoreResourceKeyTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStateStoreResourceKeyTypes(input string) (*StateStoreResourceKeyTypes, error) {
	vals := map[string]StateStoreResourceKeyTypes{
		"binary":  StateStoreResourceKeyTypesBinary,
		"pattern": StateStoreResourceKeyTypesPattern,
		"string":  StateStoreResourceKeyTypesString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StateStoreResourceKeyTypes(input)
	return &out, nil
}
