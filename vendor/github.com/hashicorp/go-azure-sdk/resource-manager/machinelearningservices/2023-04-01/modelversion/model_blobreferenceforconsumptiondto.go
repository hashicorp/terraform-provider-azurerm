package modelversion

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobReferenceForConsumptionDto struct {
	BlobUri             *string                    `json:"blobUri,omitempty"`
	Credential          PendingUploadCredentialDto `json:"credential"`
	StorageAccountArmId *string                    `json:"storageAccountArmId,omitempty"`
}

var _ json.Unmarshaler = &BlobReferenceForConsumptionDto{}

func (s *BlobReferenceForConsumptionDto) UnmarshalJSON(bytes []byte) error {
	type alias BlobReferenceForConsumptionDto
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into BlobReferenceForConsumptionDto: %+v", err)
	}

	s.BlobUri = decoded.BlobUri
	s.StorageAccountArmId = decoded.StorageAccountArmId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BlobReferenceForConsumptionDto into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["credential"]; ok {
		impl, err := unmarshalPendingUploadCredentialDtoImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Credential' for 'BlobReferenceForConsumptionDto': %+v", err)
		}
		s.Credential = impl
	}
	return nil
}
