package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OneLakeArtifact = LakeHouseArtifact{}

type LakeHouseArtifact struct {

	// Fields inherited from OneLakeArtifact

	ArtifactName string              `json:"artifactName"`
	ArtifactType OneLakeArtifactType `json:"artifactType"`
}

func (s LakeHouseArtifact) OneLakeArtifact() BaseOneLakeArtifactImpl {
	return BaseOneLakeArtifactImpl{
		ArtifactName: s.ArtifactName,
		ArtifactType: s.ArtifactType,
	}
}

var _ json.Marshaler = LakeHouseArtifact{}

func (s LakeHouseArtifact) MarshalJSON() ([]byte, error) {
	type wrapper LakeHouseArtifact
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LakeHouseArtifact: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LakeHouseArtifact: %+v", err)
	}

	decoded["artifactType"] = "LakeHouse"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LakeHouseArtifact: %+v", err)
	}

	return encoded, nil
}
