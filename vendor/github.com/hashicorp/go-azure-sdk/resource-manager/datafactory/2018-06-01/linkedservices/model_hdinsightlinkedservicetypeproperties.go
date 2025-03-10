package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HDInsightLinkedServiceTypeProperties struct {
	ClusterUri                interface{}             `json:"clusterUri"`
	EncryptedCredential       *string                 `json:"encryptedCredential,omitempty"`
	FileSystem                *interface{}            `json:"fileSystem,omitempty"`
	HcatalogLinkedServiceName *LinkedServiceReference `json:"hcatalogLinkedServiceName,omitempty"`
	IsEspEnabled              *bool                   `json:"isEspEnabled,omitempty"`
	LinkedServiceName         *LinkedServiceReference `json:"linkedServiceName,omitempty"`
	Password                  SecretBase              `json:"password"`
	UserName                  *interface{}            `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &HDInsightLinkedServiceTypeProperties{}

func (s *HDInsightLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ClusterUri                interface{}             `json:"clusterUri"`
		EncryptedCredential       *string                 `json:"encryptedCredential,omitempty"`
		FileSystem                *interface{}            `json:"fileSystem,omitempty"`
		HcatalogLinkedServiceName *LinkedServiceReference `json:"hcatalogLinkedServiceName,omitempty"`
		IsEspEnabled              *bool                   `json:"isEspEnabled,omitempty"`
		LinkedServiceName         *LinkedServiceReference `json:"linkedServiceName,omitempty"`
		UserName                  *interface{}            `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ClusterUri = decoded.ClusterUri
	s.EncryptedCredential = decoded.EncryptedCredential
	s.FileSystem = decoded.FileSystem
	s.HcatalogLinkedServiceName = decoded.HcatalogLinkedServiceName
	s.IsEspEnabled = decoded.IsEspEnabled
	s.LinkedServiceName = decoded.LinkedServiceName
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling HDInsightLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'HDInsightLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
