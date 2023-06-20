package recordsets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecordType string

const (
	RecordTypeA     RecordType = "A"
	RecordTypeAAAA  RecordType = "AAAA"
	RecordTypeCAA   RecordType = "CAA"
	RecordTypeCNAME RecordType = "CNAME"
	RecordTypeMX    RecordType = "MX"
	RecordTypeNS    RecordType = "NS"
	RecordTypePTR   RecordType = "PTR"
	RecordTypeSOA   RecordType = "SOA"
	RecordTypeSRV   RecordType = "SRV"
	RecordTypeTXT   RecordType = "TXT"
)

func PossibleValuesForRecordType() []string {
	return []string{
		string(RecordTypeA),
		string(RecordTypeAAAA),
		string(RecordTypeCAA),
		string(RecordTypeCNAME),
		string(RecordTypeMX),
		string(RecordTypeNS),
		string(RecordTypePTR),
		string(RecordTypeSOA),
		string(RecordTypeSRV),
		string(RecordTypeTXT),
	}
}

func (s *RecordType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRecordType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRecordType(input string) (*RecordType, error) {
	vals := map[string]RecordType{
		"a":     RecordTypeA,
		"aaaa":  RecordTypeAAAA,
		"caa":   RecordTypeCAA,
		"cname": RecordTypeCNAME,
		"mx":    RecordTypeMX,
		"ns":    RecordTypeNS,
		"ptr":   RecordTypePTR,
		"soa":   RecordTypeSOA,
		"srv":   RecordTypeSRV,
		"txt":   RecordTypeTXT,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RecordType(input)
	return &out, nil
}
