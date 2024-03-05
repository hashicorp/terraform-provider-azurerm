package customipprefixes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommissionedState string

const (
	CommissionedStateCommissioned                    CommissionedState = "Commissioned"
	CommissionedStateCommissionedNoInternetAdvertise CommissionedState = "CommissionedNoInternetAdvertise"
	CommissionedStateCommissioning                   CommissionedState = "Commissioning"
	CommissionedStateDecommissioning                 CommissionedState = "Decommissioning"
	CommissionedStateDeprovisioned                   CommissionedState = "Deprovisioned"
	CommissionedStateDeprovisioning                  CommissionedState = "Deprovisioning"
	CommissionedStateProvisioned                     CommissionedState = "Provisioned"
	CommissionedStateProvisioning                    CommissionedState = "Provisioning"
)

func PossibleValuesForCommissionedState() []string {
	return []string{
		string(CommissionedStateCommissioned),
		string(CommissionedStateCommissionedNoInternetAdvertise),
		string(CommissionedStateCommissioning),
		string(CommissionedStateDecommissioning),
		string(CommissionedStateDeprovisioned),
		string(CommissionedStateDeprovisioning),
		string(CommissionedStateProvisioned),
		string(CommissionedStateProvisioning),
	}
}

func (s *CommissionedState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCommissionedState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCommissionedState(input string) (*CommissionedState, error) {
	vals := map[string]CommissionedState{
		"commissioned":                    CommissionedStateCommissioned,
		"commissionednointernetadvertise": CommissionedStateCommissionedNoInternetAdvertise,
		"commissioning":                   CommissionedStateCommissioning,
		"decommissioning":                 CommissionedStateDecommissioning,
		"deprovisioned":                   CommissionedStateDeprovisioned,
		"deprovisioning":                  CommissionedStateDeprovisioning,
		"provisioned":                     CommissionedStateProvisioned,
		"provisioning":                    CommissionedStateProvisioning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CommissionedState(input)
	return &out, nil
}

type CustomIPPrefixType string

const (
	CustomIPPrefixTypeChild    CustomIPPrefixType = "Child"
	CustomIPPrefixTypeParent   CustomIPPrefixType = "Parent"
	CustomIPPrefixTypeSingular CustomIPPrefixType = "Singular"
)

func PossibleValuesForCustomIPPrefixType() []string {
	return []string{
		string(CustomIPPrefixTypeChild),
		string(CustomIPPrefixTypeParent),
		string(CustomIPPrefixTypeSingular),
	}
}

func (s *CustomIPPrefixType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomIPPrefixType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomIPPrefixType(input string) (*CustomIPPrefixType, error) {
	vals := map[string]CustomIPPrefixType{
		"child":    CustomIPPrefixTypeChild,
		"parent":   CustomIPPrefixTypeParent,
		"singular": CustomIPPrefixTypeSingular,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomIPPrefixType(input)
	return &out, nil
}

type Geo string

const (
	GeoAFRI    Geo = "AFRI"
	GeoAPAC    Geo = "APAC"
	GeoAQ      Geo = "AQ"
	GeoEURO    Geo = "EURO"
	GeoGLOBAL  Geo = "GLOBAL"
	GeoLATAM   Geo = "LATAM"
	GeoME      Geo = "ME"
	GeoNAM     Geo = "NAM"
	GeoOCEANIA Geo = "OCEANIA"
)

func PossibleValuesForGeo() []string {
	return []string{
		string(GeoAFRI),
		string(GeoAPAC),
		string(GeoAQ),
		string(GeoEURO),
		string(GeoGLOBAL),
		string(GeoLATAM),
		string(GeoME),
		string(GeoNAM),
		string(GeoOCEANIA),
	}
}

func (s *Geo) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGeo(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGeo(input string) (*Geo, error) {
	vals := map[string]Geo{
		"afri":    GeoAFRI,
		"apac":    GeoAPAC,
		"aq":      GeoAQ,
		"euro":    GeoEURO,
		"global":  GeoGLOBAL,
		"latam":   GeoLATAM,
		"me":      GeoME,
		"nam":     GeoNAM,
		"oceania": GeoOCEANIA,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Geo(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
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
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
