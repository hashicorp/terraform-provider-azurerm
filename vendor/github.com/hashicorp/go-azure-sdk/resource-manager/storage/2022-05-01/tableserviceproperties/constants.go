package tableserviceproperties

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AllowedMethods string

const (
	AllowedMethodsDELETE  AllowedMethods = "DELETE"
	AllowedMethodsGET     AllowedMethods = "GET"
	AllowedMethodsHEAD    AllowedMethods = "HEAD"
	AllowedMethodsMERGE   AllowedMethods = "MERGE"
	AllowedMethodsOPTIONS AllowedMethods = "OPTIONS"
	AllowedMethodsPATCH   AllowedMethods = "PATCH"
	AllowedMethodsPOST    AllowedMethods = "POST"
	AllowedMethodsPUT     AllowedMethods = "PUT"
)

func PossibleValuesForAllowedMethods() []string {
	return []string{
		string(AllowedMethodsDELETE),
		string(AllowedMethodsGET),
		string(AllowedMethodsHEAD),
		string(AllowedMethodsMERGE),
		string(AllowedMethodsOPTIONS),
		string(AllowedMethodsPATCH),
		string(AllowedMethodsPOST),
		string(AllowedMethodsPUT),
	}
}

func parseAllowedMethods(input string) (*AllowedMethods, error) {
	vals := map[string]AllowedMethods{
		"delete":  AllowedMethodsDELETE,
		"get":     AllowedMethodsGET,
		"head":    AllowedMethodsHEAD,
		"merge":   AllowedMethodsMERGE,
		"options": AllowedMethodsOPTIONS,
		"patch":   AllowedMethodsPATCH,
		"post":    AllowedMethodsPOST,
		"put":     AllowedMethodsPUT,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AllowedMethods(input)
	return &out, nil
}
