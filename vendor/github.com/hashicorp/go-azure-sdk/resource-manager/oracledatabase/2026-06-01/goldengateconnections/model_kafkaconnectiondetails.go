package goldengateconnections

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectionBaseProperties = KafkaConnectionDetails{}

type KafkaConnectionDetails struct {
	BootstrapServers           *[]KafkaBootstrapServer       `json:"bootstrapServers,omitempty"`
	ClusterId                  *string                       `json:"clusterId,omitempty"`
	ConsumerProperties         *string                       `json:"consumerProperties,omitempty"`
	KeyStorePasswordSecretId   *string                       `json:"keyStorePasswordSecretId,omitempty"`
	KeyStoreSecretId           *string                       `json:"keyStoreSecretId,omitempty"`
	PasswordSecretId           *string                       `json:"passwordSecretId,omitempty"`
	ProducerProperties         *string                       `json:"producerProperties,omitempty"`
	SecurityProtocol           *string                       `json:"securityProtocol,omitempty"`
	ShouldUseResourcePrincipal *bool                         `json:"shouldUseResourcePrincipal,omitempty"`
	SslKeyPasswordSecretId     *string                       `json:"sslKeyPasswordSecretId,omitempty"`
	StreamPoolId               *string                       `json:"streamPoolId,omitempty"`
	TechnologyType             KafkaConnectionTechnologyType `json:"technologyType"`
	TrustStorePasswordSecretId *string                       `json:"trustStorePasswordSecretId,omitempty"`
	TrustStoreSecretId         *string                       `json:"trustStoreSecretId,omitempty"`
	Username                   *string                       `json:"username,omitempty"`

	// Fields inherited from ConnectionBaseProperties

	CompartmentId     *string                         `json:"compartmentId,omitempty"`
	ConnectionType    ConnectionType                  `json:"connectionType"`
	DisplayName       string                          `json:"displayName"`
	DoesUseSecretIds  *bool                           `json:"doesUseSecretIds,omitempty"`
	KeyId             *string                         `json:"keyId,omitempty"`
	LifecycleDetails  *string                         `json:"lifecycleDetails,omitempty"`
	LifecycleState    *ConnectionLifecycleState       `json:"lifecycleState,omitempty"`
	NetworkAnchorId   string                          `json:"networkAnchorId"`
	Ocid              *string                         `json:"ocid,omitempty"`
	ProvisioningState *AzureResourceProvisioningState `json:"provisioningState,omitempty"`
	ResourceAnchorId  string                          `json:"resourceAnchorId"`
	RoutingMethod     *RoutingMethod                  `json:"routingMethod,omitempty"`
	TimeCreated       *string                         `json:"timeCreated,omitempty"`
	TimeUpdated       *string                         `json:"timeUpdated,omitempty"`
	VaultId           *string                         `json:"vaultId,omitempty"`
}

func (s KafkaConnectionDetails) ConnectionBaseProperties() BaseConnectionBasePropertiesImpl {
	return BaseConnectionBasePropertiesImpl{
		CompartmentId:     s.CompartmentId,
		ConnectionType:    s.ConnectionType,
		DisplayName:       s.DisplayName,
		DoesUseSecretIds:  s.DoesUseSecretIds,
		KeyId:             s.KeyId,
		LifecycleDetails:  s.LifecycleDetails,
		LifecycleState:    s.LifecycleState,
		NetworkAnchorId:   s.NetworkAnchorId,
		Ocid:              s.Ocid,
		ProvisioningState: s.ProvisioningState,
		ResourceAnchorId:  s.ResourceAnchorId,
		RoutingMethod:     s.RoutingMethod,
		TimeCreated:       s.TimeCreated,
		TimeUpdated:       s.TimeUpdated,
		VaultId:           s.VaultId,
	}
}

var _ json.Marshaler = KafkaConnectionDetails{}

func (s KafkaConnectionDetails) MarshalJSON() ([]byte, error) {
	type wrapper KafkaConnectionDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KafkaConnectionDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KafkaConnectionDetails: %+v", err)
	}

	decoded["connectionType"] = "KAFKA"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KafkaConnectionDetails: %+v", err)
	}

	return encoded, nil
}
