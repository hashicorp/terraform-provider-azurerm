package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerDataIdentity = SearchIndexerDataNoneIdentity{}

type SearchIndexerDataNoneIdentity struct {

	// Fields inherited from SearchIndexerDataIdentity

	OdataType string `json:"@odata.type"`
}

func (s SearchIndexerDataNoneIdentity) SearchIndexerDataIdentity() BaseSearchIndexerDataIdentityImpl {
	return BaseSearchIndexerDataIdentityImpl{
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = SearchIndexerDataNoneIdentity{}

func (s SearchIndexerDataNoneIdentity) MarshalJSON() ([]byte, error) {
	type wrapper SearchIndexerDataNoneIdentity
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SearchIndexerDataNoneIdentity: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SearchIndexerDataNoneIdentity: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.DataNoneIdentity"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SearchIndexerDataNoneIdentity: %+v", err)
	}

	return encoded, nil
}
