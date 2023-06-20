package recordsets

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecordType string

const (
	RecordTypeA     RecordType = "A"
	RecordTypeAAAA  RecordType = "AAAA"
	RecordTypeCNAME RecordType = "CNAME"
	RecordTypeMX    RecordType = "MX"
	RecordTypePTR   RecordType = "PTR"
	RecordTypeSOA   RecordType = "SOA"
	RecordTypeSRV   RecordType = "SRV"
	RecordTypeTXT   RecordType = "TXT"
)

func PossibleValuesForRecordType() []string {
	return []string{
		string(RecordTypeA),
		string(RecordTypeAAAA),
		string(RecordTypeCNAME),
		string(RecordTypeMX),
		string(RecordTypePTR),
		string(RecordTypeSOA),
		string(RecordTypeSRV),
		string(RecordTypeTXT),
	}
}

func parseRecordType(input string) (*RecordType, error) {
	vals := map[string]RecordType{
		"a":     RecordTypeA,
		"aaaa":  RecordTypeAAAA,
		"cname": RecordTypeCNAME,
		"mx":    RecordTypeMX,
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
