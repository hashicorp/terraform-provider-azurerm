package apidefinitions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiSpecExportResultFormat string

const (
	ApiSpecExportResultFormatInline ApiSpecExportResultFormat = "inline"
	ApiSpecExportResultFormatLink   ApiSpecExportResultFormat = "link"
)

func PossibleValuesForApiSpecExportResultFormat() []string {
	return []string{
		string(ApiSpecExportResultFormatInline),
		string(ApiSpecExportResultFormatLink),
	}
}

func (s *ApiSpecExportResultFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApiSpecExportResultFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApiSpecExportResultFormat(input string) (*ApiSpecExportResultFormat, error) {
	vals := map[string]ApiSpecExportResultFormat{
		"inline": ApiSpecExportResultFormatInline,
		"link":   ApiSpecExportResultFormatLink,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApiSpecExportResultFormat(input)
	return &out, nil
}

type ApiSpecImportSourceFormat string

const (
	ApiSpecImportSourceFormatInline ApiSpecImportSourceFormat = "inline"
	ApiSpecImportSourceFormatLink   ApiSpecImportSourceFormat = "link"
)

func PossibleValuesForApiSpecImportSourceFormat() []string {
	return []string{
		string(ApiSpecImportSourceFormatInline),
		string(ApiSpecImportSourceFormatLink),
	}
}

func (s *ApiSpecImportSourceFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApiSpecImportSourceFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApiSpecImportSourceFormat(input string) (*ApiSpecImportSourceFormat, error) {
	vals := map[string]ApiSpecImportSourceFormat{
		"inline": ApiSpecImportSourceFormatInline,
		"link":   ApiSpecImportSourceFormatLink,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApiSpecImportSourceFormat(input)
	return &out, nil
}
