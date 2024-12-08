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
	RecordTypeDS    RecordType = "DS"
	RecordTypeMX    RecordType = "MX"
	RecordTypeNAPTR RecordType = "NAPTR"
	RecordTypeNS    RecordType = "NS"
	RecordTypePTR   RecordType = "PTR"
	RecordTypeSOA   RecordType = "SOA"
	RecordTypeSRV   RecordType = "SRV"
	RecordTypeTLSA  RecordType = "TLSA"
	RecordTypeTXT   RecordType = "TXT"
)

func PossibleValuesForRecordType() []string {
	return []string{
		string(RecordTypeA),
		string(RecordTypeAAAA),
		string(RecordTypeCAA),
		string(RecordTypeCNAME),
		string(RecordTypeDS),
		string(RecordTypeMX),
		string(RecordTypeNAPTR),
		string(RecordTypeNS),
		string(RecordTypePTR),
		string(RecordTypeSOA),
		string(RecordTypeSRV),
		string(RecordTypeTLSA),
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
		"ds":    RecordTypeDS,
		"mx":    RecordTypeMX,
		"naptr": RecordTypeNAPTR,
		"ns":    RecordTypeNS,
		"ptr":   RecordTypePTR,
		"soa":   RecordTypeSOA,
		"srv":   RecordTypeSRV,
		"tlsa":  RecordTypeTLSA,
		"txt":   RecordTypeTXT,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RecordType(input)
	return &out, nil
}
