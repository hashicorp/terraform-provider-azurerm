package dedicatedhsms

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JsonWebKeyType string

const (
	JsonWebKeyTypeAllocating    JsonWebKeyType = "Allocating"
	JsonWebKeyTypeCheckingQuota JsonWebKeyType = "CheckingQuota"
	JsonWebKeyTypeConnecting    JsonWebKeyType = "Connecting"
	JsonWebKeyTypeDeleting      JsonWebKeyType = "Deleting"
	JsonWebKeyTypeFailed        JsonWebKeyType = "Failed"
	JsonWebKeyTypeProvisioning  JsonWebKeyType = "Provisioning"
	JsonWebKeyTypeSucceeded     JsonWebKeyType = "Succeeded"
)

func PossibleValuesForJsonWebKeyType() []string {
	return []string{
		string(JsonWebKeyTypeAllocating),
		string(JsonWebKeyTypeCheckingQuota),
		string(JsonWebKeyTypeConnecting),
		string(JsonWebKeyTypeDeleting),
		string(JsonWebKeyTypeFailed),
		string(JsonWebKeyTypeProvisioning),
		string(JsonWebKeyTypeSucceeded),
	}
}

func (s *JsonWebKeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJsonWebKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJsonWebKeyType(input string) (*JsonWebKeyType, error) {
	vals := map[string]JsonWebKeyType{
		"allocating":    JsonWebKeyTypeAllocating,
		"checkingquota": JsonWebKeyTypeCheckingQuota,
		"connecting":    JsonWebKeyTypeConnecting,
		"deleting":      JsonWebKeyTypeDeleting,
		"failed":        JsonWebKeyTypeFailed,
		"provisioning":  JsonWebKeyTypeProvisioning,
		"succeeded":     JsonWebKeyTypeSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JsonWebKeyType(input)
	return &out, nil
}

type SkuName string

const (
	SkuNamePayShieldOneZeroKLMKOneCPSSixZero         SkuName = "payShield10K_LMK1_CPS60"
	SkuNamePayShieldOneZeroKLMKOneCPSTwoFiveZero     SkuName = "payShield10K_LMK1_CPS250"
	SkuNamePayShieldOneZeroKLMKOneCPSTwoFiveZeroZero SkuName = "payShield10K_LMK1_CPS2500"
	SkuNamePayShieldOneZeroKLMKTwoCPSSixZero         SkuName = "payShield10K_LMK2_CPS60"
	SkuNamePayShieldOneZeroKLMKTwoCPSTwoFiveZero     SkuName = "payShield10K_LMK2_CPS250"
	SkuNamePayShieldOneZeroKLMKTwoCPSTwoFiveZeroZero SkuName = "payShield10K_LMK2_CPS2500"
	SkuNameSafeNetLunaNetworkHSMASevenNineZero       SkuName = "SafeNet Luna Network HSM A790"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNamePayShieldOneZeroKLMKOneCPSSixZero),
		string(SkuNamePayShieldOneZeroKLMKOneCPSTwoFiveZero),
		string(SkuNamePayShieldOneZeroKLMKOneCPSTwoFiveZeroZero),
		string(SkuNamePayShieldOneZeroKLMKTwoCPSSixZero),
		string(SkuNamePayShieldOneZeroKLMKTwoCPSTwoFiveZero),
		string(SkuNamePayShieldOneZeroKLMKTwoCPSTwoFiveZeroZero),
		string(SkuNameSafeNetLunaNetworkHSMASevenNineZero),
	}
}

func (s *SkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"payshield10k_lmk1_cps60":       SkuNamePayShieldOneZeroKLMKOneCPSSixZero,
		"payshield10k_lmk1_cps250":      SkuNamePayShieldOneZeroKLMKOneCPSTwoFiveZero,
		"payshield10k_lmk1_cps2500":     SkuNamePayShieldOneZeroKLMKOneCPSTwoFiveZeroZero,
		"payshield10k_lmk2_cps60":       SkuNamePayShieldOneZeroKLMKTwoCPSSixZero,
		"payshield10k_lmk2_cps250":      SkuNamePayShieldOneZeroKLMKTwoCPSTwoFiveZero,
		"payshield10k_lmk2_cps2500":     SkuNamePayShieldOneZeroKLMKTwoCPSTwoFiveZeroZero,
		"safenet luna network hsm a790": SkuNameSafeNetLunaNetworkHSMASevenNineZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}
