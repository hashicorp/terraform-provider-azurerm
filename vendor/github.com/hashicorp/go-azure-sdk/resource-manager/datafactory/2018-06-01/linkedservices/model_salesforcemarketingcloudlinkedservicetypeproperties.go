package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SalesforceMarketingCloudLinkedServiceTypeProperties struct {
	ClientId              *interface{} `json:"clientId,omitempty"`
	ClientSecret          SecretBase   `json:"clientSecret"`
	ConnectionProperties  *interface{} `json:"connectionProperties,omitempty"`
	EncryptedCredential   *string      `json:"encryptedCredential,omitempty"`
	UseEncryptedEndpoints *bool        `json:"useEncryptedEndpoints,omitempty"`
	UseHostVerification   *bool        `json:"useHostVerification,omitempty"`
	UsePeerVerification   *bool        `json:"usePeerVerification,omitempty"`
}

var _ json.Unmarshaler = &SalesforceMarketingCloudLinkedServiceTypeProperties{}

func (s *SalesforceMarketingCloudLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ClientId              *interface{} `json:"clientId,omitempty"`
		ConnectionProperties  *interface{} `json:"connectionProperties,omitempty"`
		EncryptedCredential   *string      `json:"encryptedCredential,omitempty"`
		UseEncryptedEndpoints *bool        `json:"useEncryptedEndpoints,omitempty"`
		UseHostVerification   *bool        `json:"useHostVerification,omitempty"`
		UsePeerVerification   *bool        `json:"usePeerVerification,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ClientId = decoded.ClientId
	s.ConnectionProperties = decoded.ConnectionProperties
	s.EncryptedCredential = decoded.EncryptedCredential
	s.UseEncryptedEndpoints = decoded.UseEncryptedEndpoints
	s.UseHostVerification = decoded.UseHostVerification
	s.UsePeerVerification = decoded.UsePeerVerification

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SalesforceMarketingCloudLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["clientSecret"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ClientSecret' for 'SalesforceMarketingCloudLinkedServiceTypeProperties': %+v", err)
		}
		s.ClientSecret = impl
	}

	return nil
}
