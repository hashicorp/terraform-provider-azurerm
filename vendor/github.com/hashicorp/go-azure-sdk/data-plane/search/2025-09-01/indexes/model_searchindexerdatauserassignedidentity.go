package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerDataIdentity = SearchIndexerDataUserAssignedIdentity{}

type SearchIndexerDataUserAssignedIdentity struct {
	UserAssignedIdentity string `json:"userAssignedIdentity"`

	// Fields inherited from SearchIndexerDataIdentity

	OdataType string `json:"@odata.type"`
}

func (s SearchIndexerDataUserAssignedIdentity) SearchIndexerDataIdentity() BaseSearchIndexerDataIdentityImpl {
	return BaseSearchIndexerDataIdentityImpl{
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = SearchIndexerDataUserAssignedIdentity{}

func (s SearchIndexerDataUserAssignedIdentity) MarshalJSON() ([]byte, error) {
	type wrapper SearchIndexerDataUserAssignedIdentity
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SearchIndexerDataUserAssignedIdentity: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SearchIndexerDataUserAssignedIdentity: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.DataUserAssignedIdentity"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SearchIndexerDataUserAssignedIdentity: %+v", err)
	}

	return encoded, nil
}
