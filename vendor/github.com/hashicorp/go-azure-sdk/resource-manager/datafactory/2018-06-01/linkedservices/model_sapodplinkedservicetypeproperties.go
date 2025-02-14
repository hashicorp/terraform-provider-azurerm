package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapOdpLinkedServiceTypeProperties struct {
	ClientId             *string    `json:"clientId,omitempty"`
	EncryptedCredential  *string    `json:"encryptedCredential,omitempty"`
	Language             *string    `json:"language,omitempty"`
	LogonGroup           *string    `json:"logonGroup,omitempty"`
	MessageServer        *string    `json:"messageServer,omitempty"`
	MessageServerService *string    `json:"messageServerService,omitempty"`
	Password             SecretBase `json:"password"`
	Server               *string    `json:"server,omitempty"`
	SncLibraryPath       *string    `json:"sncLibraryPath,omitempty"`
	SncMode              *bool      `json:"sncMode,omitempty"`
	SncMyName            *string    `json:"sncMyName,omitempty"`
	SncPartnerName       *string    `json:"sncPartnerName,omitempty"`
	SncQop               *string    `json:"sncQop,omitempty"`
	SubscriberName       *string    `json:"subscriberName,omitempty"`
	SystemId             *string    `json:"systemId,omitempty"`
	SystemNumber         *string    `json:"systemNumber,omitempty"`
	UserName             *string    `json:"userName,omitempty"`
	X509CertificatePath  *string    `json:"x509CertificatePath,omitempty"`
}

var _ json.Unmarshaler = &SapOdpLinkedServiceTypeProperties{}

func (s *SapOdpLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ClientId             *string `json:"clientId,omitempty"`
		EncryptedCredential  *string `json:"encryptedCredential,omitempty"`
		Language             *string `json:"language,omitempty"`
		LogonGroup           *string `json:"logonGroup,omitempty"`
		MessageServer        *string `json:"messageServer,omitempty"`
		MessageServerService *string `json:"messageServerService,omitempty"`
		Server               *string `json:"server,omitempty"`
		SncLibraryPath       *string `json:"sncLibraryPath,omitempty"`
		SncMode              *bool   `json:"sncMode,omitempty"`
		SncMyName            *string `json:"sncMyName,omitempty"`
		SncPartnerName       *string `json:"sncPartnerName,omitempty"`
		SncQop               *string `json:"sncQop,omitempty"`
		SubscriberName       *string `json:"subscriberName,omitempty"`
		SystemId             *string `json:"systemId,omitempty"`
		SystemNumber         *string `json:"systemNumber,omitempty"`
		UserName             *string `json:"userName,omitempty"`
		X509CertificatePath  *string `json:"x509CertificatePath,omitempty"`
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
	s.SubscriberName = decoded.SubscriberName
	s.SystemId = decoded.SystemId
	s.SystemNumber = decoded.SystemNumber
	s.UserName = decoded.UserName
	s.X509CertificatePath = decoded.X509CertificatePath

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SapOdpLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'SapOdpLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
