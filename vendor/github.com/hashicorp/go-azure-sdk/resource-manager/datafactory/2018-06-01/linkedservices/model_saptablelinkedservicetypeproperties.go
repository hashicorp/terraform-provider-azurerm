package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapTableLinkedServiceTypeProperties struct {
	ClientId             *interface{} `json:"clientId,omitempty"`
	EncryptedCredential  *string      `json:"encryptedCredential,omitempty"`
	Language             *interface{} `json:"language,omitempty"`
	LogonGroup           *interface{} `json:"logonGroup,omitempty"`
	MessageServer        *interface{} `json:"messageServer,omitempty"`
	MessageServerService *interface{} `json:"messageServerService,omitempty"`
	Password             SecretBase   `json:"password"`
	Server               *interface{} `json:"server,omitempty"`
	SncLibraryPath       *interface{} `json:"sncLibraryPath,omitempty"`
	SncMode              *bool        `json:"sncMode,omitempty"`
	SncMyName            *interface{} `json:"sncMyName,omitempty"`
	SncPartnerName       *interface{} `json:"sncPartnerName,omitempty"`
	SncQop               *interface{} `json:"sncQop,omitempty"`
	SystemId             *interface{} `json:"systemId,omitempty"`
	SystemNumber         *interface{} `json:"systemNumber,omitempty"`
	UserName             *interface{} `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &SapTableLinkedServiceTypeProperties{}

func (s *SapTableLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ClientId             *interface{} `json:"clientId,omitempty"`
		EncryptedCredential  *string      `json:"encryptedCredential,omitempty"`
		Language             *interface{} `json:"language,omitempty"`
		LogonGroup           *interface{} `json:"logonGroup,omitempty"`
		MessageServer        *interface{} `json:"messageServer,omitempty"`
		MessageServerService *interface{} `json:"messageServerService,omitempty"`
		Server               *interface{} `json:"server,omitempty"`
		SncLibraryPath       *interface{} `json:"sncLibraryPath,omitempty"`
		SncMode              *bool        `json:"sncMode,omitempty"`
		SncMyName            *interface{} `json:"sncMyName,omitempty"`
		SncPartnerName       *interface{} `json:"sncPartnerName,omitempty"`
		SncQop               *interface{} `json:"sncQop,omitempty"`
		SystemId             *interface{} `json:"systemId,omitempty"`
		SystemNumber         *interface{} `json:"systemNumber,omitempty"`
		UserName             *interface{} `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ClientId = decoded.ClientId
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Language = decoded.Language
	s.LogonGroup = decoded.LogonGroup
	s.MessageServer = decoded.MessageServer
	s.MessageServerService = decoded.MessageServerService
	s.Server = decoded.Server
	s.SncLibraryPath = decoded.SncLibraryPath
	s.SncMode = decoded.SncMode
	s.SncMyName = decoded.SncMyName
	s.SncPartnerName = decoded.SncPartnerName
	s.SncQop = decoded.SncQop
	s.SystemId = decoded.SystemId
	s.SystemNumber = decoded.SystemNumber
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SapTableLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'SapTableLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
