package datasources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Prefer string

const (
	PreferReturnRepresentation Prefer = "return=representation"
)

func PossibleValuesForPrefer() []string {
	return []string{
		string(PreferReturnRepresentation),
	}
}

func (s *Prefer) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrefer(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrefer(input string) (*Prefer, error) {
	vals := map[string]Prefer{
		"return=representation": PreferReturnRepresentation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Prefer(input)
	return &out, nil
}

type SearchIndexerDataSourceType string

const (
	SearchIndexerDataSourceTypeAdlsgenTwo SearchIndexerDataSourceType = "adlsgen2"
	SearchIndexerDataSourceTypeAzureblob  SearchIndexerDataSourceType = "azureblob"
	SearchIndexerDataSourceTypeAzuresql   SearchIndexerDataSourceType = "azuresql"
	SearchIndexerDataSourceTypeAzuretable SearchIndexerDataSourceType = "azuretable"
	SearchIndexerDataSourceTypeCosmosdb   SearchIndexerDataSourceType = "cosmosdb"
	SearchIndexerDataSourceTypeMysql      SearchIndexerDataSourceType = "mysql"
	SearchIndexerDataSourceTypeOnelake    SearchIndexerDataSourceType = "onelake"
)

func PossibleValuesForSearchIndexerDataSourceType() []string {
	return []string{
		string(SearchIndexerDataSourceTypeAdlsgenTwo),
		string(SearchIndexerDataSourceTypeAzureblob),
		string(SearchIndexerDataSourceTypeAzuresql),
		string(SearchIndexerDataSourceTypeAzuretable),
		string(SearchIndexerDataSourceTypeCosmosdb),
		string(SearchIndexerDataSourceTypeMysql),
		string(SearchIndexerDataSourceTypeOnelake),
	}
}

func (s *SearchIndexerDataSourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSearchIndexerDataSourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSearchIndexerDataSourceType(input string) (*SearchIndexerDataSourceType, error) {
	vals := map[string]SearchIndexerDataSourceType{
		"adlsgen2":   SearchIndexerDataSourceTypeAdlsgenTwo,
		"azureblob":  SearchIndexerDataSourceTypeAzureblob,
		"azuresql":   SearchIndexerDataSourceTypeAzuresql,
		"azuretable": SearchIndexerDataSourceTypeAzuretable,
		"cosmosdb":   SearchIndexerDataSourceTypeCosmosdb,
		"mysql":      SearchIndexerDataSourceTypeMysql,
		"onelake":    SearchIndexerDataSourceTypeOnelake,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SearchIndexerDataSourceType(input)
	return &out, nil
}
