package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureDataLakeStoreLinkedServiceTypeProperties struct {
	AccountName         *interface{}         `json:"accountName,omitempty"`
	AzureCloudType      *interface{}         `json:"azureCloudType,omitempty"`
	Credential          *CredentialReference `json:"credential,omitempty"`
	DataLakeStoreUri    interface{}          `json:"dataLakeStoreUri"`
	EncryptedCredential *string              `json:"encryptedCredential,omitempty"`
	ResourceGroupName   *interface{}         `json:"resourceGroupName,omitempty"`
	ServicePrincipalId  *interface{}         `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey SecretBase           `json:"servicePrincipalKey"`
	SubscriptionId      *interface{}         `json:"subscriptionId,omitempty"`
	Tenant              *interface{}         `json:"tenant,omitempty"`
}

var _ json.Unmarshaler = &AzureDataLakeStoreLinkedServiceTypeProperties{}

func (s *AzureDataLakeStoreLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccountName         *interface{}         `json:"accountName,omitempty"`
		AzureCloudType      *interface{}         `json:"azureCloudType,omitempty"`
		Credential          *CredentialReference `json:"credential,omitempty"`
		DataLakeStoreUri    interface{}          `json:"dataLakeStoreUri"`
		EncryptedCredential *string              `json:"encryptedCredential,omitempty"`
		ResourceGroupName   *interface{}         `json:"resourceGroupName,omitempty"`
		ServicePrincipalId  *interface{}         `json:"servicePrincipalId,omitempty"`
		SubscriptionId      *interface{}         `json:"subscriptionId,omitempty"`
		Tenant              *interface{}         `json:"tenant,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AccountName = decoded.AccountName
	s.AzureCloudType = decoded.AzureCloudType
	s.Credential = decoded.Credential
	s.DataLakeStoreUri = decoded.DataLakeStoreUri
	s.EncryptedCredential = decoded.EncryptedCredential
	s.ResourceGroupName = decoded.ResourceGroupName
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.SubscriptionId = decoded.SubscriptionId
	s.Tenant = decoded.Tenant

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureDataLakeStoreLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'AzureDataLakeStoreLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
