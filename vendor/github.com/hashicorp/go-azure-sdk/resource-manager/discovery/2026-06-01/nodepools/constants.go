package nodepools

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

type ScaleSetPriority string

const (
	ScaleSetPriorityRegular ScaleSetPriority = "Regular"
	ScaleSetPrioritySpot    ScaleSetPriority = "Spot"
)

func PossibleValuesForScaleSetPriority() []string {
	return []string{
		string(ScaleSetPriorityRegular),
		string(ScaleSetPrioritySpot),
	}
}

func (s *ScaleSetPriority) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScaleSetPriority(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScaleSetPriority(input string) (*ScaleSetPriority, error) {
	vals := map[string]ScaleSetPriority{
		"regular": ScaleSetPriorityRegular,
		"spot":    ScaleSetPrioritySpot,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScaleSetPriority(input)
	return &out, nil
}

type VMSize string

const (
	VMSizeStandardNCEightasTFourVThree           VMSize = "Standard_NC8as_T4_v3"
	VMSizeStandardNCFourEightadsAOneHundredVFour VMSize = "Standard_NC48ads_A100_v4"
	VMSizeStandardNCFourasTFourVThree            VMSize = "Standard_NC4as_T4_v3"
	VMSizeStandardNCNineSixadsAOneHundredVFour   VMSize = "Standard_NC96ads_A100_v4"
	VMSizeStandardNCOneSixasTFourVThree          VMSize = "Standard_NC16as_T4_v3"
	VMSizeStandardNCSixFourasTFourVThree         VMSize = "Standard_NC64as_T4_v3"
	VMSizeStandardNCTwoFouradsAOneHundredVFour   VMSize = "Standard_NC24ads_A100_v4"
	VMSizeStandardNDFourZerorsVTwo               VMSize = "Standard_ND40rs_v2"
	VMSizeStandardNVOneTwoadsAOneZeroVFive       VMSize = "Standard_NV12ads_A10_v5"
	VMSizeStandardNVSevenTwoadsAOneZeroVFive     VMSize = "Standard_NV72ads_A10_v5"
	VMSizeStandardNVSixadsAOneZeroVFive          VMSize = "Standard_NV6ads_A10_v5"
	VMSizeStandardNVThreeSixadmsAOneZeroVFive    VMSize = "Standard_NV36adms_A10_v5"
	VMSizeStandardNVThreeSixadsAOneZeroVFive     VMSize = "Standard_NV36ads_A10_v5"
	VMSizeStandardNVTwoFouradsAOneZeroVFive      VMSize = "Standard_NV24ads_A10_v5"
)

func PossibleValuesForVMSize() []string {
	return []string{
		string(VMSizeStandardNCEightasTFourVThree),
		string(VMSizeStandardNCFourEightadsAOneHundredVFour),
		string(VMSizeStandardNCFourasTFourVThree),
		string(VMSizeStandardNCNineSixadsAOneHundredVFour),
		string(VMSizeStandardNCOneSixasTFourVThree),
		string(VMSizeStandardNCSixFourasTFourVThree),
		string(VMSizeStandardNCTwoFouradsAOneHundredVFour),
		string(VMSizeStandardNDFourZerorsVTwo),
		string(VMSizeStandardNVOneTwoadsAOneZeroVFive),
		string(VMSizeStandardNVSevenTwoadsAOneZeroVFive),
		string(VMSizeStandardNVSixadsAOneZeroVFive),
		string(VMSizeStandardNVThreeSixadmsAOneZeroVFive),
		string(VMSizeStandardNVThreeSixadsAOneZeroVFive),
		string(VMSizeStandardNVTwoFouradsAOneZeroVFive),
	}
}

func (s *VMSize) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMSize(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMSize(input string) (*VMSize, error) {
	vals := map[string]VMSize{
		"standard_nc8as_t4_v3":     VMSizeStandardNCEightasTFourVThree,
		"standard_nc48ads_a100_v4": VMSizeStandardNCFourEightadsAOneHundredVFour,
		"standard_nc4as_t4_v3":     VMSizeStandardNCFourasTFourVThree,
		"standard_nc96ads_a100_v4": VMSizeStandardNCNineSixadsAOneHundredVFour,
		"standard_nc16as_t4_v3":    VMSizeStandardNCOneSixasTFourVThree,
		"standard_nc64as_t4_v3":    VMSizeStandardNCSixFourasTFourVThree,
		"standard_nc24ads_a100_v4": VMSizeStandardNCTwoFouradsAOneHundredVFour,
		"standard_nd40rs_v2":       VMSizeStandardNDFourZerorsVTwo,
		"standard_nv12ads_a10_v5":  VMSizeStandardNVOneTwoadsAOneZeroVFive,
		"standard_nv72ads_a10_v5":  VMSizeStandardNVSevenTwoadsAOneZeroVFive,
		"standard_nv6ads_a10_v5":   VMSizeStandardNVSixadsAOneZeroVFive,
		"standard_nv36adms_a10_v5": VMSizeStandardNVThreeSixadmsAOneZeroVFive,
		"standard_nv36ads_a10_v5":  VMSizeStandardNVThreeSixadsAOneZeroVFive,
		"standard_nv24ads_a10_v5":  VMSizeStandardNVTwoFouradsAOneZeroVFive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMSize(input)
	return &out, nil
}
