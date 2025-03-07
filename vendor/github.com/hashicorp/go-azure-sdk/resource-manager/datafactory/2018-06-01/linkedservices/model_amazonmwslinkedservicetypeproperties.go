package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmazonMWSLinkedServiceTypeProperties struct {
	AccessKeyId           interface{} `json:"accessKeyId"`
	EncryptedCredential   *string     `json:"encryptedCredential,omitempty"`
	Endpoint              interface{} `json:"endpoint"`
	MarketplaceID         interface{} `json:"marketplaceID"`
	MwsAuthToken          SecretBase  `json:"mwsAuthToken"`
	SecretKey             SecretBase  `json:"secretKey"`
	SellerID              interface{} `json:"sellerID"`
	UseEncryptedEndpoints *bool       `json:"useEncryptedEndpoints,omitempty"`
	UseHostVerification   *bool       `json:"useHostVerification,omitempty"`
	UsePeerVerification   *bool       `json:"usePeerVerification,omitempty"`
}

var _ json.Unmarshaler = &AmazonMWSLinkedServiceTypeProperties{}

func (s *AmazonMWSLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccessKeyId           interface{} `json:"accessKeyId"`
		EncryptedCredential   *string     `json:"encryptedCredential,omitempty"`
		Endpoint              interface{} `json:"endpoint"`
		MarketplaceID         interface{} `json:"marketplaceID"`
		SellerID              interface{} `json:"sellerID"`
		UseEncryptedEndpoints *bool       `json:"useEncryptedEndpoints,omitempty"`
		UseHostVerification   *bool       `json:"useHostVerification,omitempty"`
		UsePeerVerification   *bool       `json:"usePeerVerification,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AccessKeyId = decoded.AccessKeyId
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Endpoint = decoded.Endpoint
	s.MarketplaceID = decoded.MarketplaceID
	s.SellerID = decoded.SellerID
	s.UseEncryptedEndpoints = decoded.UseEncryptedEndpoints
	s.UseHostVerification = decoded.UseHostVerification
	s.UsePeerVerification = decoded.UsePeerVerification

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AmazonMWSLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["mwsAuthToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'MwsAuthToken' for 'AmazonMWSLinkedServiceTypeProperties': %+v", err)
		}
		s.MwsAuthToken = impl
	}

	if v, ok := temp["secretKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'SecretKey' for 'AmazonMWSLinkedServiceTypeProperties': %+v", err)
		}
		s.SecretKey = impl
	}

	return nil
}
