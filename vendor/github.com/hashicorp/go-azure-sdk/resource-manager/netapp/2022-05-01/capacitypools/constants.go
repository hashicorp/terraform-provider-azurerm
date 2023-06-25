package capacitypools

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionType string

const (
	EncryptionTypeDouble EncryptionType = "Double"
	EncryptionTypeSingle EncryptionType = "Single"
)

func PossibleValuesForEncryptionType() []string {
	return []string{
		string(EncryptionTypeDouble),
		string(EncryptionTypeSingle),
	}
}

func parseEncryptionType(input string) (*EncryptionType, error) {
	vals := map[string]EncryptionType{
		"double": EncryptionTypeDouble,
		"single": EncryptionTypeSingle,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionType(input)
	return &out, nil
}

type QosType string

const (
	QosTypeAuto   QosType = "Auto"
	QosTypeManual QosType = "Manual"
)

func PossibleValuesForQosType() []string {
	return []string{
		string(QosTypeAuto),
		string(QosTypeManual),
	}
}

func parseQosType(input string) (*QosType, error) {
	vals := map[string]QosType{
		"auto":   QosTypeAuto,
		"manual": QosTypeManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QosType(input)
	return &out, nil
}

type ServiceLevel string

const (
	ServiceLevelPremium     ServiceLevel = "Premium"
	ServiceLevelStandard    ServiceLevel = "Standard"
	ServiceLevelStandardZRS ServiceLevel = "StandardZRS"
	ServiceLevelUltra       ServiceLevel = "Ultra"
)

func PossibleValuesForServiceLevel() []string {
	return []string{
		string(ServiceLevelPremium),
		string(ServiceLevelStandard),
		string(ServiceLevelStandardZRS),
		string(ServiceLevelUltra),
	}
}

func parseServiceLevel(input string) (*ServiceLevel, error) {
	vals := map[string]ServiceLevel{
		"premium":     ServiceLevelPremium,
		"standard":    ServiceLevelStandard,
		"standardzrs": ServiceLevelStandardZRS,
		"ultra":       ServiceLevelUltra,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceLevel(input)
	return &out, nil
}
