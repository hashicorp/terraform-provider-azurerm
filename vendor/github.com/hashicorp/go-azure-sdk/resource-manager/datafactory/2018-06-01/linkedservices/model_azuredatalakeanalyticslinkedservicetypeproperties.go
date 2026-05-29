package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureDataLakeAnalyticsLinkedServiceTypeProperties struct {
	AccountName          interface{}  `json:"accountName"`
	DataLakeAnalyticsUri *interface{} `json:"dataLakeAnalyticsUri,omitempty"`
	EncryptedCredential  *string      `json:"encryptedCredential,omitempty"`
	ResourceGroupName    *interface{} `json:"resourceGroupName,omitempty"`
	ServicePrincipalId   *interface{} `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey  SecretBase   `json:"servicePrincipalKey"`
	SubscriptionId       *interface{} `json:"subscriptionId,omitempty"`
	Tenant               interface{}  `json:"tenant"`
}

var _ json.Unmarshaler = &AzureDataLakeAnalyticsLinkedServiceTypeProperties{}

func (s *AzureDataLakeAnalyticsLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccountName          interface{}  `json:"accountName"`
		DataLakeAnalyticsUri *interface{} `json:"dataLakeAnalyticsUri,omitempty"`
		EncryptedCredential  *string      `json:"encryptedCredential,omitempty"`
		ResourceGroupName    *interface{} `json:"resourceGroupName,omitempty"`
		ServicePrincipalId   *interface{} `json:"servicePrincipalId,omitempty"`
		SubscriptionId       *interface{} `json:"subscriptionId,omitempty"`
		Tenant               interface{}  `json:"tenant"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AccountName = decoded.AccountName
	s.DataLakeAnalyticsUri = decoded.DataLakeAnalyticsUri
	s.EncryptedCredential = decoded.EncryptedCredential
	s.ResourceGroupName = decoded.ResourceGroupName
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.SubscriptionId = decoded.SubscriptionId
	s.Tenant = decoded.Tenant

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureDataLakeAnalyticsLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'AzureDataLakeAnalyticsLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
