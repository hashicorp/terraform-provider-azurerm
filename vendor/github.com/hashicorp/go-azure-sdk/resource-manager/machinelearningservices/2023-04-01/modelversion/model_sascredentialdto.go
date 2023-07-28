package modelversion

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PendingUploadCredentialDto = SASCredentialDto{}

type SASCredentialDto struct {
	SasUri *string `json:"sasUri,omitempty"`

	// Fields inherited from PendingUploadCredentialDto
}

var _ json.Marshaler = SASCredentialDto{}

func (s SASCredentialDto) MarshalJSON() ([]byte, error) {
	type wrapper SASCredentialDto
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SASCredentialDto: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SASCredentialDto: %+v", err)
	}
	decoded["credentialType"] = "SAS"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SASCredentialDto: %+v", err)
	}

	return encoded, nil
}
