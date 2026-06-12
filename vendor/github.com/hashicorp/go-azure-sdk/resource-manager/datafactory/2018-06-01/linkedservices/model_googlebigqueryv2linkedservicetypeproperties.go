package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GoogleBigQueryV2LinkedServiceTypeProperties struct {
	AuthenticationType  GoogleBigQueryV2AuthenticationType `json:"authenticationType"`
	ClientId            *interface{}                       `json:"clientId,omitempty"`
	ClientSecret        SecretBase                         `json:"clientSecret"`
	EncryptedCredential *string                            `json:"encryptedCredential,omitempty"`
	KeyFileContent      SecretBase                         `json:"keyFileContent"`
	ProjectId           interface{}                        `json:"projectId"`
	RefreshToken        SecretBase                         `json:"refreshToken"`
}

var _ json.Unmarshaler = &GoogleBigQueryV2LinkedServiceTypeProperties{}

func (s *GoogleBigQueryV2LinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType  GoogleBigQueryV2AuthenticationType `json:"authenticationType"`
		ClientId            *interface{}                       `json:"clientId,omitempty"`
		EncryptedCredential *string                            `json:"encryptedCredential,omitempty"`
		ProjectId           interface{}                        `json:"projectId"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.ClientId = decoded.ClientId
	s.EncryptedCredential = decoded.EncryptedCredential
	s.ProjectId = decoded.ProjectId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling GoogleBigQueryV2LinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["clientSecret"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ClientSecret' for 'GoogleBigQueryV2LinkedServiceTypeProperties': %+v", err)
		}
		s.ClientSecret = impl
	}

	if v, ok := temp["keyFileContent"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'KeyFileContent' for 'GoogleBigQueryV2LinkedServiceTypeProperties': %+v", err)
		}
		s.KeyFileContent = impl
	}

	if v, ok := temp["refreshToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RefreshToken' for 'GoogleBigQueryV2LinkedServiceTypeProperties': %+v", err)
		}
		s.RefreshToken = impl
	}

	return nil
}
