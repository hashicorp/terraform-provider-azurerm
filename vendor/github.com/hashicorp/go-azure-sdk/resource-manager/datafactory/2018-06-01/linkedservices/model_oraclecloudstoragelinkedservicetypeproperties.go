package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OracleCloudStorageLinkedServiceTypeProperties struct {
	AccessKeyId         *string    `json:"accessKeyId,omitempty"`
	EncryptedCredential *string    `json:"encryptedCredential,omitempty"`
	SecretAccessKey     SecretBase `json:"secretAccessKey"`
	ServiceURL          *string    `json:"serviceUrl,omitempty"`
}

var _ json.Unmarshaler = &OracleCloudStorageLinkedServiceTypeProperties{}

func (s *OracleCloudStorageLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccessKeyId         *string `json:"accessKeyId,omitempty"`
		EncryptedCredential *string `json:"encryptedCredential,omitempty"`
		ServiceURL          *string `json:"serviceUrl,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AccessKeyId = decoded.AccessKeyId
	s.EncryptedCredential = decoded.EncryptedCredential
	s.ServiceURL = decoded.ServiceURL

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling OracleCloudStorageLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["secretAccessKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'SecretAccessKey' for 'OracleCloudStorageLinkedServiceTypeProperties': %+v", err)
		}
		s.SecretAccessKey = impl
	}

	return nil
}
